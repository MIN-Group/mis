package MetaData

import "fmt"

//go:generate msgp
type ElectNewWorkerTeam struct {
	WorkerPubList map[string]uint64 `msg:"WorkerPubList"`
	ElectResult   map[string]int    `msg:"electresult"`
}

func (itmsg ElectNewWorkerTeam) ToByteArray() []byte {
	data, _ := itmsg.MarshalMsg(nil)
	return data
}

func (itmsg *ElectNewWorkerTeam) FromByteArray(data []byte) {
	_, err := itmsg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
