package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zenanet-network/harmonia/helper"
	"github.com/zenanet-network/harmonia/types"
)

// MsgEventRecord - state msg
type MsgEventRecord struct {
	From            types.HarmoniaAddress `json:"from"`
	TxHash          types.HarmoniaHash    `json:"tx_hash"`
	LogIndex        uint64                `json:"log_index"`
	BlockNumber     uint64                `json:"block_number"`
	ContractAddress types.HarmoniaAddress `json:"contract_address"`
	Data            types.HexBytes        `json:"data"`
	ID              uint64                `json:"id"`
	ChainID         string                `json:"bor_chain_id"`
}

var _ sdk.Msg = MsgEventRecord{}

// NewMsgEventRecord - construct state msg
func NewMsgEventRecord(
	from types.HarmoniaAddress,
	txHash types.HarmoniaHash,
	logIndex uint64,
	blockNumber uint64,
	id uint64,
	contractAddress types.HarmoniaAddress,
	data types.HexBytes,
	chainID string,

) MsgEventRecord {
	return MsgEventRecord{
		From:            from,
		TxHash:          txHash,
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		ID:              id,
		ContractAddress: contractAddress,
		Data:            data,
		ChainID:         chainID,
	}
}

// Route Implements Msg.
func (msg MsgEventRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgEventRecord) Type() string { return "event-record" }

// ValidateBasic Implements Msg.
func (msg MsgEventRecord) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.TxHash.Empty() {
		return sdk.ErrInvalidAddress("missing tx hash")
	}

	// DO NOT REMOVE THIS CHANGE
	if len(msg.Data) > helper.LegacyMaxStateSyncSize {
		return ErrSizeExceed(sdk.CodespaceType(fmt.Sprintf("length is larger than %d bytes", helper.LegacyMaxStateSyncSize)))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEventRecord) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgEventRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HarmoniaAddressToAccAddress(msg.From)}
}

// GetTxHash Returns tx hash
func (msg MsgEventRecord) GetTxHash() types.HarmoniaHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgEventRecord) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgEventRecord) GetSideSignBytes() []byte {
	return nil
}
