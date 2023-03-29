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

const (
	MessageType_AckIDRemoteLogin int32 = 3
	MessageType_CrcIDHeartbeat   int32 = 10
	MessageType_AckIDHeartbeat   int32 = 11
	MessageType_CrcIDSingleChat  int32 = 12
	MessageType_AckIDSingleChat  int32 = 13
)

func main() {
	dove.SetMode(dove.DebugMode)
	dove.SetConnMax(100)
	dove.DefaultWsPort = ":10888"
	client = dove.NewDove()
	client.EventNotify(func(protocol api.EventType, cli network.Conn) {
		logger := dove.Logger(cli, nil)
		switch protocol {
		case api.EventType_ConnAccept:
			logger.Debug().Msg("设备上线")
		case api.EventType_ConnClose:
			logger.Debug().Msg("设备离线")
		}
	})
	client.RegisterHandleFunc(MessageType_CrcIDHeartbeat, func(cli network.Conn, data *api.Dove) {
		logger := dove.Logger(cli, data)
		logger.Debug().Msg("接收到心跳")
		resData, _ := dove.NewDoveRes().Metadata(data.GetMetadata().GetCrcId(), MessageType_AckIDHeartbeat).BodyOk().Result()
		_ = cli.Write(resData)
	})
	client.RegisterHandleFunc(MessageType_CrcIDSingleChat, func(cli network.Conn, data *api.Dove) {
		logger := dove.Logger(cli, data)
		logger.Debug().Any("data", data).Msg("接收到消息")
		resData, _ := dove.NewDoveRes().Metadata(data.GetMetadata().GetCrcId(), MessageType_AckIDSingleChat).BodyOk().Result()
		_ = cli.Write(resData)
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
	if err = client.CanAccept("room-001-user-666"); err != nil {
		switch err {
		case dove.ErrExceedsLengthLimit:
			//连接数已经满了不能再连接
			_ = wsConn.Close()
			return
		case dove.ErrIdentityAlreadyExists:
			//可能存在异地登录，需要踢掉之前的连接
			if conn, ok := client.Manage().GetConn("room-001-user-666"); ok {
				client.Manage().Del(conn.Identity())
				log.Warn().Str("old", conn.ConnID()).Msg("可能存在异地登录，需要踢掉之前的连接")
				resData, _ := dove.NewDoveRes().MetadataAckId(MessageType_AckIDRemoteLogin).BodyOk().Result()
				conn.Close(resData)
			}
		default:
			log.Warn().Msg("unknown error type")
		}
	}
	err = client.Accept(network.WithWsConn(wsConn), network.WithIdentity("room-001-user-666"), network.WithGroup("room-001"), network.WithLength(4))
	if err != nil {
		log.Error().Err(err).Msg("accept failed")
	}
}
