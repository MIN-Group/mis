package Message

import "fmt"

//go:generate msgp
type TransactionMessage struct {
	Data []byte `msg:"data"`
}

func (gm TransactionMessage) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *TransactionMessage) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
