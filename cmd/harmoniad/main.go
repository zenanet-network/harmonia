package main

import (
	"context"

	"github.com/zenanet-network/harmonia/cmd/harmoniad/service"
)

func main() {
	service.NewHarmoniadService(context.Background(), nil)
}
