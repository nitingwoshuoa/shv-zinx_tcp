package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/nitingwoshuoa/shv-zinx_tcp/utils"
	"github.com/nitingwoshuoa/shv-zinx_tcp/ziface"
)

/*
	封包拆包模块

	直接面向TCP链接的数据流， 用于处理TCPP粘包问题
*/

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节)+ ID uint32(4字节)
	return 8
}

// 封包方法： /datalen/msgid/data/
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	fmt.Println("msg id p ==", msg.GetMsgId())
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil

}

func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
