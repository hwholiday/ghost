package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
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

func NewWsConn(opt ...Option) (Conn, error) {
	var (
		err error
		c   = getWsConn()
	)
	c.opts, err = NewOptions(opt...)
	if err != nil {
		return nil, err
	}
	c.stopChan = make(chan struct{})
	c.cache = NewCache()
	c.writerChan = make(chan []byte, c.opts.readChanLen)
	c.readChan = make(chan []byte, c.opts.witerChanLen)
	c.cache.Save(Identity, c.opts.identity)
	c.cache.Save(Group, c.opts.group)
	c.isOpen.Store(true)
	go c.readChannel()
	go c.witerChannel()
	return c, nil
}
func (c *wsConn) WsConn() *websocket.Conn {
	return c.WsConn()
}

func (c *wsConn) Identity() string {
	return c.cache.Get(Identity).String()
}
func (c *wsConn) Group() string {
	return c.cache.Get(Group).String()
}

func (c *wsConn) RemoteAddr() string {
	return c.opts.wsConn.RemoteAddr().String()
}

func (c *wsConn) LocalAddr() string {
	return c.opts.wsConn.LocalAddr().String()
}

func (c *wsConn) Cache() *Cache {
	return c.cache
}

func (c *wsConn) Close() {
	if c.isOpen.Load() {
		c.isOpen.Store(false)
		_ = c.opts.wsConn.Close()
		putWsConn(c)
		c.stopChan <- struct{}{}
		close(c.stopChan)
		close(c.readChan)
		close(c.writerChan)
	}
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
func (c *wsConn) ResetConnDeadline() error {
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
		messageType, reader, err := c.opts.wsConn.NextReader()
		if err != nil {
			log.Printf("[Dove] readChannel NextReader Close wsConn identity : %s , err: %s ", c.opts.identity, err.Error())
			c.Close()
			return
		}
		if messageType != websocket.BinaryMessage {
			log.Printf("[Dove] readChannel wsConn only support BinaryMessage identity : %s , err: %s ", c.opts.identity)
			c.Close()
			return
		}
		byt, err := c.read(reader)
		if err != nil {
			log.Printf("[Dove] readChannel read Close wsConn identity : %s , err: %s ", c.opts.identity, err.Error())
			c.Close()
			return
		}
		if c.opts.autoResetConnDeadline {
			_ = c.ResetConnDeadline()
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
				log.Printf("[Dove] witerChannel wsConn identity : %s , err: %s ", c.opts.identity, err.Error())
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
