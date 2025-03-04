package simulation

import (
	"time"

	"github.com/zenanet-network/harmonia/checkpoint/types"
	hmTypes "github.com/zenanet-network/harmonia/types"
	"github.com/zenanet-network/harmonia/types/module"
)

// RandomizedGenState return dummy genesis
func RandomizedGenState(simState *module.SimulationState) {
	lastNoACK := 0
	ackCount := 1
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")

	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	eireneChainID := "1234"

	bufferedCheckpoint := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		eireneChainID,
		timestamp,
	)

	Checkpoints := make([]hmTypes.Checkpoint, ackCount)

	for i := range Checkpoints {
		Checkpoints[i] = bufferedCheckpoint
	}

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		&bufferedCheckpoint,
		uint64(lastNoACK),
		uint64(ackCount),
		Checkpoints,
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
