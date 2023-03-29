package dove

import (
	"github.com/google/uuid"
	api "github.com/hwholiday/ghost/dove/api/dove"
	"github.com/hwholiday/ghost/dove/network"
	"github.com/hwholiday/ghost/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DefaultWsPort               = ":8081"
	DefaultConnMax        int64 = 10000
	DefaultDoveBodyCodeOK int32 = 200
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

var doveMode = ReleaseMode

func SetConnMax(value int64) {
	DefaultConnMax = value
}
func SetMode(value string) {
	switch value {
	case DebugMode:
		doveMode = DebugMode
	case ReleaseMode:
		doveMode = ReleaseMode
	default:
		doveMode = ReleaseMode
	}
}

func ModeName() string {
	return doveMode
}

func setup() {
	utils.SetUpGlobalZeroLogConf(doveMode == DebugMode)
	log.Info().Str("dove run mode :", ModeName()).Send()
}

func Logger(cli network.Conn, data *api.Dove) zerolog.Logger {
	traceId := data.GetMetadata().GetSeq()
	if traceId == "" {
		traceId = uuid.NewString()
	}
	logger := log.With().Str("trace-id", traceId).Logger()
	if cli == nil {
		return logger
	}
	return logger.With().Str("conn-id", cli.ConnID()).Str("identity", cli.Identity()).Str("group", cli.Group()).Logger()
}
