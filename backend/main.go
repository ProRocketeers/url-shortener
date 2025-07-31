package main

import (
	_ "github.com/ProRocketeers/url-shortener/docs"
	"github.com/ProRocketeers/url-shortener/infrastructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "dev"

func main() {
	// logger setup is kinda weird, before we parse the config, we don't know logger level, so let's assume it's JSON output, level info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	config, err := infrastructure.ParseServerConfig(Version)
	if err != nil {
		// Fatal already calls `os.Exit(1)`
		log.Fatal().Err(err).Msg("error while parsing server config")
	}

	if err := infrastructure.RunServerGracefully(config); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
