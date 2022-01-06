package Message

import (
	"MIS-BC/MetaData"
	"fmt"
)

//go:generate msgp
type RequestBlockGroupMsg struct {
	Height int `msg:"Height"`
	//NodeID uint64 `msg:"NodeID"`
}

type RespondBlockGroupMsg struct {
	Height int                 `msg:"Height"`
	Group  MetaData.BlockGroup `msg:"-"`
	Data   []byte              `msg:"BlockGroup"`
}

func (msg *RequestBlockGroupMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestBlockGroupMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRequestBlockGroupMsg(receiver uint64, Height int) (header MessageHeader, msg RequestBlockGroupMsg) {
	msg.Height = Height
	header = manager.CreateHeader(receiver, RequestBlockGroup, 0, 0)
	return
}

func (msg *RespondBlockGroupMsg) ToByteArray() ([]byte, error) {
	temp_data, err := msg.Group.ToBytes(nil)
	msg.Data = temp_data
	data, err := msg.MarshalMsg(nil)
	return data, err
}

func (msg *RespondBlockGroupMsg) FromByteArray(data []byte) error {
	data, err := msg.UnmarshalMsg(data)
	data, err = msg.Group.FromBytes(msg.Data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRespondBlockGroupMsg(receiver uint64, height int, group MetaData.BlockGroup) (header MessageHeader, msg RespondBlockGroupMsg) {
	msg.Height = height
	msg.Group = group
	header = manager.CreateHeader(receiver, RespondBlockGroup, 0, 0)
	return
}
