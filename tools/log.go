package tools

import (
	"github.com/rs/zerolog"
	"time"
)

func LogGlobalConf(development bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if development {
		zerolog.TimeFieldFormat = time.DateTime
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
