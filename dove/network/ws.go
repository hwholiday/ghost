package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"sync/atomic"
	"time"
)

var (
	_ Conn = (*wsConn)(nil)
)

// wsConn length + data 模式
type wsConn struct {
	opts       *options
	cache      *Cache
	stopChan   chan struct{}
	writerChan chan []byte
	readChan   chan []byte
	isOpen     atomic.Bool
}

func (c *wsConn) ConnID() string {
	return c.opts.connID
}
func (c *wsConn) logger() zerolog.Logger {
	return log.With().Str("conn-id", c.opts.connID).Str("identity", c.opts.identity).Str("group", c.opts.group).Logger()
}

func NewWsConnWithOpt(opt *options) Conn {
	var (
		c = getWsConn()
	)
	c.opts = opt
	c.stopChan = make(chan struct{})
	c.cache = NewCache()
	c.writerChan = make(chan []byte, c.opts.readChanLen)
	c.readChan = make(chan []byte, c.opts.witerChanLen)
	c.isOpen.Store(true)
	c.saveCacheByOpts(c.opts)
	go c.readChannel()
	go c.witerChannel()
	return c
}

func (c *wsConn) saveCacheByOpts(opts *options) {
	for key, val := range c.opts.cache {
		c.cache.Save(key, val)
	}
}

func NewWsConn(opt ...Option) (Conn, error) {
	opts, err := NewOptions(opt...)
	if err != nil {
		return nil, err
	}
	return NewWsConnWithOpt(opts), nil
}
func (c *wsConn) WsConn() *websocket.Conn {
	return c.WsConn()
}

func (c *wsConn) Identity() string {
	return c.opts.identity
}
func (c *wsConn) Group() string {
	return c.opts.group
}

func (c *wsConn) Cache() *Cache {
	return c.cache
}

func (c *wsConn) Close(byt ...[]byte) {
	if c.isOpen.Load() {
		c.isOpen.Store(false)
		c.stopChan <- struct{}{}
		close(c.stopChan)
		close(c.readChan)
		close(c.writerChan)
		c.witerClose(byt...)
		_ = c.opts.wsConn.Close()
		putWsConn(c)
	}
}
func (c *wsConn) witerClose(byt ...[]byte) {
	if len(byt) <= 0 {
		return
	}
	_ = c.witer(byt[0])
}

func (c *wsConn) Read() (byt []byte, err error) {
	select {
	case byt = <-c.readChan:
		return byt, nil
	case <-c.stopChan:
		return nil, AlreadyCloseErr
	}
}

func (c *wsConn) Conn() net.Conn {
	return nil
}
func (c *wsConn) ResetHeartbeat() error {
	return c.opts.wsConn.SetReadDeadline(time.Now().Add(c.opts.heartbeatInterval))
}

func (c *wsConn) Write(byt []byte) error {
	select {
	case c.writerChan <- byt:
	case <-c.stopChan:
		return AlreadyCloseErr
	}
	return nil
}

func (c *wsConn) readChannel() {
	for {
		logger := c.logger()
		messageType, reader, err := c.opts.wsConn.NextReader()
		if err != nil {
			logger.Error().Err(err).Msg("[Dove] ws conn readChannel NextReader failed")
			c.Close()
			return
		}
		if messageType != websocket.BinaryMessage {
			logger.Error().Msg("[Dove] ws conn readChannel only support BinaryMessage failed")

			log.Printf("[Dove] readChannel wsConn only support BinaryMessage id %s , identity : %s , err: %s ", c.opts.connID, c.opts.identity, err.Error())
			c.Close()
			return
		}
		byt, err := c.read(reader)
		if err != nil {
			logger.Error().Msg("[Dove] ws conn readChannel read failed")
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

func (c *wsConn) witerChannel() {
	for {
		select {
		case byt := <-c.writerChan:
			if err := c.witer(byt); err != nil {
				logger := c.logger()
				logger.Error().Err(err).Msg("[Dove] ws conn witerChannel failed")
			}
		case <-c.stopChan:
			return
		}
	}
}

func (c *wsConn) read(input io.Reader) ([]byte, error) {
	rd := bufio.NewReader(input)
	lengthByte, err := rd.Peek(c.opts.length)
	if err != nil {
		return nil, err
	}
	var length int32
	if err = binary.Read(bytes.NewReader(lengthByte), c.opts.endian, &length); err != nil {
		return nil, MayBeCloseErr
	}
	if rd.Buffered() < c.opts.length+int(length) {
		return nil, errors.New("the corresponding data cannot be read")
	}
	pack := make([]byte, c.opts.length+int(length))

	if _, err = rd.Read(pack); err != nil {
		return nil, err
	}
	return pack[c.opts.length:], err
}

func (c *wsConn) witer(byt []byte) error {
	var (
		length = int32(len(byt))
		pkg    = new(bytes.Buffer)
	)

	if err := binary.Write(pkg, c.opts.endian, length); err != nil {
		return err
	}
	if err := binary.Write(pkg, c.opts.endian, byt); err != nil {
		return err
	}
	writer, err := c.opts.wsConn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		c.Close()
		return err
	}
	if _, err = writer.Write(pkg.Bytes()); err != nil {
		return err
	}
	_ = writer.Close()
	return nil
}
