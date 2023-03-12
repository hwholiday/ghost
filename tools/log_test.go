package tools

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func Test_Log(t *testing.T) {
	LogGlobalConf(true)
	log.Info().Any("test", 11).Msg("data")
	LogGlobalConf(false)
	log.Info().Any("test", 11).Msg("data")
}
