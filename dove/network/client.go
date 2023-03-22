package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
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
	tmp        any
}

func NewConn(opt ...Option) (Conn, error) {
	var (
		err error
		c   = getConn()
	)
	c.opts, err = newOptions(opt...)
	if err != nil {
		return nil, err
	}
	c.readWriter = bufio.NewReadWriter(
		bufio.NewReaderSize(c.opts.conn, c.opts.readBufferSize),
		bufio.NewWriterSize(c.opts.conn, c.opts.witerBufferSize))
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

func (c *conn) Identity() string {
	return c.cache.Get(Identity).String()
}
func (c *conn) Group() string {
	return c.cache.Get(Group).String()
}

func (c *conn) RemoteAddr() string {
	return c.opts.conn.RemoteAddr().String()
}

func (c *conn) LocalAddr() string {
	return c.opts.conn.LocalAddr().String()
}

func (c *conn) Cache() *Cache {
	return c.cache
}

func (c *conn) Close() {
	if c.isOpen.Load() {
		c.isOpen.Store(false)
		_ = c.opts.conn.Close()
		putConn(c)
		c.stopChan <- struct{}{}
		close(c.stopChan)
		close(c.readChan)
		close(c.writerChan)
	}
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
func (c *conn) ResetConnDeadline() error {
	return c.opts.conn.SetDeadline(time.Now().Add(c.opts.heartbeatInterval))
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
				log.Printf("[Dove] readChannel Close conn identity : %s , err: %s ", c.opts.identity, err.Error())
			}
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

func (c *conn) witerChannel() {
	for {
		select {
		case byt := <-c.writerChan:
			if err := c.witer(byt); err != nil {
				log.Printf("[Dove] witerChannel identity : %s , err: %s ", c.opts.identity, err.Error())
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
	var length int64
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
		length = int64(len(byt))
	)
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, length); err != nil {
		return err
	}
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, byt); err != nil {
		return err
	}
	return c.readWriter.Writer.Flush()
}
