package Message

import "fmt"

//go:generate msgp
type ElectNewWorkerMsg struct {
	NewWorker []string `msg:"newworker"`
	MyPubkey  string   `msg:"mypubkey"`
}

func (gm ElectNewWorkerMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *ElectNewWorkerMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
