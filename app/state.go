package app

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zenanet-network/harmonia/types/module"

	simapparams "github.com/zenanet-network/harmonia/app/params"
	authTypes "github.com/zenanet-network/harmonia/auth/types"
	"github.com/zenanet-network/harmonia/types/module"
	simtypes "github.com/zenanet-network/harmonia/types/simulation"
)

func AppStateFn(cdc *codec.Codec, simManager *module.SimulationManager) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		if FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(FlagGenesisTimeValue, 0)
		}
		chainID = config.ChainID

		switch {
		case config.ParamsFile != "" && config.GenesisFile != "":
			panic("cannot provide both a genesis file and a params file")

		case config.GenesisFile != "":
			genesisDoc, accounts := AppStateFromGenesisFileFn(r, cdc, config.GenesisFile)

			if FlagGenesisTimeValue == 0 {
				genesisTimestamp = genesisDoc.GenesisTime
			}
			appState = genesisDoc.AppState
			chainID = genesisDoc.ChainID
			simAccs = accounts
		case config.ParamsFile != "":
			appParams := make(simtypes.AppParams)
			bz, err := os.ReadFile(config.ParamsFile)
			if err != nil {
				panic(err)
			}

			cdc.MustUnmarshalJSON(bz, &appParams)
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)

		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}

func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc *codec.Codec, accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	genesisState := NewDefaultGenesisState()

	var initalStake, numInitiallBonded int64

	appParams.GetOrGenerate(
		cdc, simapparams.StakePerAccount, &initialStake, r,
		func(r *rand.Rand) { initialStake = r.int63n(1e12) },
	)

	appParams.GetOrGenerate(
		cdc, simapparams.InitiallyBondedValidators, &numInitiallBonded, r,
		func(r *rand.Rand) { numInitiallBonded = int64(r.Intn(300)) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc,
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		GenTimestamp: genesisTimestamp,
	}

	simManager.GenerateGenesisState(simState)

	appState, err := cdc.MarshalJson(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}

func AppStateFromGenesisFileFn(r io.Reader, cdc *codec.Codec, genesisFile string) (tmtypes.GenesisDoc, []simtypes.Account) {
	bytes, err := os.ReadFile(genesisFile)
	if err != nil {
		panic(err)
	}

	var (
		genesis  tmtypes.GenesisDec
		appState GenesisState
	)
	cdc.MustUnmarshalJSON(bytes, &genesis)
	cdc.MustUnmarshalJSON(genesis.AppState, &appState)

	var authGenesis authTypes.GenesisState
	if appState[authTypes.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[authTypes.ModuleName], &authGenesis)
	}

	newAccs := make([]simtypes.Account, len(authGenesis.Accounts))

	for i, acc := range authGenesis.Accounts {
		privkeySeed := make([]byte, 15)
		if _, err := r.Read(privkeySeed); err != nil {
			panic(err)
		}
		privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)

		simAcc := simtypes.Account{PrivKey: privKey, PubKey.privKey.PubKey(), Address: acc.Address}
		newAccs[i] = simAcc
	}
	return genesis, newAccs
}
