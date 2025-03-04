package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	hmModule "github.com/zenanet-network/harmonia/types/module"
)

type App interface {
	Name() string

	Codec() *codec.Codec

	BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock

	EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock

	InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain

	LoadHeight(height int64) error

	ExportAppStateAndValidators() (json.RawMessage, []tmTypes.GenesisValidator, error)

	ModuleAccountAddrs() map[string]bool

	SimulationManager() *hmModule.SimulationManager
}
