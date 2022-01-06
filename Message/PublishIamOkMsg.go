package Message

import "fmt"

//go:generate msgp
type PublishIamOkMsg struct {
	Type   string `msg:"type"`
	Pubkey string `msg:"pubkey"`
}

func (gm PublishIamOkMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *PublishIamOkMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
