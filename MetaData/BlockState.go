package MetaData

import "fmt"

//go:generate msgp
type BlockState struct {
	Height int `msg:"height"`
}

func (min BlockState) ToByteArray() []byte {
	data, _ := min.MarshalMsg(nil)
	return data
}

func (min *BlockState) FromByteArray(data []byte) {
	_, err := min.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
