package main

import (
	"github.com/rs/zerolog/log"

	"github.com/chenyanchen/redis-cleaner/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		log.Fatal().Err(err).Msg("Execute")
	}
}
