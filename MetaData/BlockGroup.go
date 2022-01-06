package MetaData

/*func init() {
	// Registering an extension is as simple as matching the
	// appropriate type number with a function that initializes
	// a freshly-allocated object of that type
	msgp.RegisterExtension(10, func() msgp.Extension { return new(Block) } )
}*/
//go:generate msgp
type BlockGroup struct {
	Height         int          `msg:"height"`
	Generator      string       `msg:"generator"`
	PreviousHash   string       `msg:"preHash"`
	MerkleRoot     string       `msg:"merkleRoot"`
	VoteTickets    []VoteTicket `msg:"-"`
	Timestamp      float64      `msg:"timestamp"`
	VoteResult     []int        `msg:"VoteResult"`
	BlockHashes    []string     `msg:"BlockHashes"`
	NextDutyWorker uint32       `msg:"NextDutyWorker"`
	Sig            string       `msg:"Sig"`

	ReceivedBlockGroupHeader bool `msg:"-"`

	Blocks            []Block `msg:"-"`
	CheckTransactions []int   `msg:"-"`
	CheckHeader       []int   `msg:"-"` //对每一个区块的投票
}

func (bg *BlockGroup) ToBytes(raw_data []byte) (data []byte, err error) {
	//封装区块头
	var bytehelper3 ByteArrayHelper
	temp_data, err := bg.MarshalMsg(nil)
	bytehelper3.Data = temp_data
	data, err = bytehelper3.MarshalMsg(data)
	//封装投票集合
	var bytehelper2 ByteArrayHelper
	var helper2 DoubleByteArrayHelper
	for i := 0; i < len(bg.VoteTickets); i++ {
		ticket, _ := bg.VoteTickets[i].MarshalMsg(nil)
		//fmt.Println("ticket=",ticket)
		helper2.Data = append(helper2.Data, ticket)
	}
	//fmt.Println("helper2=",helper2)
	temp_data, err = helper2.MarshalMsg(nil)
	bytehelper2.Data = temp_data
	//fmt.Println("bytehelper2",bytehelper2)
	data, err = bytehelper2.MarshalMsg(data)

	//封装区块集合
	var bytehelper1 ByteArrayHelper
	var helper DoubleByteArrayHelper
	for i := 0; i < len(bg.Blocks); i++ {
		block, _ := bg.Blocks[i].MarshalMsg(nil)
		helper.Data = append(helper.Data, block)
	}
	temp_data, err = helper.MarshalMsg(nil)
	bytehelper1.Data = temp_data
	data, err = bytehelper1.MarshalMsg(data)
	return
}

func (bg *BlockGroup) FromBytes(raw_data []byte) (data []byte, err error) {
	//读取区块头
	var bytehelper3 ByteArrayHelper
	data, err = bytehelper3.UnmarshalMsg(raw_data)
	_, err = bg.UnmarshalMsg(bytehelper3.Data)

	//读取投票集合
	var bytehelper2 ByteArrayHelper
	data, err = bytehelper2.UnmarshalMsg(data)
	var helper2 DoubleByteArrayHelper
	_, err = helper2.UnmarshalMsg(bytehelper2.Data)
	for i := 0; i < len(helper2.Data); i++ {
		var ticket VoteTicket
		ticket.UnmarshalMsg(helper2.Data[i])
		bg.VoteTickets = append(bg.VoteTickets, ticket)
	}

	//读取区块集合
	var bytehelper1 ByteArrayHelper
	data, err = bytehelper1.UnmarshalMsg(data)
	var helper DoubleByteArrayHelper
	_, err = helper.UnmarshalMsg(bytehelper1.Data)
	for i := 0; i < len(helper.Data); i++ {
		var block Block
		block.UnmarshalMsg(helper.Data[i])
		bg.Blocks = append(bg.Blocks, block)
	}
	return
}

func (bg *BlockGroup) ToHeaderBytes(raw_data []byte) (data []byte, err error) {
	//封装区块头
	var bytehelper3 ByteArrayHelper
	temp_data, err := bg.MarshalMsg(nil)
	bytehelper3.Data = temp_data
	data, err = bytehelper3.MarshalMsg(data)
	//封装投票集合
	var bytehelper2 ByteArrayHelper
	var helper2 DoubleByteArrayHelper
	for i := 0; i < len(bg.VoteTickets); i++ {
		ticket, _ := bg.VoteTickets[i].MarshalMsg(nil)
		//fmt.Println("ticket=",ticket)
		helper2.Data = append(helper2.Data, ticket)
	}
	//fmt.Println("helper2=",helper2)
	temp_data, err = helper2.MarshalMsg(nil)
	bytehelper2.Data = temp_data
	//fmt.Println("bytehelper2",bytehelper2)
	data, err = bytehelper2.MarshalMsg(data)
	return
}

func (bg *BlockGroup) FromHeaderBytes(raw_data []byte) (data []byte, err error) {
	//读取区块头
	var bytehelper3 ByteArrayHelper
	data, err = bytehelper3.UnmarshalMsg(raw_data)
	_, err = bg.UnmarshalMsg(bytehelper3.Data)

	//读取投票集合
	var bytehelper2 ByteArrayHelper
	data, err = bytehelper2.UnmarshalMsg(data)
	var helper2 DoubleByteArrayHelper
	_, err = helper2.UnmarshalMsg(bytehelper2.Data)
	for i := 0; i < len(helper2.Data); i++ {
		var ticket VoteTicket
		ticket.UnmarshalMsg(helper2.Data[i])
		bg.VoteTickets = append(bg.VoteTickets, ticket)
	}
	return
}

func (bg *BlockGroup) ToBlocksBytes(raw_data []byte) (data []byte, err error) {
	//封装区块集合
	var bytehelper1 ByteArrayHelper
	var helper DoubleByteArrayHelper
	for i := 0; i < len(bg.Blocks); i++ {
		block, _ := bg.Blocks[i].MarshalMsg(nil)
		helper.Data = append(helper.Data, block)
	}
	temp_data, err := helper.MarshalMsg(nil)
	bytehelper1.Data = temp_data
	data, err = bytehelper1.MarshalMsg(raw_data)
	return
}

func (bg *BlockGroup) FromBlocksBytes(raw_data []byte) (data []byte, err error) {
	//读取区块集合
	var bytehelper1 ByteArrayHelper
	data, err = bytehelper1.UnmarshalMsg(raw_data)
	var helper DoubleByteArrayHelper
	_, err = helper.UnmarshalMsg(bytehelper1.Data)
	for i := 0; i < len(helper.Data); i++ {
		var block Block
		block.UnmarshalMsg(helper.Data[i])
		bg.Blocks = append(bg.Blocks, block)
	}
	return
}
