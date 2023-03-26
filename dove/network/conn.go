package network

import (
	"github.com/gorilla/websocket"
	"net"
)

const (
	Identity = "identity" // 链接的ID
	Group    = "group"    // 链接的ID分组
)

type Conn interface {
	ConnId() string
	Identity() string
	Group() string
	RemoteAddr() string
	LocalAddr() string
	Conn() net.Conn
	WsConn() *websocket.Conn
	Read() (byt []byte, err error)
	Write(byt []byte) error
	WiterNoChan(byt []byte) error
	Cache() *Cache
	Close()
	ResetConnDeadline() error
}
