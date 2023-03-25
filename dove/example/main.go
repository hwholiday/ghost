package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/hwholiday/ghost/dove"
	api "github.com/hwholiday/ghost/dove/api/dove"
	"github.com/hwholiday/ghost/dove/network"
	"github.com/rs/zerolog/log"
	"net/http"
)

// golang 通过Reader.Peek读取 Uint8Array 数据
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var client dove.Dove

func main() {
	dove.SetMode(dove.DebugMode)
	dove.SetConnMax(100)
	dove.DefaultWsPort = ":10888"
	client = dove.NewDove()
	client.RegisterHandleFunc(dove.DefaultConnAcceptCrcId, func(cli network.Conn, data *api.Dove) {
		log.Info().Str("Identity", cli.Cache().Get(network.Identity).String()).Msg("设备上线")
		for _, v := range client.Manage().FindConnByGroup("user-001") {
			//v.Write([]byte("1111"))
			log.Info().Str("group", "user-001").Str("identity", v.Identity()).Msg("FindConnByGroup succeed")
		}
	})
	client.RegisterHandleFunc(dove.DefaultConnCloseCrcId, func(cli network.Conn, data *api.Dove) {
		log.Info().Str("Identity", cli.Cache().Get(network.Identity).String()).Msg("设备离线")
		for _, v := range client.Manage().FindConnByGroup("user-001") {
			log.Info().Str("group", "user-001").Str("identity", v.Identity()).Msg("FindConnByGroup succeed")
		}
	})
	client.RegisterHandleFunc(3, func(cli network.Conn, data *api.Dove) {
		logger := log.With().Str("Identity", cli.Cache().Get(network.Identity).String()).Interface("data", data).Logger()
		logger.Info().Msg("func id 3")
		res, err := dove.NewBuild().BuildMetadata(data.GetMetadata().GetCrcId(), data.GetMetadata().GetAckId()).BuildDoveBodyOk().Result()
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
		_ = cli.Write(res)
	})
	Listen()
}

func Listen() {
	http.HandleFunc("/socket.io", HandleSocketForWs)
	http.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		obj := client.Manage().GetMapStatus()
		bytes, _ := json.Marshal(obj)
		_, _ = writer.Write(bytes)
	})
	log.Info().Str("addr", ":10888").Msg("example service start succeed")
	err := http.ListenAndServe(dove.DefaultWsPort, nil)
	if err != nil {
		panic(err)
	}
}

func HandleSocketForWs(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		wsConn *websocket.Conn
	)
	if wsConn, err = wsUpgrader.Upgrade(res, req, nil); err != nil {
		log.Error().Err(err).Msg("Upgrade failed")
		return
	}
	err = client.Accept(network.WithWsConn(wsConn), network.WithGroup("user-001"), network.WithLength(4))
	if err != nil {
		log.Error().Err(err).Msg("Accept failed")
	}
}
