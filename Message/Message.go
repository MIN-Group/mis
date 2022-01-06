package Message

import (
	_ "encoding/binary"
	_ "unsafe"

	_ "github.com/tinylib/msgp/msgp"
)

type MessagerManager struct {
	Index  uint64
	Pubkey string
	ID     uint64
}

type MessageInterface interface {
	//将包含消息头的消息转换为byte数组
	ToByteArray() ([]byte, error)
	FromByteArray([]byte) error
}

const (
	Zero                     = 0
	QueryPubkey              = 1
	TransactionMsg           = 2
	GenesisBlock             = 3
	NormalPublishBlock       = 4
	NormalBlockVoteMsg       = 5
	RequestHeight            = 6
	RespondHeight            = 7
	RequestBlockGroup        = 8
	RespondBlockGroup        = 9
	ElectNewWorker           = 10
	RequestBlockGroupHeader  = 11
	ResponseBlockGroupHeader = 12
	RequestBlock             = 13
	ResponseBlock            = 14
	BlockGroupHeaderMsg      = 15
	RequestAddNewNode        = 17
	QueryAddNodeVote         = 18
	ResponseQueryAddNodeVote = 19
	PublishIamOk             = 20
	RequestQuitNode          = 21
	RequestRemoveNode        = 22
	ResponseRemoveNode       = 23
	NodeStatus               = 26
	TransactionStatistics    = 27
	TPSMsg                   = 28
)

//go:generate msgp
type MessageHeader struct {
	Index     uint64 `msg:"Index"`
	Sender    uint64 `msg:"Sender"`
	Receiver  uint64 `msg:"Receiver"`
	Pubkey    string `msg:"Pubkey"`
	MsgType   int    `msg:"MsgType"`
	ChildType int    `msg:"ChildType"`
	RespondTo uint64 `msg:"RespondIndex"`
	Sig       string `msg:"Sig"`
	Data      []byte `msg:"Data"`
}

func ReadHeaderFromByteArray(data []byte) (header MessageHeader, body []byte, err error) {
	body, err = header.UnmarshalMsg(data)
	return
}

func (manager *MessagerManager) CreateHeader(receiver uint64, msg_type int, child_type int, respond_to uint64) MessageHeader {
	manager.Index++
	header := MessageHeader{
		Receiver:  receiver,
		MsgType:   msg_type,
		Sender:    manager.ID,
		ChildType: child_type,
		RespondTo: respond_to,
		Pubkey:    manager.Pubkey}
	return header
}
