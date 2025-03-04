package types

import (
	"github.com/zenanet-network/harmonia/auth/exported"
)

// AccountProcessor is an interface to process account as per module
type AccountProcessor func(*GenesisAccount, *BaseAccount) exported.Account
