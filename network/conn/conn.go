package conn

import (
	"net"
)

const (
	IP       = "ip"       // 链接的IP地址
	Identity = "identity" // 链接的UUID
)

type Client interface {
	Conn() net.Conn
	Read() (byt []byte, err error)
	Write(byt []byte) error
	SaveCache(k string, v any)
	GetCache(k string) (v any, ok bool)
	GetCacheString(k string) string
	ResetConnDeadline() error
	Close()
}
