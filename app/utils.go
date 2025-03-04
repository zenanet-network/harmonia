package app

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/zenanet-network/harmonia/types/module"
	simTypes "github.com/zenanet-network/harmonia/types/simulation"
)

const (
	SimAppChainID = "simulation-app"
)

func SetupSimulation(dirPreFix, dbName string) (simTypes.Config, dbm.DB, string, log.Logger, bool, error) {
	config := NewConfigFromFlags()
	config.ChainID = SimAppChainID

	var lggr log.Logger
	if FlagVerboseValue {
		lggr = log.TestingLogger()
	} else {
		lggr = log.NewNopLogger()
	}

	dir, err := os.MkdirTemp("", dirPrefix)
	if err != nil {
		return simTypes.Config{}, nil, "", nil, false, err
	}

	db, err := sdk.NewLevelDB(dbName, dir)
	if err != nil {
		return simTypes.Config{}, nil, "", nil, false, err
	}

	return config, db, dir, lggr, false, nil
}

func SimulationOperations(app App, cdc *codec.Codec, config simTypes.Config) []simTypes.WeightedOperation {
	simState := module.SimulationState{
		AppParams: make(simTypes.AppParams),
		Cdc:       cdc,
	}

	if config.ParamsFile != "" {
		bz, err := os.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}
		app.Codec().MustUnmarshalJSON(bz, &simState.AppParams)
	}

	simState.ParamsChanges = app.SimulationManager().GenerateParamChanges(config.Seed)
	simState.Contents = app.SimulationManager().GetProposalContents(simState)

	return app.SimulationManager().WeightedOperation()

}

func CheckExportSimulation(app App, config simTypes.Config, params simTypes.Params) error {
	if config.ExportStatePath != "" {
		fmt.Println("exporting app state...")

		appState, _, err := app.ExportAppStateAndValidators()
		if err != nil {
			return nil
		}

		if err = os.WriteFile(config.ExportStatePath, appState, 0644); err != nil {
			return err
		}
	}

	if config.ExportParamsPath != "" {
		fmt.Println("exporting simulation params...")

		paramsBz, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(params, "", " ")
		if err != nil {
			return err
		}

		if err = os.WriteFile(config.ExportParamsPath, paramsBz, 0644); err != nil {
			return err
		}
	}
	return nil
}

func PrintStats(db dbm.DB) {
	fmt.Println("\nLevelDB Stats")
	fmt.Println(db.Stats()["leveldb.stats"])
	fmt.Println("LevelDB cached block size", db.Stats()["leveldb.cachedblock"])
}

func GetSimulationLog(storeName string, sdr module.StoreDecoderRegistry, cdc *codec.Codec, kvAs, kvBs []sdk.KVPair) (log string) {
	for i := 0; i < len(kvAs); i++ {
		if len(kvAs[i].Value) == 0 && len(kvBs[i].Value) == 0 {
			continue
		}

		decoder, ok := sdr[storeName]
		if ok {
			log += decoder(cdc, kvAs[i], kvBs[i])
		} else {
			log += fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", kvAs[i].Key, kvAs[i].Value, kvBs[i].Key, kvBs[i].Value)
		}
	}

	return
}
