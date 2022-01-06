package Node

import (
	"MIS-BC/Message"
	"MIS-BC/MetaData"
)

func (node *Node) ApplyForVoter() {
	var itmsg MetaData.IdentityTransformation
	itmsg.Type = "ApplyForVoter"
	itmsg.Pubkey = node.config.MyPubkey
	itmsg.SetNodeId(node.network.MyNodeInfo.ID)
	itmsg.IPAddr = node.network.MyNodeInfo.IP
	itmsg.Port = node.network.MyNodeInfo.PORT

	var header MetaData.TransactionHeader
	header.TXType = MetaData.IdTransformation

	var message Message.TransactionMessage
	message.Data = MetaData.EncodeTransaction(header, &itmsg)

	var messageheader Message.MessageHeader
	messageheader.Pubkey = node.config.MyPubkey
	messageheader.Sender = node.network.MyNodeInfo.ID
	messageheader.Receiver = node.accountManager.VoterSet[node.accountManager.WorkerNumberSet[node.dutyWorkerNumber]]
	messageheader.MsgType = Message.TransactionMsg

	node.SendMessage(messageheader, &message)
}

func (node *Node) ApplyForWorkerCandidate() {
	var itmsg MetaData.IdentityTransformation
	itmsg.Type = "ApplyForWorkerCandidate"
	itmsg.Pubkey = node.config.MyPubkey
	itmsg.SetNodeId(node.network.MyNodeInfo.ID)
	itmsg.IPAddr = node.network.MyNodeInfo.IP
	itmsg.Port = node.network.MyNodeInfo.PORT

	var header MetaData.TransactionHeader
	header.TXType = MetaData.IdTransformation

	var message Message.TransactionMessage
	message.Data = MetaData.EncodeTransaction(header, &itmsg)

	var messageheader Message.MessageHeader
	messageheader.Pubkey = node.config.MyPubkey
	messageheader.Sender = node.network.MyNodeInfo.ID
	messageheader.Receiver = node.accountManager.VoterSet[node.accountManager.WorkerNumberSet[node.dutyWorkerNumber]]
	messageheader.MsgType = Message.TransactionMsg

	node.SendMessage(messageheader, &message)
}

func (node *Node) QuitVoter() {
	var itmsg MetaData.IdentityTransformation
	itmsg.Type = "QuitVoter"
	itmsg.Pubkey = node.config.MyPubkey
	itmsg.SetNodeId(node.network.MyNodeInfo.ID)
	itmsg.IPAddr = node.network.MyNodeInfo.IP
	itmsg.Port = node.network.MyNodeInfo.PORT

	var header MetaData.TransactionHeader
	header.TXType = MetaData.IdTransformation

	var message Message.TransactionMessage
	message.Data = MetaData.EncodeTransaction(header, &itmsg)

	var messageheader Message.MessageHeader
	messageheader.Pubkey = node.config.MyPubkey
	messageheader.Sender = node.network.MyNodeInfo.ID
	messageheader.Receiver = node.accountManager.VoterSet[node.accountManager.WorkerNumberSet[node.dutyWorkerNumber]]
	messageheader.MsgType = Message.TransactionMsg

	node.SendMessage(messageheader, &message)
}

func (node *Node) QuitWorkerCandidate() {
	var itmsg MetaData.IdentityTransformation
	itmsg.Type = "QuitWorkerCandidate"
	itmsg.Pubkey = node.config.MyPubkey
	itmsg.SetNodeId(node.network.MyNodeInfo.ID)
	itmsg.IPAddr = node.network.MyNodeInfo.IP
	itmsg.Port = node.network.MyNodeInfo.PORT

	var header MetaData.TransactionHeader
	header.TXType = MetaData.IdTransformation

	var message Message.TransactionMessage
	message.Data = MetaData.EncodeTransaction(header, &itmsg)

	var messageheader Message.MessageHeader
	messageheader.Pubkey = node.config.MyPubkey
	messageheader.Sender = node.network.MyNodeInfo.ID
	messageheader.Receiver = node.accountManager.VoterSet[node.accountManager.WorkerNumberSet[node.dutyWorkerNumber]]
	messageheader.MsgType = Message.TransactionMsg

	node.SendMessage(messageheader, &message)
}
