package conn

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	_               Client = (*Conn)(nil)
	AlreadyCloseErr        = errors.New("conn already close")
)

// Conn length + data 模式
type Conn struct {
	opts       *options
	readWriter *bufio.ReadWriter
	cache      sync.Map
	stopChan   chan bool
	writerChan chan []byte
	readChan   chan []byte
	once       sync.Once
}

func NewConnClient(opt ...Option) (Client, error) {
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
	c.stopChan = make(chan bool)
	c.writerChan = make(chan []byte, c.opts.readChanLen)
	c.readChan = make(chan []byte, c.opts.witerChanLen)
	c.SaveCache(IP, c.opts.conn.RemoteAddr().String())
	c.SaveCache(Identity, c.opts.id)
	go c.readChannel()
	go c.witerChannel()
	return c, nil
}
func (c *Conn) Close() {
	c.once.Do(func() {
		c.stopChan <- true
		close(c.stopChan)
		close(c.readChan)
		close(c.writerChan)
		_ = c.opts.conn.Close()
		putConn(c)
	})
}

func (c *Conn) Read() (byt []byte, err error) {
	select {
	case byt = <-c.readChan:
		return byt, nil
	case <-c.stopChan:
		return nil, AlreadyCloseErr
	}
}

func (c *Conn) SaveCache(k string, v any) {
	c.cache.Store(k, v)
}

func (c *Conn) GetCache(k string) (v any, ok bool) {
	return c.cache.Load(k)
}

func (c *Conn) Conn() net.Conn {
	return c.opts.conn
}
func (c *Conn) ResetConnDeadline() error {
	return c.opts.conn.SetDeadline(time.Now().Add(time.Duration(c.opts.HeartbeatInterval) * time.Second))
}

func (c *Conn) GetCacheString(k string) string {
	v, ok := c.cache.Load(k)
	if !ok {
		return ""
	}
	vStr, ok := v.(string)
	if !ok {
		return ""
	}
	return vStr
}

func (c *Conn) Write(byt []byte) error {
	select {
	case c.writerChan <- byt:
	case <-c.stopChan:
		return AlreadyCloseErr
	}
	return nil
}

func (c *Conn) readChannel() {
	for {
		byt, err := c.read()
		if err != nil {
			log.Printf("[Conn] readChannel Close conn id : %s , err: %s \n", c.opts.id, err.Error())
			c.Close()
			return
		}
		if c.opts.AutoHeartbeat {
			_ = c.ResetConnDeadline()
		}
		select {
		case c.readChan <- byt:
		case <-c.stopChan:
			return
		}
	}
}

func (c *Conn) witerChannel() {
	for {
		select {
		case byt := <-c.writerChan:
			if err := c.witer(byt); err != nil {
				log.Printf("[Conn] witerChannel id : %s , err: %s \n", c.opts.id, err.Error())
			}
		case <-c.stopChan:
			break
		}
	}
}

func (c *Conn) read() ([]byte, error) {
	lengthByte, err := c.readWriter.Reader.Peek(c.opts.length)
	if err != nil {
		return nil, err
	}
	var length int
	if err = binary.Read(bytes.NewReader(lengthByte), c.opts.endian, &length); err != nil {
		return nil, err
	}

	if c.readWriter.Reader.Buffered() < int(c.opts.length+length) {
		return nil, errors.New("the corresponding data cannot be read")
	}
	pack := make([]byte, int(c.opts.length+length))

	if _, err = c.readWriter.Reader.Read(pack); err != nil {
		return nil, err
	}
	return pack[c.opts.length:], err
}

func (c *Conn) witer(byt []byte) error {
	var (
		length = len(byt)
	)
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, length); err != nil {
		return err
	}
	if err := binary.Write(c.readWriter.Writer, c.opts.endian, byt); err != nil {
		return err
	}
	return c.readWriter.Writer.Flush()
}
