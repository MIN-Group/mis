package Message

import "fmt"

//go:generate msgp
type RequestHeightMsg struct {
	Height int `msg:"Height"`
}

func (msg *RequestHeightMsg) ToByteArray() ([]byte, error) {
	//data, _ := msg.MarshalMsg(nil)
	return msg.MarshalMsg(nil)
}
func (msg *RequestHeightMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("RequestHeightMsg-FromByteArray err=", err)
	}
	return err
}

func (manager *MessagerManager) CreateRequestHeightMsg(receiver uint64) (header MessageHeader, msg RequestHeightMsg) {
	header = manager.CreateHeader(receiver, RequestHeight, 0, 0)
	return
}

func (manager *MessagerManager) CreateRespondHeightMsg(receiver uint64, height int) (header MessageHeader, msg RequestHeightMsg) {
	msg.Height = height
	header = manager.CreateHeader(receiver, RespondHeight, 0, 0)
	//fmt.Println("response header sender=",header.Sender)
	return
}
