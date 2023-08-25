package main

import (
	"fmt"

	"yellow-jersey/internal/setup"
	"yellow-jersey/pkg/logs"
)

// TODO: Should be in config
const port = 8080 // port of local demo server

func main() {
	srv, err := setup.Setup()
	if err != nil {
		// panic as it's not a recoverable for the API
		panic(err)
	}

	logs.Logger.Info().Msgf("starting api at %d", port)
	logs.Logger.Fatal().Err(srv.Start(fmt.Sprintf(":%d", port))).Msg("")
}
