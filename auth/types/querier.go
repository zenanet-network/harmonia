package types

import (
	"github.com/zenanet-network/harmonia/types"
)

// query endpoints supported by the auth Querier
const (
	QueryParams  = "params"
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address types.HarmoniaAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr types.HarmoniaAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
