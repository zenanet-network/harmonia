package simulation

import (
	"github.com/zenanet-network/harmonia/bank/types"
	"github.com/zenanet-network/harmonia/types/module"
)

// RandomizedGenState returns bank genesis
func RandomizedGenState(simState *module.SimulationState) {
	bankGenesis := types.NewGenesisState(true)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bankGenesis)
}
