package Message

import (
	"fmt"
)

//go:generate msgp
type RequestAddNewNodeMsg struct {
	Type   string `msg:"type"`
	Pubkey string `msg:"pubkey"`
	NodeID uint64 `msg:"nodeid"`
	IPAddr string `msg:"ip"`
	Port   int    `msg:"port"`
}

func (msg *RequestAddNewNodeMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestAddNewNodeMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRequestAddNewNodeMsg(receiver uint64, typ string, Pubkey string, NodeID uint64, IPAddr string, Port int) (header MessageHeader, msg RequestAddNewNodeMsg) {
	msg.Type = typ
	msg.Pubkey = Pubkey
	msg.NodeID = NodeID
	msg.IPAddr = IPAddr
	msg.Port = Port
	header = manager.CreateHeader(receiver, RequestAddNewNode, 0, 0)
	return
}
