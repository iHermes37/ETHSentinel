// Package hd 提供 BIP44 HD 钱包功能。
// 基于原始 account/hdwallet.go 整合重构。
package hd

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"sync"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bip39 "github.com/tyler-smith/go-bip39"
)

const issue179FixEnvar = "GO_ETHEREUM_HDWALLET_FIX_ISSUE_179"

// Wallet BIP44 HD 钱包
type Wallet struct {
	masterKey   *hdkeychain.ExtendedKey
	fixIssue172 bool
	seedBytes   []byte
	mnemonic    string

	stateLock sync.RWMutex
	paths     map[common.Address]accounts.DerivationPath
	accounts  []accounts.Account
	url       accounts.URL
}

// ─────────────────────────────────────────────
//  构造函数
// ─────────────────────────────────────────────

// NewFromMnemonic 从助记词创建钱包
func NewFromMnemonic(mnemonic string, passOpt ...string) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("hd: mnemonic is required")
	}
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("hd: mnemonic is invalid")
	}

	password := ""
	if len(passOpt) > 0 {
		password = passOpt[0]
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}

	w, err := newFromSeed(seed)
	if err != nil {
		return nil, err
	}
	w.mnemonic = mnemonic
	return w, nil
}

// NewFromSeed 从种子创建钱包
func NewFromSeed(seed []byte) (*Wallet, error) {
	if len(seed) == 0 {
		return nil, errors.New("hd: seed is required")
	}
	return newFromSeed(seed)
}

// GenerateMnemonic 生成新的助记词（128/160/192/224/256 bits）
func GenerateMnemonic(bits int) (string, error) {
	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// GenerateSeed 生成随机种子
func GenerateSeed() ([]byte, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	return b, err
}

func newFromSeed(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		masterKey:   masterKey,
		seedBytes:   seed,
		paths:       make(map[common.Address]accounts.DerivationPath),
		accounts:    []accounts.Account{},
		fixIssue172: len(os.Getenv(issue179FixEnvar)) > 0,
	}, nil
}

// ─────────────────────────────────────────────
//  派生
// ─────────────────────────────────────────────

// Derive 派生指定路径的账户，pin=true 则加入账户列表
func (w *Wallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	w.stateLock.RLock()
	address, err := w.deriveAddress(path)
	w.stateLock.RUnlock()

	if err != nil {
		return accounts.Account{}, err
	}

	account := accounts.Account{
		Address: address,
		URL:     accounts.URL{Path: path.String()},
	}

	if !pin {
		return account, nil
	}

	w.stateLock.Lock()
	defer w.stateLock.Unlock()
	if _, ok := w.paths[address]; !ok {
		w.accounts = append(w.accounts, account)
		w.paths[address] = path
	}
	return account, nil
}

// DeriveAt 派生第 index 个账户（使用默认路径 m/44'/60'/0'/0/index）
func (w *Wallet) DeriveAt(index uint32, pin bool) (accounts.Account, error) {
	path := accounts.DefaultBaseDerivationPath
	path[len(path)-1] = index
	return w.Derive(path, pin)
}

// ─────────────────────────────────────────────
//  账户信息
// ─────────────────────────────────────────────

func (w *Wallet) Accounts() []accounts.Account {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()
	cpy := make([]accounts.Account, len(w.accounts))
	copy(cpy, w.accounts)
	return cpy
}

func (w *Wallet) Contains(account accounts.Account) bool {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()
	_, exists := w.paths[account.Address]
	return exists
}

func (w *Wallet) Mnemonic() string { return w.mnemonic }

// ─────────────────────────────────────────────
//  密钥导出
// ─────────────────────────────────────────────

func (w *Wallet) PrivateKey(account accounts.Account) (*ecdsa.PrivateKey, error) {
	path, err := accounts.ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}
	return w.derivePrivateKey(path)
}

func (w *Wallet) PrivateKeyHex(account accounts.Account) (string, error) {
	key, err := w.PrivateKey(account)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(crypto.FromECDSA(key))[2:], nil
}

func (w *Wallet) PublicKey(account accounts.Account) (*ecdsa.PublicKey, error) {
	path, err := accounts.ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}
	return w.derivePublicKey(path)
}

func (w *Wallet) Address(account accounts.Account) (common.Address, error) {
	pub, err := w.PublicKey(account)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pub), nil
}

// ─────────────────────────────────────────────
//  签名
// ─────────────────────────────────────────────

// SignHash 对任意 hash 签名
func (w *Wallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}
	key, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}
	return crypto.Sign(hash, key)
}

// SignTx 签名交易（自动选择最新签名器）
func (w *Wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.signTxWithSigner(account, tx, types.LatestSignerForChainID(chainID))
}

// SignTxEIP1559 使用 London Signer 签名（EIP-1559）
func (w *Wallet) SignTxEIP1559(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.signTxWithSigner(account, tx, types.NewLondonSigner(chainID))
}

func (w *Wallet) signTxWithSigner(account accounts.Account, tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	if tx == nil {
		return nil, errors.New("hd: nil transaction")
	}
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}
	key, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}
	signed, err := types.SignTx(tx, signer, key)
	if err != nil {
		return nil, err
	}
	sender, err := types.Sender(signer, signed)
	if err != nil {
		return nil, err
	}
	if sender != account.Address {
		return nil, errors.New("hd: signer mismatch")
	}
	return signed, nil
}

// SelfDerive 自动派生（占位，实现参考 go-ethereum accounts.Wallet 接口）
func (w *Wallet) SelfDerive(_ []accounts.DerivationPath, _ geth.ChainStateReader) {}

// ─────────────────────────────────────────────
//  内部：派生私钥
// ─────────────────────────────────────────────

func (w *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	key := w.masterKey
	var err error
	for _, n := range path {
		if w.fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}
	priv, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privECDSA := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: crypto.S256()},
		D:         priv.ToECDSA().D,
	}
	privECDSA.PublicKey.X, privECDSA.PublicKey.Y = privECDSA.PublicKey.Curve.ScalarBaseMult(privECDSA.D.Bytes())
	return privECDSA, nil
}

func (w *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	priv, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}
	pub, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("hd: failed to get public key")
	}
	return pub, nil
}

func (w *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	pub, err := w.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pub), nil
}
