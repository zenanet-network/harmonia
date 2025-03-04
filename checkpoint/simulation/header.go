package simulation

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zenanet-network/harmonia/staking"
	stakingSim "github.com/zenanet-network/harmonia/staking/simulation"

	"github.com/zenanet-network/harmonia/types"
)

// GenRandCheckpoint return headers
func GenRandCheckpoint(start uint64, headerSize uint64, maxCheckpointLength uint64) (headerBlock types.Checkpoint, err error) {
	end := start + headerSize
	eireneChainID := "1234"
	rootHash := types.HexToHarmoniaHash("123")
	proposer := types.HarmoniaAddress{}

	headerBlock = types.CreateBlock(
		start,
		end,
		rootHash,
		proposer,
		eireneChainID,
		uint64(time.Now().UTC().Unix()))

	return headerBlock, nil
}

// LoadValidatorSet loads validator set
func LoadValidatorSet(t *testing.T, count int, keeper staking.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	t.Helper()

	var valSet types.ValidatorSet

	validators := stakingSim.GenRandomVal(count, 0, 10, uint64(timeAlive), randomise, 1)
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.NoError(t, err, "Unable to set validator, Error: %v", err)

		err = valSet.UpdateWithChangeSet([]*types.Validator{&validator})
		require.NoError(t, err)
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.NoError(t, err, "Unable to update validator set")

	vals := keeper.GetAllValidators(ctx)
	require.NotNil(t, vals)

	return valSet
}
