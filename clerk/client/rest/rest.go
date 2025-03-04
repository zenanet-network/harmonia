package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/zenanet-network/harmonia/helper"
)

// RestLogger for staking module logger
var RestLogger tmLog.Logger

func init() {
	RestLogger = helper.Logger.With("module", "clerk/rest")
}

// RegisterRoutes registers checkpoint-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
