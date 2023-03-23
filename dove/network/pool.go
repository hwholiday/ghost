package network

import "sync"

var connPool = sync.Pool{
	New: func() interface{} {
		return &conn{}
	},
}

func getConn() *conn {
	cli := connPool.Get().(*conn)
	return cli
}

func putConn(cli *conn) {
	connPool.Put(cli)
}

var connWsPool = sync.Pool{
	New: func() interface{} {
		return &wsConn{}
	},
}

func getWsConn() *wsConn {
	cli := connWsPool.Get().(*wsConn)
	return cli
}

func putWsConn(cli *wsConn) {
	connWsPool.Put(cli)
}
