package Message

import "fmt"

//go:generate msgp
type QueryPubkeyMsg struct {
	Type        int    `msg:"type"`
	Information string `msg:"info"`
	NodeID      uint64 `msg:"nodeid"`
}

func (gm QueryPubkeyMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *QueryPubkeyMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
