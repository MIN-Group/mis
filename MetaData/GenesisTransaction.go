package MetaData

import "fmt"

//type NodeID = uint64
//go:generate msgp
type GenesisTransaction struct {
	WorkerNum              int               `msg:"WorkerNum"`              //记账节点数量
	VotedNum               int               `msg:"VotedNum"`               //投票选举的节点数量
	BlockGroupPerCycle     int               `msg:"BlockGroupPerCycle"`     //任职周期
	Tcut                   float64           `msg:"Tcut"`                   //共识截止时间
	WorkerPubList          map[string]uint64 `msg:"WorkerPubList"`          //记账节点公钥列表
	WorkerCandidatePubList map[string]uint64 `msg:"WorkerCandidatePubList"` //记账候选节点公钥列表
	VoterPubList           map[string]uint64 `msg:"VoterPubList"`           //投票节点公钥列表
	WorkerSet              []string          `msg:"WNS"`                    //记账节点列表
	VoterSet               []string          `msg:"VS"`                     //投票节点列表
}

func (gt GenesisTransaction) ToByteArray() []byte {
	data, _ := gt.MarshalMsg(nil)
	return data
}

func (gt *GenesisTransaction) FromByteArray(data []byte) {
	_, err := gt.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
