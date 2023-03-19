package dove

import (
	"errors"
	"github.com/hwholiday/ghost/dove/network"
	"github.com/hwholiday/ghost/utils"
	"sync"
	"sync/atomic"
)

var ErrExceedsLengthLimit = errors.New("exceeds length limit")

type manage struct {
	maxConn  int64
	connNum  int64
	connMap  sync.Map
	groupMap sync.Map
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

func (m *manage) Add(conn network.Conn) error {
	if m.Full() {
		return ErrExceedsLengthLimit
	}
	identity := conn.Identity()
	if old, ok := m.GetConn(identity); ok {
		//关闭老的链接信息，这里可能是异地登陆
		old.(network.Conn).Close()
		m.Del(identity)
	}
	m.connMap.Store(identity, conn)
	m.saveGroup(conn)
	atomic.AddInt64(&m.connNum, 1)
	return nil
}

func (m *manage) saveGroup(conn network.Conn) {
	if conn.Group() == "" {
		return
	}
	var (
		identity = conn.Identity()
		group    = conn.Group()
	)
	arr := m.loadGroup(group)
	if !utils.InStrArr(group, arr) {
		arr = append(arr, identity)
		m.groupMap.Store(group, arr)
	}
}

func (m *manage) delGroup(conn network.Conn) {
	if conn.Group() == "" {
		return
	}
	var (
		identity = conn.Identity()
		group    = conn.Group()
	)
	arr := m.loadGroup(group)
	if utils.InStrArr(group, arr) {
		m.groupMap.Store(group, utils.DelStrArr(identity, arr))
	}
}

func (m *manage) loadGroup(group string) []string {
	identityAny, ok := m.groupMap.Load(group)
	if !ok {
		return nil
	}
	identityArr, ok := identityAny.([]string)
	if !ok {
		return nil
	}
	return identityArr
}

func (m *manage) Del(identity string) {
	conn, ok := m.GetConn(identity)
	if !ok {
		return
	}
	atomic.AddInt64(&m.connNum, -1)
	m.connMap.Delete(conn.Identity())
	m.delGroup(conn)
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

func (m *manage) FindConnByGroup(group string) []network.Conn {
	identityArr := m.loadGroup(group)
	if len(identityArr) <= 0 {
		return nil
	}
	var coonArr = make([]network.Conn, 0, len(identityArr))
	for _, item := range identityArr {
		if conn, ok := m.GetConn(item); ok {
			coonArr = append(coonArr, conn)
		}
	}
	return coonArr
}

func (m *manage) FindAllConn() []network.Conn {
	var clientArr = make([]network.Conn, 0, m.GetConnNum())
	m.connMap.Range(func(key, value any) bool {
		clientArr = append(clientArr, value.(network.Conn))
		return true
	})
	return clientArr
}
