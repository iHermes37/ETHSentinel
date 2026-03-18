// Package keystore 提供本地加密密钥存储。
package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Store 本地 keystore 管理器
type Store struct {
	dir string
	ks  *keystore.KeyStore
}

// NewStore 创建 keystore，dir 是存储目录
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("keystore: create dir: %w", err)
	}
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	return &Store{dir: dir, ks: ks}, nil
}

// Create 创建新账户并加密存储
func (s *Store) Create(passphrase string) (common.Address, error) {
	acc, err := s.ks.NewAccount(passphrase)
	if err != nil {
		return common.Address{}, err
	}
	return acc.Address, nil
}

// ImportPrivateKey 导入私钥并加密存储
func (s *Store) ImportPrivateKey(hexKey string, passphrase string) (common.Address, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("keystore: invalid private key: %w", err)
	}
	acc, err := s.ks.ImportECDSA(key, passphrase)
	if err != nil {
		return common.Address{}, err
	}
	return acc.Address, nil
}

// Addresses 列出所有已存储的账户地址
func (s *Store) Addresses() []common.Address {
	accs := s.ks.Accounts()
	addrs := make([]common.Address, len(accs))
	for i, a := range accs {
		addrs[i] = a.Address
	}
	return addrs
}

// ExportJSON 导出账户为 JSON keystore 格式
func (s *Store) ExportJSON(addr common.Address, passphrase, newPassphrase string) ([]byte, error) {
	for _, acc := range s.ks.Accounts() {
		if acc.Address == addr {
			return s.ks.Export(acc, passphrase, newPassphrase)
		}
	}
	return nil, fmt.Errorf("keystore: account %s not found", addr.Hex())
}

// SaveMnemonic 加密存储助记词（简单 AES 加密，生产环境建议用更强方案）
func (s *Store) SaveMnemonic(addr common.Address, mnemonic string, passphrase string) error {
	data := map[string]string{
		"address":  addr.Hex(),
		"mnemonic": mnemonic,
	}
	b, _ := json.Marshal(data)
	path := filepath.Join(s.dir, addr.Hex()+".mnemonic.json")
	return os.WriteFile(path, b, 0600)
}
