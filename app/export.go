package app

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

func (a *HarmoniaApp) ExportAppStateAndValidators() (
	appState json.RawMessage,
	validators []tmTypes.GenesisValidator,
	err error,
) {
	ctx := a.NewContext(true, abci.Header{Height: a.LastBlockHeight()})
	result := a.mm.ExportGenesis(ctx)

	appState, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)

	return appState, validators, err
}
