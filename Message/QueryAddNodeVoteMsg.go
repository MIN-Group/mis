package Message

import "fmt"

//go:generate msgp
type QueryAddNodeVoteMsg struct {
	Type   string `msg:"type"`
	NodeID uint64 `msg:"nodeid"`
	Pubkey string `msg:"pubkey"`
	Result int    `msg:"result"`
	Sign   string `msg:"sign"`
}

func (gm QueryAddNodeVoteMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *QueryAddNodeVoteMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
