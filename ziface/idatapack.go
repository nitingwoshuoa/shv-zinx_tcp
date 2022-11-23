package ziface

/*
	pack unpack
*/

type IDataPack interface {
	//获取包长度
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
