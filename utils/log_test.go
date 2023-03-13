package utils

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func Test_Log(t *testing.T) {
	SetUpGlobalLogConf(true)
	log.Info().Any("test", 11).Msg("data")
	SetUpGlobalLogConf(false)
	log.Info().Any("test", 11).Msg("data")
}
