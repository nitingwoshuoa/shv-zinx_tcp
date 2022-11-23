package znet

type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

//创建一个message 消息包

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgID(id uint32) {
	m.Id = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
