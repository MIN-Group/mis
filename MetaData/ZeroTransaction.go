package MetaData

import "fmt"

//go:generate msgp
type ZeroTransaction struct {
	Content []byte `msg:"zerotx"`
}

func (zt ZeroTransaction) ToByteArray() []byte {
	data, _ := zt.MarshalMsg(nil)
	return data
}

func (zt *ZeroTransaction) FromByteArray(data []byte) {
	_, err := zt.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("--------------------------------")
		fmt.Println("err=", err)
	}
}
