package Message

import "fmt"

//go:generate msgp
type BlockMsg struct {
	Data []byte `msg:"data"`
}

func (bm BlockMsg) ToByteArray() ([]byte, error) {
	data, _ := bm.MarshalMsg(nil)
	return data, nil
}

func (bm *BlockMsg) FromByteArray(data []byte) error {
	_, err := bm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
