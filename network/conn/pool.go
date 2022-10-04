package conn

import "sync"

var connPool = sync.Pool{
	New: func() interface{} {
		return &Conn{}
	},
}

func getConn() *Conn {
	cli := connPool.Get().(*Conn)
	return cli
}

func putConn(cli *Conn) {
	connPool.Put(cli)
}
