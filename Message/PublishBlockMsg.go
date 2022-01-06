package Message

import (
	"MIS-BC/MetaData"
	"fmt"
)

//go:generate msgp
type PublishBlockMsg struct {
	Height    int    `msg:"Height"`
	Block_num uint32 `msg:"BlockNum"`
	Block     []byte `msg:"Block"`
}

func (msg *PublishBlockMsg) SetBlock(b MetaData.Block) {
	msg.Block, _ = b.MarshalMsg(nil)
	msg.Block_num = b.BlockNum
	msg.Height = b.Height
}

func (msg *PublishBlockMsg) GetBlock() (block MetaData.Block) {
	_, err := block.UnmarshalMsg(msg.Block)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (msg *PublishBlockMsg) GetHeight() int {
	return msg.Height
}

func (msg *PublishBlockMsg) GetBlockNum() uint32 {
	return msg.Block_num
}

func (msg *PublishBlockMsg) ToByteArray() ([]byte, error) {
	return msg.MarshalMsg(nil)
}

func (msg *PublishBlockMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	return err
}

func (manager *MessagerManager) CreatePublishBlockMsg(b MetaData.Block, receiver uint64) (header MessageHeader, msg PublishBlockMsg) {
	msg.SetBlock(b)
	header = manager.CreateHeader(receiver, NormalPublishBlock, 0, 0)
	return
}
