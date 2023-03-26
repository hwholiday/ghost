package network

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type Option func(*options)

type options struct {
	identity              string
	group                 string
	conn                  net.Conn
	wsConn                *websocket.Conn
	useBigEndian          bool
	endian                binary.ByteOrder
	length                int
	readBufferSize        int
	witerBufferSize       int
	witerChanLen          int
	readChanLen           int
	heartbeatInterval     time.Duration
	autoResetConnDeadline bool
}

func WithConn(conn net.Conn) Option {
	return func(o *options) {
		o.conn = conn
	}
}
func WithWsConn(wsConn *websocket.Conn) Option {
	return func(o *options) {
		o.wsConn = wsConn
	}
}
func WithUseBigEndian(useBigEndian bool) Option {
	return func(o *options) {
		o.useBigEndian = useBigEndian
	}
}

func WithAutoResetConnDeadline(auto bool) Option {
	return func(o *options) {
		o.autoResetConnDeadline = auto
	}
}

func WithReadChanLen(len int) Option {
	return func(o *options) {
		o.readChanLen = len
	}
}

func WithWiterChanLen(len int) Option {
	return func(o *options) {
		o.witerChanLen = len
	}
}

func WithReadBufferSize(size int) Option {
	return func(o *options) {
		o.readBufferSize = size
	}
}

func WithWiterBufferSize(size int) Option {
	return func(o *options) {
		o.witerBufferSize = size
	}
}

func WithHeartbeatInterval(t time.Duration) Option {
	return func(o *options) {
		o.heartbeatInterval = t
	}
}

func WithLength(length int) Option {
	return func(o *options) {
		o.length = length
	}
}

func WithIdentity(identity string) Option {
	return func(o *options) {
		o.identity = identity
	}
}
func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}

func (o *options) HasWsConn() bool {
	return o.wsConn != nil
}

func (o *options) GetIdentity() string {
	return o.identity
}

func (o *options) HasConn() bool {
	return o.conn != nil
}
func NewOptions(opts ...Option) (*options, error) {
	o := &options{
		identity:              uuid.New().String(),
		witerBufferSize:       4096,
		readBufferSize:        4096,
		witerChanLen:          1,
		readChanLen:           1,
		length:                4,
		heartbeatInterval:     time.Second * 30,
		useBigEndian:          true,
		autoResetConnDeadline: true,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.conn == nil && o.wsConn == nil {
		return nil, errors.New("conn is nil")
	}
	if o.conn != nil && o.wsConn != nil {
		_ = o.conn.Close()
		_ = o.wsConn.Close()
		return nil, errors.New("only support conn or wsConn")
	}
	if o.useBigEndian {
		o.endian = binary.BigEndian
	} else {
		o.endian = binary.LittleEndian
	}
	return o, nil
}
