package ziface

/*
	message
*/

type IMessage interface {
	//获取消息的ID
	GetMsgId() uint32
	//获取消息的长度
	GetMsgLen() uint32
	//获取消息的内容
	GetData() []byte

	//  type/len/value
	SetMsgID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
