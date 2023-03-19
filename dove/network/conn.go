package network

import (
	"net"
)

const (
	Identity = "identity" // 链接的ID
	Group    = "group"    // 链接的ID分组
)

type Conn interface {
	Identity() string
	Group() string
	RemoteAddr() string
	LocalAddr() string
	Conn() net.Conn
	Read() (byt []byte, err error)
	Write(byt []byte) error
	Cache() *Cache
	Close()
	ResetConnDeadline() error
}
