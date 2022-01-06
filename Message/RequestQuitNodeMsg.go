package Message

import "fmt"

//go:generate msgp
type RequestQuitNodeMsg struct {
	Type   string `msg:"type"`
	Pubkey string `msg:"pubkey"`
	NodeID uint64 `msg:"nodeid"`
}

func (gm RequestQuitNodeMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *RequestQuitNodeMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
