package dove

import (
	"errors"
	"github.com/hwholiday/ghost/dove/network"
	"sync"
	"sync/atomic"
)

var ErrExceedsLengthLimit = errors.New("exceeds length limit")

type manage struct {
	maxConn int64
	connNum int64
	connMap sync.Map
}

func newManage() *manage {
	return &manage{maxConn: DefaultConnMax}
}
func (m *manage) Full() bool {
	if m.GetConnNum() >= m.maxConn {
		return true
	}
	return false
}

func (m *manage) Add(identity string, conn network.Conn) error {
	if m.Full() {
		return ErrExceedsLengthLimit
	}
	if old, ok := m.GetConn(identity); ok {
		//关闭老的链接信息，这里可能是异地登陆
		old.(network.Conn).Close()
		m.Del(identity)
	}
	m.connMap.Store(identity, conn)
	atomic.AddInt64(&m.connNum, 1)
	return nil
}

func (m *manage) Del(identity string) {
	if _, ok := m.connMap.Load(identity); !ok {
		return
	}
	atomic.AddInt64(&m.connNum, -1)
	m.connMap.Delete(identity)
}

func (m *manage) GetConnNum() int64 {
	return atomic.LoadInt64(&m.connNum)
}

func (m *manage) GetMapStatus() map[string]interface{} {
	var result = make(map[string]interface{}, 2)
	result["connNum"] = m.GetConnNum()
	result["maxNum"] = m.maxConn
	return result
}

func (m *manage) GetConn(identity string) (network.Conn, bool) {
	val, ok := m.connMap.Load(identity)
	if !ok {
		return nil, false
	}
	return val.(network.Conn), true
}

func (m *manage) GetAllConn() []network.Conn {
	var clientArr = make([]network.Conn, 0, m.GetConnNum())
	m.connMap.Range(func(key, value any) bool {
		clientArr = append(clientArr, value.(network.Conn))
		return true
	})
	return clientArr
}
