package utils

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func Test_Log(t *testing.T) {
	SetUpGlobalZeroLogConf(true)
	log.Info().Any("test", 11).Msg("data")
	SetUpGlobalZeroLogConf(false)
	log.Info().Any("test", 11).Msg("data")
}
