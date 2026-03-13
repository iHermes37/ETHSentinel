package account

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// DefaultRootDerivationPath is the root path to which custom derivation endpoints
	// are appended. As such, the first account will be at m/44'/60'/0'/0, the second
	// at m/44'/60'/0'/1, etc.
	DefaultRootDerivationPath = accounts.DefaultRootDerivationPath

	// DefaultBaseDerivationPath is the base path from which custom derivation endpoints
	// are incremented. As such, the first account will be at m/44'/60'/0'/0, the second
	// at m/44'/60'/0'/1, etc
	DefaultBaseDerivationPath = accounts.DefaultBaseDerivationPath

	issue179FixEnvar = "GO_ETHEREUM_HDWALLET_FIX_ISSUE_179"
)

type Wallet struct {
	masterKey   *hdkeychain.ExtendedKey
	fixIssue172 bool
	seed        Seed
	paths       map[common.Address]accounts.DerivationPath
	accounts    []accounts.Account

	mnemonic Mnemonic
}

type Mnemonic struct {
	MnemonicStr string
}
type Seed struct {
	seedByte []byte
}
type DerivationPath struct{}
