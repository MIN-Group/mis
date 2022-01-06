package Message

import (
	"fmt"
)

//go:generate msgp
type RequestRemoveNodeMsg struct {
	Type   string `msg:"type"`
	Pubkey string `msg:"pubkey"`
	NodeID uint64 `msg:"nodeid"`
}

func (msg *RequestRemoveNodeMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestRemoveNodeMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}
