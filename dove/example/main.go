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
	client = dove.NewDove()
	client.RegisterHandleFunc(dove.DefaultConnAcceptCrcId, func(cli network.Conn, data *api.Dove) {
		log.Info().Str("Identity", cli.Cache().Get(network.Identity).String()).Msg("设备上线")
		for _, v := range client.Manage().FindConnByGroup("user-001") {
			//_ = v.Write([]byte("succeed"))
			log.Info().Str("group", "user-001").Str("identity", v.Identity()).Msg("FindConnByGroup succeed")
		}
	})
	client.RegisterHandleFunc(dove.DefaultConnCloseCrcId, func(cli network.Conn, data *api.Dove) {
		log.Info().Str("Identity", cli.Cache().Get(network.Identity).String()).Msg("设备离线")
		for _, v := range client.Manage().FindConnByGroup("user-001") {
			log.Info().Str("group", "user-001").Str("identity", v.Identity()).Msg("FindConnByGroup succeed")
		}
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
	log.Info().Str("addr", dove.DefaultWsPort).Msg("example service start succeed")
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
	err = client.Accept(network.WithConn(wsConn.UnderlyingConn()), network.WithGroup("user-001"))
	if err != nil {
		log.Error().Err(err).Msg("Accept failed")
	}
}
