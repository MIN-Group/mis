package Message

import "fmt"

//go:generate msgp
type ResponseRemoveNodeMsg struct {
	Type   string `msg:"type"`
	NodeID uint64 `msg:"nodeid"`
	Pubkey string `msg:"pubkey"`
	Result int    `msg:"result"`
	Sign   string `msg:"sign"`
}

func (gm ResponseRemoveNodeMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *ResponseRemoveNodeMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
