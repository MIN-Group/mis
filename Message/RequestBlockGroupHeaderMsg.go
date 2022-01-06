package Message

import (
	"fmt"

	"MIS-BC/MetaData"
)

//go:generate msgp
type RequestBlockGroupHeaderMsg struct {
	Height int `msg:"Height"`
	//NodeID uint64 `msg:"NodeID"`
}

type RespondBlockGroupHeaderMsg struct {
	Height                int    `msg:"Height"`
	BlockGroupHeaderBytes []byte `msg:"BlockGroupHeader"`
}

func (msg *RequestBlockGroupHeaderMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestBlockGroupHeaderMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestBlockGroupHeaderMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRequestBlockGroupHeaderMsg(receiver uint64, Height int) (header MessageHeader, msg RequestBlockGroupHeaderMsg) {
	msg.Height = Height
	header = manager.CreateHeader(receiver, RequestBlockGroupHeader, 0, 0)
	return
}

func (msg *RespondBlockGroupHeaderMsg) ToByteArray() ([]byte, error) {
	data, err := msg.MarshalMsg(nil)
	return data, err
}

func (msg *RespondBlockGroupHeaderMsg) FromByteArray(data []byte) error {
	data, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRespondBlockGroupHeaderMsg(receiver uint64, height int, group *MetaData.BlockGroup) (header MessageHeader, msg RespondBlockGroupHeaderMsg) {
	msg.Height = height
	msg.BlockGroupHeaderBytes, _ = group.ToHeaderBytes(nil)
	header = manager.CreateHeader(receiver, ResponseBlockGroupHeader, 0, 0)
	return
}
