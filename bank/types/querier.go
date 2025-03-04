package types

import (
	hmTyps "github.com/zenanet-network/harmonia/types"
)

const (
	QueryBalance = "balances"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address hmTyps.HarmoniaAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr hmTyps.HarmoniaAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
