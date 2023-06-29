package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/chenyanchen/redis-cleaner/cmd"
)

func main() {
	// Use console logger.
	options := func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
	}
	log.Logger = log.Output(zerolog.NewConsoleWriter(options))

	// New command and execute.
	if err := cmd.New().Execute(); err != nil {
		log.Fatal().Err(err).Msg("Execute")
	}
}
