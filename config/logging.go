package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitLogger() zerolog.Logger {
	if viper.GetBool(LogJson) {
		log.Logger = log.Output(zerolog.New(os.Stdout))
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if viper.GetBool(LogDebug) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if viper.GetBool(LogTrace) {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	return log.With().Logger()
}

func WithComponent(logger zerolog.Logger, name string) zerolog.Logger {
	return logger.With().Str("component", name).Logger()
}
