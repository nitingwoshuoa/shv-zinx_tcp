package ziface

/*
	IRequest接口



*/
type IRequest interface {

	// get curr connection
	GetConnection() IConnection

	GetData() []byte
	GetMsgID() uint32
}
