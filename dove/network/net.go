package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

var (
	_               Conn = (*conn)(nil)
	AlreadyCloseErr      = errors.New("conn already close")
	MayBeCloseErr        = errors.New("conn may be closed or request data format error")
)

// conn length + data 模式
type conn struct {
	opts       *options
	readWriter *bufio.ReadWriter
	cache      *Cache
	stopChan   chan struct{}
	writerChan chan []byte
	readChan   chan []byte
	isOpen     atomic.Bool
}

func NewConnWithOpt(opt *options) Conn {
	var (
		c = getConn()
	)
	c.opts = opt
	c.readWriter = bufio.NewReadWriter(
		bufio.NewReaderSize(c.opts.conn, c.opts.readBufferSize),
		bufio.NewWriterSize(c.opts.conn, c.opts.witerBufferSize))
	c.stopChan = make(chan struct{})
	c.cache = NewCache()
	c.writerChan = make(chan []byte, c.opts.readChanLen)
	c.readChan = make(chan []byte, c.opts.witerChanLen)
	c.isOpen.Store(true)
	go c.readChannel()
	go c.witerChannel()
	return c
}

func NewConn(opt ...Option) (Conn, error) {
	opts, err := NewOptions(opt...)
	if err != nil {
		return nil, err
	}
	return NewConnWithOpt(opts), nil
}

func (c *conn) ConnID() string {
	return c.opts.connID
}
func (c *conn) WsConn() *websocket.Conn {
	return nil
}
func (c *conn) logger() zerolog.Logger {
	return log.With().Str("connID", c.opts.connID).Str("identity", c.opts.identity).Str("group", c.opts.group).Logger()
}
func (c *conn) Identity() string {
	return c.opts.identity
}
func (c *conn) Group() string {
	return c.opts.group
}

func (c *conn) Cache() *Cache {
	return c.cache
}

func (c *conn) Close(byt ...[]byte) {
	if c.isOpen.Load() {
		c.isOpen.Store(false)
		c.stopChan <- struct{}{}
		close(c.stopChan)
		close(c.readChan)
		close(c.writerChan)
		c.witerClose(byt...)
		_ = c.opts.conn.Close()
		putConn(c)
	}
}

func (c *conn) witerClose(byt ...[]byte) {
	if len(byt) <= 0 {
		return
	}
	_ = c.witer(byt[0])
}

func (c *conn) Read() (byt []byte, err error) {
	select {
	case byt = <-c.readChan:
		return byt, nil
	case <-c.stopChan:
		return nil, AlreadyCloseErr
	}
}

func (c *conn) Conn() net.Conn {
	return c.opts.conn
}
func (c *conn) ResetHeartbeat() error {
	return c.opts.conn.SetReadDeadline(time.Now().Add(c.opts.heartbeatInterval))
}

func (c *conn) Write(byt []byte) error {
	select {
	case c.writerChan <- byt:
	case <-c.stopChan:
		return AlreadyCloseErr
	}
	return nil
}

func (c *conn) readChannel() {
	for {
		byt, err := c.read()
		if err != nil {
			if !errors.Is(err, AlreadyCloseErr) && !strings.Contains(err.Error(), "use of closed network connection") {
				logger := c.logger()
				logger.Error().Err(err).Msg("[Dove] conn readChannel failed")
			}
			c.Close()
			return
		}
		if c.opts.autoHeartbeat {
			_ = c.ResetHeartbeat()
		}
		select {
		case c.readChan <- byt:
		case <-c.stopChan:
			return
		}
	}
}

func (c *conn) witerChannel() {
	for {
		select {
		case byt := <-c.writerChan:
			if err := c.witer(byt); err != nil {
				logger := c.logger()
				logger.Error().Err(err).Msg("[Dove] conn witerChannel failed")
			}
		case <-c.stopChan:
			return
		}
	}
}

func (c *conn) read() ([]byte, error) {
	lengthByte, err := c.readWriter.Reader.Peek(c.opts.length)
	if err != nil {
		return nil, err
	}
	var length int32
	if err = binary.Read(bytes.NewReader(lengthByte), c.opts.endian, &length); err != nil {
		return nil, MayBeCloseErr
	}
	if c.readWriter.Reader.Buffered() < c.opts.length+int(length) {
		return nil, errors.New("the corresponding data cannot be read")
	}
	pack := make([]byte, c.opts.length+int(length))

	if _, err = c.readWriter.Reader.Read(pack); err != nil {
		return nil, err
	}
	return pack[c.opts.length:], err
}

func (c *conn) witer(byt []byte) error {
	var (
		length = int32(len(byt))
	)
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, length); err != nil {
		return err
	}
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, byt); err != nil {
		return err
	}
	return c.readWriter.Writer.Flush()
}
