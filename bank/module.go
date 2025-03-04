package bank

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/zenanet-network/harmonia/bank/types"
	"github.com/zenanet-network/harmonia/checkpoint/types"
	"github.com/zenanet-network/harmonia/supply/types"

	bankCli "github.com/zenanet-network/harmonia/bank/client/cli"
	bankRest "github.com/zenanet-network/harmonia/bank/client/rest"
	"github.com/zenanet-network/harmonia/bank/simulation"
	"github.com/zenanet-network/harmonia/bank/types"
	"github.com/zenanet-network/harmonia/helper"
	hmModule "github.com/zenanet-network/harmonia/types/module"
	simTypes "github.com/zenanet-network/harmonia/types/simulation"
)

var (
	_ module.AppModule = AppModule{}

	_ module.AppModuleBasic        = AppModuleBasic{}
	_ hmModule.HarmoniaModuleBasic = AppModule{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustUnmarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}

	return types.ValidateGenesis(data)
}

func (AppModuleBasic) VerifyGenesis(bz map[string]json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	bankRest.ReigsterRoutes(ctx, rtr)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return bankCli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return bankCli.GetQueryCmd(cdc)
}

//=================================================================

type AppModule struct {
	AppModuleBasic

	keeper         Keeper
	contractCaller helper.IContractCaller
}

func NewAppModule(keeper Keeper, contractCaller helper.IcontractCaller) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		contractCaller: contractCaller,
	}
}

func (AppModule) Name() string {
	return types.ModuleName
}

func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() string {
	return types.RouterKey
}

func (am AppModule) NewHanlder() sdk.Handler {
	return NewHandler(am.keeper, am.contractCaller)
}

func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the auth module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the auth module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState creates a randomized GenState of the chainManager module
func (AppModule) GenerateGenesisState(simState *hmModule.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions.
func (AppModule) ProposalContents(simState hmModule.SimulationState) []simTypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simTypes.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for chainmanager module's types
func (AppModule) RegisterStoreDecoder(sdr hmModule.StoreDecoderRegistry) {
}

// WeightedOperations doesn't return any chainmanager module operation.
func (AppModule) WeightedOperations(_ hmModule.SimulationState) []simTypes.WeightedOperation {
	return nil
}
