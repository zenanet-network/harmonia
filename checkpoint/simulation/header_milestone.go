package simulation

import (
	"time"

	"github.com/zenanet-network/harmonia/types"
)

const TestMilestoneID = "17ce48fe-0a18-41a8-ab7e-59d8002f027b - 0x901a64406d97a3fa9b87b320cbeb86b3c62328f5"

// GenRandMilestone return headers
func GenRandMilestone(start uint64, sprintLength uint64) (milestone types.Milestone, err error) {
	end := start + sprintLength - 1
	eireneChainID := "1234"
	rootHash := types.HexToHarmoniaHash("123")
	proposer := types.HarmoniaAddress{}
	milestoneID := TestMilestoneID
	milestone = types.CreateMilestone(
		start,
		end,
		rootHash,
		proposer,
		eireneChainID,
		milestoneID,
		uint64(time.Now().UTC().Unix()))

	return milestone, nil
}
