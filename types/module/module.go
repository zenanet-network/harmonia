package module

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/zenanet-network/harmonia/types"
)

type HarmoniaModuleBasic interface {
	module.AppModuleBasic
	VerifyGensis(map[string]json.RawMessage) error
}

type SideModule interface {
	NewSideTxHandler() types.SideTxHandler
	NewPostTxHandler() types.PostTxHandler
}

type ModuleGenesisData struct {
	Path string

	Data json.RawMessage

	NextKey []byte
}

type StreamedGenesisExporter interface {
	ExportPartialGenesis(ctx sdk.Context) (json.RawMessage, error)

	NextGenesisData(ctx sdk.Context, nextKey []byte, max int) (*ModuleGenesisData, error)
}
