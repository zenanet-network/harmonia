package module

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zenanet-network/harmonia/types/simulation"
)

type StoreDecoderRegistry map[string]func(cdc *codec.Codec, kvA, kvB sdk.KVPair) string

type AppModuleSimulation interface {
	GenerateGenesisState(input *SimulationState)

	ProposalContents(simState SimulationState) []simulation.WeightedProposalContent

	RandomizedParams(r *rand.Rand) []simulation.ParamChange

	RegisterStoreDecoder(StoreDecoderRegistry)

	WeightedOperations(simState SimulationState) []simulation.WeightedOperation
}

type SimulationManager struct {
	Modules       []AppModuleSimulation
	StoreDecoders StoreDecoderRegistry
}

func NewSimulationManager(modules ...AppModuleSimulation) *SimulationManager {
	return &SimulationManager{
		Modules:       modules,
		StoreDecoders: make(StoreDecoderRegistry),
	}
}

func (sm *SimulationManager) GetProposalContents(simState SimulationState) []simulation.WeightedProposalContent {
	wContents := make([]simulation.WeightedProposalContent, 0, len(sm.Modules))
	for _, module := range sm.Modules {
		wContents = append(wContents, module.ProposalContents(simState)...)
	}
	return wContents
}

// RegisterStoreDecoders registers each of the modules' store decoders into a map
func (sm *SimulationManager) RegisterStoreDecoders() {
	for _, module := range sm.Modules {
		module.RegisterStoreDecoder(sm.StoreDecoders)
	}
}

// GenerateGenesisStates generates a randomized GenesisState for each of the
// registered modules
func (sm *SimulationManager) GenerateGenesisStates(simState *SimulationState) {
	for _, module := range sm.Modules {
		module.GenerateGenesisState(simState)
	}
}

// GenerateParamChanges generates randomized contents for creating params change
// proposal transactions
func (sm *SimulationManager) GenerateParamChanges(seed int64) (paramChanges []simulation.ParamChange) {
	r := rand.New(rand.NewSource(seed)) //nolint

	for _, module := range sm.Modules {
		paramChanges = append(paramChanges, module.RandomizedParams(r)...)
	}

	return
}

// WeightedOperations returns all the modules' weighted operations of an application
func (sm *SimulationManager) WeightedOperations(simState SimulationState) []simulation.WeightedOperation {
	wOps := make([]simulation.WeightedOperation, 0, len(sm.Modules))
	for _, module := range sm.Modules {
		wOps = append(wOps, module.WeightedOperations(simState)...)
	}

	return wOps
}

// SimulationState is the input parameters used on each of the module's randomized
// GenesisState generator function
type SimulationState struct {
	AppParams    simulation.AppParams
	Cdc          *codec.Codec                         // application codec
	Rand         *rand.Rand                           // random number
	GenState     map[string]json.RawMessage           // genesis state
	Accounts     []simulation.Account                 // simulation accounts
	GenTimestamp time.Time                            // genesis timestamp
	ParamChanges []simulation.ParamChange             // simulated parameter changes from modules
	Contents     []simulation.WeightedProposalContent // proposal content generator functions with their default weight and app sim key
}
