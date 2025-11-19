package signature

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// SignHash implements accounts.Wallet, which allows signing arbitrary data.
func (w *Wallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(hash, privateKey)
}

// SignTxWithSigner allows the account to sign a transaction using a custom signer.
func (w *Wallet) SignTxWithSigner(account accounts.Account, tx *types.Transaction, signer types.Signer) (*types.Transaction, error) {
	if tx == nil {
		return nil, errors.New("nil transaction")
	}
	if signer == nil {
		return nil, errors.New("nil signer")
	}

	w.stateLock.RLock() // Comms have own mutex, this is for the state fields
	defer w.stateLock.RUnlock()

	// Make sure the requested account is contained within
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}

	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

// SignTxEIP1559 uses the London Signer which supports
// EIP-1559 dynamic fee, EIP-2930 access list, EIP-155 replay protected, and legacy Homestead Transactions
func (w *Wallet) SignTxEIP1559(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.SignTxWithSigner(account, tx, types.NewLondonSigner(chainID))
}

// SignTxEIP155 implements accounts.Wallet, which allows the account to sign an ERC-20 transaction.
func (w *Wallet) SignTxEIP155(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.SignTxWithSigner(account, tx, types.NewEIP155Signer(chainID))
}

// SignTx implements accounts.Wallet, which allows the account to sign an Ethereum transaction.
func (w *Wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.SignTxWithSigner(account, tx, types.LatestSignerForChainID(chainID))
}
