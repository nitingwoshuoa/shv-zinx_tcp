package ziface

/*
	conn manager
*/

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, error)
	//conn total count
	Len() int
	ClearConn()
}
