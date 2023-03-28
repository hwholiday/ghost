package network

import (
	"github.com/gorilla/websocket"
	"net"
)

type Conn interface {
	ConnID() string
	Identity() string
	Group() string
	Conn() net.Conn
	WsConn() *websocket.Conn
	Read() (byt []byte, err error)
	Write(byt []byte) error
	Cache() *Cache
	Close(...[]byte)
	ResetHeartbeat() error
}
