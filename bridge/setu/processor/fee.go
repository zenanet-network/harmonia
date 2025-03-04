package processor

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/zenanet-network/go-zenanet/accounts/abi"
	"github.com/zenanet-network/go-zenanet/core/types"

	"github.com/zenanet-network/harmonia/bridge/setu/util"
	"github.com/zenanet-network/harmonia/contracts/stakinginfo"
	"github.com/zenanet-network/harmonia/helper"
	topupTypes "github.com/zenanet-network/harmonia/topup/types"
	hmTypes "github.com/zenanet-network/harmonia/types"
)

// FeeProcessor - process fee related events
type FeeProcessor struct {
	BaseProcessor
	stakingInfoAbi *abi.ABI
}

// NewFeeProcessor - add  abi to clerk processor
func NewFeeProcessor(stakingInfoAbi *abi.ABI) *FeeProcessor {
	return &FeeProcessor{
		stakingInfoAbi: stakingInfoAbi,
	}
}

// Start starts new block subscription
func (fp *FeeProcessor) Start() error {
	fp.Logger.Info("Starting")
	return nil
}

// RegisterTasks - Registers clerk related tasks with machinery
func (fp *FeeProcessor) RegisterTasks() {
	fp.Logger.Info("Registering fee related tasks")

	if err := fp.queueConnector.Server.RegisterTask("sendTopUpFeeToHeimdall", fp.sendTopUpFeeToHeimdall); err != nil {
		fp.Logger.Error("RegisterTasks | sendTopUpFeeToHeimdall", "error", err)
	}
}

// processTopupFeeEvent - processes topup fee event
func (fp *FeeProcessor) sendTopUpFeeToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		fp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoTopUpFee)
	if err := helper.UnpackLog(fp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		fp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := fp.isOldTx(fp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.TopupEvent, event); isOld {
			fp.Logger.Info("Ignoring task to send topup to heimdall as already processed",
				"event", eventName,
				"user", event.User,
				"Fee", event.Fee,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		fp.Logger.Info("✅ sending topup to heimdall",
			"event", eventName,
			"user", event.User,
			"Fee", event.Fee,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// create msg checkpoint ack message
		msg := topupTypes.NewMsgTopup(helper.GetFromAddress(fp.cliCtx), hmTypes.BytesToHeimdallAddress(event.User.Bytes()), sdk.NewIntFromBigInt(event.Fee), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index), vLog.BlockNumber)

		// return broadcast to heimdall
		txRes, err := fp.txBroadcaster.BroadcastToHeimdall(msg, event)
		if err != nil {
			fp.Logger.Error("Error while broadcasting TopupFee msg to heimdall", "msg", msg, "error", err)
			return err
		}

		if txRes.Code != uint32(sdk.CodeOK) {
			fp.Logger.Error("topup tx failed on heimdall", "txHash", txRes.TxHash, "code", txRes.Code)
			return fmt.Errorf("topup tx failed, tx response code: %v", txRes.Code)
		}
	}

	return nil
}
