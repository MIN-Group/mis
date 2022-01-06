package Message

import "fmt"

//go:generate msgp
type BlockGroupHeader struct {
	Data []byte `msg:"data"`
}

func (bm BlockGroupHeader) ToByteArray() ([]byte, error) {
	data, _ := bm.MarshalMsg(nil)
	return data, nil
}

func (bm *BlockGroupHeader) FromByteArray(data []byte) error {
	_, err := bm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
