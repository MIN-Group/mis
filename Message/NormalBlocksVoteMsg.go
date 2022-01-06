package Message

import "MIS-BC/MetaData"

//go:generate msgp
type NormalBlocksVoteMsg struct {
	Height   int                 `msg:"Height"`
	BlockNum uint32              `msg:"BlockNum"`
	Data     []byte              `msg:"VoteTicket"`
	Ticket   MetaData.VoteTicket `msg:"-"`
}

func (msg *NormalBlocksVoteMsg) ToByteArray() ([]byte, error) {
	temp_data, err := msg.Ticket.MarshalMsg(nil)
	msg.Data = temp_data
	data, err := msg.MarshalMsg(nil)
	return data, err
}

func (msg *NormalBlocksVoteMsg) FromByteArray(data []byte) error {
	_, err := msg.UnmarshalMsg(data)
	_, err = msg.Ticket.UnmarshalMsg(msg.Data)
	return err
}

func (manager *MessagerManager) CreateNormalBlocksVoteMsg(ticket MetaData.VoteTicket,
	receiver uint64,
	height int,
	blocknum uint32) (header MessageHeader, msg NormalBlocksVoteMsg) {
	msg.Ticket = ticket
	msg.Height = height
	msg.BlockNum = blocknum
	header = manager.CreateHeader(receiver, NormalBlockVoteMsg, 0, 0)
	return
}
