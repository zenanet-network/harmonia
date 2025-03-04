package simulation

import (
	"fmt"
	"math/rand"

	"github.com/zenanet-network/harmonia/chainmanager/types"
	"github.com/zenanet-network/harmonia/simulation"
	simtypes "github.com/zenanet-network/harmonia/types/simulation"
)

const (
	KeyMainchainTxConfirmations  = "MainchainTxConfirmations"
	KeyZenachainTxConfirmations = "ZenachainTxConfirmations"
	KeyChainParams               = "ChainParams"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, KeyMainchainTxConfirmations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMainchainTxConfirmations(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, KeyZenachainTxConfirmations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaticchainTxConfirmations(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, KeyChainParams,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBorChainId(r))
			},
		),
	}
}
