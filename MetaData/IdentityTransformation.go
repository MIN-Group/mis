package MetaData

import (
	"fmt"
	"strconv"
)

//go:generate msgp
type IdentityTransformation struct {
	Type      string `msg:"type"`
	Pubkey    string `msg:"pubkey"`
	NodeID    string `msg:"nodeid"`
	IPAddr    string `msg:"ip"`
	Port      int    `msg:"port"`
	Timestamp int64
}

func (itmsg *IdentityTransformation) SetNodeId(nodeid uint64) {
	itmsg.NodeID = strconv.FormatUint(nodeid, 10)
}

func (itmsg *IdentityTransformation) GetNodeId() uint64 {
	result, err := strconv.ParseUint(itmsg.NodeID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (itmsg IdentityTransformation) ToByteArray() []byte {
	data, _ := itmsg.MarshalMsg(nil)
	return data
}

func (itmsg *IdentityTransformation) FromByteArray(data []byte) {
	_, err := itmsg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
