package Message

import (
	"fmt"

	"MIS-BC/MetaData"
)

//go:generate msgp
type RequestBlockMsg struct {
	Height   int `msg:"Height"`
	BlockNum int `msg:"BlockNum"`
	//NodeID uint64 `msg:"NodeID"`
}

type RespondBlockMsg struct {
	Height   int            `msg:"Height"`
	BlockNum int            `msg:"BlockNum"`
	Block    MetaData.Block `msg:"-"`
	Data     []byte         `msg:"Data"`
}

func (msg *RequestBlockMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestBlockMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRequestBlockMsg(receiver uint64, Height int, BlockNum int) (header MessageHeader, msg RequestBlockMsg) {
	msg.Height = Height
	msg.BlockNum = BlockNum
	header = manager.CreateHeader(receiver, RequestBlock, 0, 0)
	return
}

func (msg *RespondBlockMsg) ToByteArray() ([]byte, error) {
	temp, err := msg.Block.MarshalMsg(nil)
	msg.Data = temp
	data, err := msg.MarshalMsg(nil)
	return data, err
}

func (msg *RespondBlockMsg) FromByteArray(data []byte) error {
	data, err := msg.UnmarshalMsg(data)
	data, err = msg.Block.UnmarshalMsg(msg.Data)
	if err != nil {
		fmt.Println("RespondBlockMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRespondBlockMsg(receiver uint64, height int, blockNum int, block MetaData.Block) (header MessageHeader, msg RespondBlockMsg) {
	msg.Height = height
	msg.BlockNum = blockNum
	msg.Block = block
	header = manager.CreateHeader(receiver, ResponseBlock, 0, 0)
	return
}
