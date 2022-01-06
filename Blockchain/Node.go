package Node

import (
	"MIS-BC/AccountManager"
	"MIS-BC/Database"
	"MIS-BC/Message"
	"MIS-BC/MetaData"
	"MIS-BC/Network"
	"MIS-BC/TransactionPool"
	"MIS-BC/common"
	"MIS-BC/security"
	"MIS-BC/security/keymanager"
	"MIS-BC/utils"
	"encoding/base64"
	"fmt"
	"github.com/karlseguin/ccache/v2"
	"github.com/patrickmn/go-cache"
	"github.com/smallnest/rpcx/server"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	Sync    = 0
	Genesis = 1
	Normal  = 2
)

type CallBackInstance struct {
	run       func(msg Message.MessageInterface, header Message.MessageHeader)
	MsgType   int
	ChildType int
	StartTime float64
	WaitTime  float64
}

type Node struct {
	network        *Network.Network
	mongo          *MongoDB.Mongo
	SessionCache   *cache.Cache
	keyManager     *keymanager.KeyManager
	accountManager *AccountManager.AccountManager
	msgManager     *Message.MessagerManager
	config         *common.Config
	txPool         *TransactionPool.TransactionPool
	syncTool       *SynchronizeModule
	//各种控制状态
	state chan int //系统状态

	//共识状态变量
	dutyWorkerNumber      uint32 //轮值记账节点编号
	true_dutyWorkerNum    uint32
	round                 uint32  //当前经过轮数
	StartTime             float64 //当前共识开始时间
	Tcut                  float64 //一轮的时间
	isTimeOut, isTimeout2 bool    //超时flag

	BlockGroups *sync.Map
	RWLock      *sync.RWMutex
	AddNodeLock *sync.Mutex

	CallBackList []CallBackInstance

	LoginTable *sync.Map

	//genesis block temp
	WorkerPubList          map[string]uint64   //记账节点公钥列表
	WorkerCandidatePubList map[string]uint64   //记账候选节点公钥列表
	VoterPubList           map[string]uint64   //投票节点公钥列表
	ElectNewWorkerList     []Message.ElectNewWorkerMsg

	GenesisBlockDone                   bool //是否生成了创世区块
	NormalGenerateBlockDone            bool
	NormalGenerateVoteDone             bool
	NormalGenerateBlockGroupHeaderDone bool
	registryList        map[string]int
	StartGetVoteTime    float64
	LastGetVoteTime     float64
	LastReqIamOkTime    float64
	//AddNode
	AddNewNodeVoteList map[string]int //新加入节点请求现有节点的投票列表
	//RemoveNode
	RemoveNodeVoteList map[uint64]int //申请删除节点投票的列表

	//延迟处理和节点增删相关的交易
	IdentityTransList map[MetaData.IdentityTransformation]int

	//在该列表中的节点可以同意加入区块链网络
	NewNodeList    *ccache.Cache
	timestampCache *sync.Map //时间戳缓存，防止重放攻击

	BCStatus *BCStatus // 区块链状态
}

func (node *Node) Init() {
	//初始化网络模块
	node.syncTool = &SynchronizeModule{}
	node.syncTool.Init()
	node.msgManager = &Message.MessagerManager{}
	node.txPool = &TransactionPool.TransactionPool{}
	key := new(keymanager.KeyManager)
	key.Init()
	node.network = &Network.Network{
		MyNodeInfo: Network.NodeInfo{
			IP:          "",
			PORT:        0,
			ID:          0,
			Prefix:      "",
			HostName:    "",
			AreaName:    "",
			CountryName: "",
			Longitude:   0,
			Latitude:    0,
		},
		NodeList:      make(map[Network.NodeID]Network.NodeInfo),
		Blacklist:     make(map[Network.NodeID]Network.Void),
		CBforBC:       nil,
		Mutex:         new(sync.RWMutex),
		Keychain:      &security.KeyChain{},
	}
	node.network.SetCB(node.HandleBCMessage)

	node.mongo = &MongoDB.Mongo{}
	node.keyManager = &keymanager.KeyManager{}
	// node.rsaManager = &KeyManager.KeyManager{}
	node.keyManager.Init()
	// node.rsaManager.Init(conf.Key.MINVPN)
	node.state = make(chan int, 3)
	node.state <- Sync
	node.accountManager = &AccountManager.AccountManager{}
	node.accountManager.WorkerSet = make(map[string]uint64)
	node.accountManager.VoterSet = make(map[string]uint64)
	node.accountManager.WorkerCandidateSet = make(map[string]uint64)
	node.accountManager.WorkerNumberSet = make(map[uint32]string)
	node.accountManager.VoterNumberSet = make(map[uint32]string)
	node.WorkerPubList = make(map[string]uint64)
	node.WorkerCandidatePubList = make(map[string]uint64)
	node.VoterPubList = make(map[string]uint64)
	node.GenesisBlockDone = false
	node.NormalGenerateBlockDone = false
	node.NormalGenerateVoteDone = false
	node.NormalGenerateBlockGroupHeaderDone = false
	node.StartTime = utils.GetCurrentTime()
	node.dutyWorkerNumber = 0
	node.StartGetVoteTime = 0
	node.LastGetVoteTime = 0
	node.LastReqIamOkTime = 0
	node.RWLock = new(sync.RWMutex)
	node.BlockGroups = new(sync.Map)
	node.LoginTable = new(sync.Map)
	node.registryList = make(map[string]int)
	node.timestampCache = new(sync.Map)
	node.AddNodeLock = new(sync.Mutex)
	node.AddNewNodeVoteList = make(map[string]int)
	node.IdentityTransList = make(map[MetaData.IdentityTransformation]int)
	node.NewNodeList = ccache.New(ccache.Configure())
	node.RemoveNodeVoteList = make(map[uint64]int)
	node.timestampCache = new(sync.Map)
	node.BCStatus = new(BCStatus)
	node.BCStatus.Mutex = new(sync.RWMutex)

}

func (node *Node) SetConfig(config common.Config) {
	node.config = &config
	node.mongo.SetConfig(config)
	node.network.SetConfig(config) //配置记账节点 候选记账节点等
	node.keyManager.SetPubkey(config.MyPubkey)
	node.keyManager.SetPriKey(config.MyPrikey)
	node.txPool.Init(node.config.TxPoolSize, node.config.TxPoolSize)
	node.Tcut = config.Tcut
	node.msgManager.Pubkey = config.MyPubkey
	node.msgManager.ID = node.network.MyNodeInfo.ID
	node.SessionCache = cache.New(time.Duration(config.DefaultExpiration)*time.Minute, time.Duration(config.CleanupInterval)*time.Minute)
	node.LoadBlockChain()
}

func (node *Node) LoadBlockChain() {
	fmt.Println("加载区块链")
	// fmt.Println("当前使用的身份为:", node.network.Keychain.GetCurrentIdentity())
	height := node.mongo.QueryHeight()
	fmt.Println("height:", height)
	//加速区块链重启 ==> 不支持节点身份转换
	if height >= 0 {
		group := node.mongo.GetBlockFromDatabase(0)
		group.ReceivedBlockGroupHeader = true
		node.UpdateGenesisBlockVariables(&group)
		if len(node.state) == 1 {
			<-node.state
		}
		if height > 0 {
			group = node.mongo.GetBlockFromDatabase(height)
			node.UpdateVariablesFromDisk(&group)
		}
		node.mongo.Height = height
		node.mongo.Block = group
		fmt.Println("高度为0-", height, "的区块组加载完成")
	}

	if node.mongo.HasAccount() {
		account := node.mongo.GetAccountFromDatabase()
		node.accountManager.VoterNumberSet = account.GetVorterNumberSet()
		node.accountManager.WorkerNumberSet = account.GetWorkerNumberSet()
		node.accountManager.WorkerCandidateList = account.WorkerCandidateList
		node.accountManager.VoterSet = account.GetVoterSet()
		node.accountManager.WorkerSet = account.GetWorkerSet()
		node.accountManager.WorkerCandidateSet = account.GetWorkerCandidateSet()
	}
	if node.mongo.HasNodelist() {
		nodelist := node.mongo.GetNodeListFromDatabase()
		node.network.NodeList = nodelist.GetNodeList()

	}
	if node.mongo.HasIdentityTransList() {
		list := node.mongo.GetIdentityTransListFromDatabase()
		node.IdentityTransList = list.GetIdentityTransList()
	}
}

func (node *Node) Start() {
	go node.network.Start()
	go node.StartRPCServer()

	node.run()
}

func (node *Node) StartRPCServer() {
	s := server.NewServer()
	var rpcServer = &RpcServer{node}
	err := s.RegisterName("Registry", rpcServer, "")
	if err != nil {
		common.Logger.Info(err)
	}
	port := strconv.Itoa(node.network.MyNodeInfo.PORT + 10)

	err = s.Serve("tcp", "localhost:"+port)
	if err != nil {
		common.Logger.Info(err)
	}
}

func (node *Node) run() {
	for {
		time.Sleep(10 * time.Millisecond)
		state := <-node.state
		switch state {
		case Sync:
			node.SynchronizeBlockGroup()
		case Genesis:
			node.GenerateGenesisBlockGroup()
		case Normal:
			node.Normal()
			node.state <- Normal
		}
	}
}

// 同步区块组
func (node *Node) SynchronizeBlockGroup() {
	//node.state <- Genesis

	if node.config.IsNewJoin == 0 {
		node.syncBlockGroupSequence()
		node.StartTime = utils.GetCurrentTime()
		return
	}

	var agree int = 0
	node.AddNodeLock.Lock()
	for _, v := range node.AddNewNodeVoteList {
		if v == 1 {
			agree++
		}
	}
	defer node.AddNodeLock.Unlock()

	if agree > len(node.accountManager.VoterSet)/2 {
		//fmt.Println(len(node.accountManager.VoterSet))
		fmt.Println("进入同步区块阶段")
		node.syncBlockGroupSequence()
		node.StartTime = utils.GetCurrentTime()
		return
	} else {
		if utils.GetCurrentTime()-node.StartGetVoteTime > node.config.GenerateBlockPeriod+1 {
			node.RequestAddNode()
			node.StartGetVoteTime = utils.GetCurrentTime()
			fmt.Println("请求加入本节点")
		}

		if utils.GetCurrentTime()-node.LastGetVoteTime > 1 {
			node.QueryAllAddNodeVote()
			fmt.Println("请求投票节点投票")
			node.LastGetVoteTime = utils.GetCurrentTime()
		}
		node.StartTime = utils.GetCurrentTime()
		node.state <- Sync
	}
}

// 新节点加入请求
func (node *Node) RequestAddNode() {
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.RequestAddNewNode
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.RequestAddNewNodeMsg //消息体
	gm.Pubkey = node.config.MyPubkey
	gm.Type = "ApplyNode"
	gm.NodeID = node.network.MyNodeInfo.ID
	gm.IPAddr = node.network.MyNodeInfo.IP
	gm.Port = node.network.MyNodeInfo.PORT
	node.SendMessage(mh, &gm)
}

func (node *Node) QueryAllAddNodeVote() { //only for 生成创世节点的节点
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.QueryAddNodeVote
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.QueryAddNodeVoteMsg //消息体
	gm.Type = "ApplyNode"
	gm.NodeID = node.network.MyNodeInfo.ID
	gm.Pubkey = node.config.MyPubkey
	node.SendMessage(mh, &gm)
}

// 生成创世区块组
func (node *Node) GenerateGenesisBlockGroup() {
	if node.config.MyAddress.IP == node.config.WorkerList[0].IP && node.config.MyAddress.Port == node.config.WorkerList[0].Port {
		//if !node.config.IsMINConn && node.config.MyAddress.IP == node.config.WorkerList[0].IP && node.config.MyAddress.Port == node.config.WorkerList[0].Port {
		if !node.GenesisBlockDone {
			node.QueryAllPubkey() //请求公钥
			time.Sleep(time.Second)
			node.RWLock.Lock()

			if !(len(node.WorkerPubList) == len(node.config.WorkerList) &&
				len(node.WorkerCandidatePubList) == len(node.config.WorkerCandidateList) &&
				len(node.VoterPubList) == len(node.config.VoterList)) {
				node.state <- Genesis
				node.RWLock.Unlock()
				return
			}

			var genesisTransaction MetaData.GenesisTransaction //创世交易
			genesisTransaction.WorkerNum = node.config.WorkerNum
			genesisTransaction.VotedNum = node.config.VotedNum
			genesisTransaction.BlockGroupPerCycle = node.config.BlockGroupPerCycle
			genesisTransaction.Tcut = node.config.Tcut
			genesisTransaction.WorkerPubList = node.WorkerPubList
			genesisTransaction.WorkerCandidatePubList = node.WorkerCandidatePubList
			genesisTransaction.VoterPubList = node.VoterPubList

			for key, _ := range genesisTransaction.WorkerPubList {
				genesisTransaction.WorkerSet = append(genesisTransaction.WorkerSet, key)
			}

			for key, _ := range genesisTransaction.VoterPubList {
				genesisTransaction.VoterSet = append(genesisTransaction.VoterSet, key)
			}

			var transactionHeader MetaData.TransactionHeader //交易头
			transactionHeader.TXType = MetaData.Genesis

			var block MetaData.Block //区块
			block.Height = 0
			block.Generator = node.config.MyPubkey
			block.Transactions = append(block.Transactions, MetaData.EncodeTransaction(transactionHeader, &genesisTransaction))
			block.MerkleRoot = keymanager.GetHash(block.GetTransactionsBytes())
			var blockgroup MetaData.BlockGroup //区块组
			blockgroup.Height = 0
			blockgroup.Generator = node.config.MyPubkey
			blockgroup.Timestamp = utils.GetCurrentTime()
			blockgroup.Blocks = append(blockgroup.Blocks, block)
			temp, _ := blockgroup.ToHeaderBytes(nil)
			blockgroup.Sig, _ = node.keyManager.Sign(temp)

			var blockmsg Message.BlockMsg //消息体
			blockmsg.Data, _ = blockgroup.ToBytes(nil)

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = node.network.MyNodeInfo.ID
			msgheader.Receiver = 0
			msgheader.Pubkey = node.config.MyPubkey
			msgheader.MsgType = Message.GenesisBlock

			node.SendMessage(msgheader, &blockmsg)
			node.GenesisBlockDone = true //修改状态
			node.RWLock.Unlock()
			return
		}
	}
}

func (node *Node) QueryAllPubkey() { //only for 生成创世节点的节点
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.QueryPubkey
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.QueryPubkeyMsg //消息体
	gm.Type = 100
	node.SendMessage(mh, &gm)
}

func (node *Node) Normal() {

	if node.config.IsNewJoin == 1 {
		if _, ok := node.accountManager.VoterSet[node.config.MyPubkey]; !ok {
			if utils.GetCurrentTime()-node.LastReqIamOkTime > node.config.Tcut*5 {
				node.SendIamOkMsgWhenNewJoin()
				node.LastReqIamOkTime = utils.GetCurrentTime()
			}
		}

	}

	var height = node.mongo.GetHeight() + 1
	//删除过时的区块组
	if height%99 == 0 {
		node.BlockGroups.Range(func(k, _ interface{}) bool {
			if k.(int) < height-10 {
				node.BlockGroups.Delete(k)
			}
			return true
		})
		node.timestampCache.Range(func(k, v interface{}) bool {
			timestamp, _ := strconv.ParseFloat(v.(string), 64)
			if utils.GetCurrentTime()-timestamp > 60 {
				node.timestampCache.Delete(k)
			}
			return true
		})
	}
	////若还没创建当前所需共识的区块组，则创建一个
	for h := height; h < height+2; h++ {
		node.BlockGroups.LoadOrStore(h, node.CreateBlockGroup())
	}
	//更新当前轮数和轮值记账节点
	round := uint32((utils.GetCurrentTime() - node.StartTime) / node.Tcut)
	if round != node.round {
		node.isTimeOut = true
		//	node.isTimeout2 = true
		node.NormalGenerateVoteDone = false
		node.round = round
		fmt.Println("轮数切换为", round)
	}
	if (utils.GetCurrentTime()-node.StartTime)/node.Tcut > 1.5 {
		node.isTimeout2 = true
	}

	//经过3轮超时后才进行轮值记账节点的更换
	if round >= 3 {
		node.true_dutyWorkerNum = (node.dutyWorkerNumber + round - 1) % uint32(node.config.WorkerNum)
	} else {
		node.true_dutyWorkerNum = node.dutyWorkerNumber
	}
	//判断自己的身份
	pubkey := node.config.MyPubkey
	_, isWorkerCandidate := node.accountManager.WorkerCandidateSet[pubkey]
	_, isWorker := node.accountManager.WorkerSet[pubkey]
	_, isVoter := node.accountManager.VoterSet[pubkey]
	//所有节点都需要执行Commit提交区块组数据
	node.Commit()
	//在主线程检查区块头是否正确
	node.CheckBlocksHeader()
	//投票节点执行生成投票过程
	if isVoter {
		node.GenerateVote()
	}
	//记账节点生成区块
	if isWorker {
		node.GenerateBlock()
		//轮值记账节点产生区块组头
		//fmt.Println("node.true_dutyWorkerNum=",node.true_dutyWorkerNum)
		duty_pubkey, _ := node.accountManager.WorkerNumberSet[node.true_dutyWorkerNum]
		if duty_pubkey == pubkey {
			node.GenerateBlockGroupHeader()
		}
	}
	//候选记账节点不需要做任何事情
	if isWorkerCandidate {

	}
}

//当本节点同步完成就可以告诉其他节点我已经同步完
func (node *Node) SendIamOkMsgWhenNewJoin() {
	fmt.Println("新节点发送同步完成消息")
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.PublishIamOk
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.PublishIamOkMsg //消息体
	gm.Type = "IamOk"
	gm.Pubkey = node.config.MyPubkey
	node.SendMessage(mh, &gm)
}

// 生成区块组
func (node *Node) CreateBlockGroup() MetaData.BlockGroup {
	var group MetaData.BlockGroup

	group.VoteTickets = make([]MetaData.VoteTicket, len(node.accountManager.VoterSet))
	group.Blocks = make([]MetaData.Block, node.config.WorkerNum)
	group.CheckTransactions = make([]int, node.config.WorkerNum)
	group.CheckHeader = make([]int, node.config.WorkerNum)
	group.ReceivedBlockGroupHeader = false
	for i, _ := range group.CheckTransactions {
		group.CheckTransactions[i] = 0
	}
	for i, _ := range group.CheckHeader {
		group.CheckHeader[i] = 0
	}
	return group
}

//区块组提交
func (node *Node) Commit() {
	node.RWLock.RLock()
	height := node.mongo.GetHeight() + 1
	value, ok := node.BlockGroups.Load(height)
	node.RWLock.RUnlock()
	flag := true
	if ok {
		item := value.(MetaData.BlockGroup)
		if item.ReceivedBlockGroupHeader {
			for k, v := range item.VoteResult {
				if v == 1 {
					if !item.Blocks[k].IsSet {
						flag = false
						break
					} else {
						block := item.Blocks[k]
						data := block.GetBlockHeaderBytes()
						if item.BlockHashes[k] != keymanager.GetHash(data) {
							item.Blocks[k].IsSet = false
							flag = false
							break
						}
					}
				}
			}
			if flag {
				node.mongo.PushbackBlockToDatabase(item) //数据落盘
				//同步区块时所有区块组都通过Commit函数执行提交操作，需要对创世区块组进行特殊处理
				if height == 0 {
					node.UpdateGenesisBlockVariables(&item)
					if len(node.state) == 1 {
						<-node.state
					}
				} else {
					if node.config.SendMsgToBCMgmt {
						node.SendNodeStatusMsgToBCManagementServer()
					}
					node.UpdateVariables(&item)
					fmt.Println("高度为", height, "的区块组成功共识并保存")

					//go node.checkCert()
				} //更新变量
				node.NormalGenerateVoteDone = false
				node.NormalGenerateBlockDone = false
				node.NormalGenerateBlockGroupHeaderDone = false
				node.isTimeout2 = false
				node.round = 0
			}
		} else {
			//fmt.Println("ReceivedBlockGroupHeader false")
			flag = false
		}
	}
	if !flag {
		if node.round > 0 && node.isTimeOut {
			//fmt.Println("执行超时请求操作")
			node.NormalTimeOutProcess()
			node.isTimeOut = false
		}
	}
}

func (node *Node) CheckBlocksHeader() {
	height := node.mongo.GetHeight() + 1
	value, _ := node.BlockGroups.LoadOrStore(height, node.CreateBlockGroup())
	item := value.(MetaData.BlockGroup)
	for i, value := range item.CheckTransactions {
		if value != 0 && item.CheckHeader[i] == 0 {
			if node.ValidateBlockHeader(&item.Blocks[i]) {
				item.CheckHeader[i] = 1
			} else {
				item.CheckHeader[i] = -1
			}
		}
	}
	node.BlockGroups.Store(height, item)
}

// 投票节点对区块投票并发送投票结果给轮值记账节点
func (node *Node) GenerateVote() {
	//fmt.Println("GenerateVote start")
	//fmt.Println("node.NormalGenerateVoteDone", node.NormalGenerateVoteDone)
	if !node.NormalGenerateVoteDone {
		height := node.mongo.GetHeight() + 1
		//检查是否存在当前共识所需高度的区块组
		value, existed := node.BlockGroups.Load(height)
		if !existed {
			return
		}
		item := value.(MetaData.BlockGroup)
		//设置投票结果
		var checkResult = make([]int, node.config.WorkerNum)
		for i := 0; i < len(node.accountManager.WorkerNumberSet); i++ {
			if item.CheckTransactions[i] == 0 || item.CheckHeader[i] == 0 {
				checkResult[i] = 0
			} else {
				if item.CheckTransactions[i] == 1 && item.CheckHeader[i] == 1 {
					checkResult[i] = 1
				} else {
					checkResult[i] = -1
				}
			}
		}
		//在未超时时需要检查是否所有投票已经产生，超时后不需要
		if node.round == 0 {
			if len(checkResult) < len(node.accountManager.WorkerNumberSet) {
				return
			}
			for _, value := range checkResult {
				if value == 0 {
					return
				}
			}
		} else {
			fmt.Println("超时投票！")
		}
		//设置各个区块hash值
		var hashes = make([]string, len(node.accountManager.WorkerNumberSet))
		for i, value := range checkResult {
			if value != 0 {
				block := item.Blocks[i]
				data := block.GetBlockHeaderBytes()
				hashes[i] = keymanager.GetHash(data)
			}
		}
		var ticket MetaData.VoteTicket
		ticket.VoteResult = checkResult
		ticket.BlockHashes = hashes
		ticket.Timestamp = utils.GetCurrentTime()
		ticket.Voter = node.config.MyPubkey
		ticket.Sig = ""
		data, _ := ticket.MarshalMsg(nil)
		ticket.Sig, _ = node.keyManager.Sign(data)
		pubkey := node.accountManager.WorkerNumberSet[node.true_dutyWorkerNum]
		var receiver uint64 = node.accountManager.WorkerSet[pubkey]
		var block_num = node.GetMyVoterNumber()
		header, msg := node.msgManager.CreateNormalBlocksVoteMsg(ticket, receiver, height, block_num)
		node.SendMessage(header, &msg)
		node.NormalGenerateVoteDone = true
		//fmt.Println("GenerateVote end")
	}
}

// 记账节点生成区块并发布
func (node *Node) GenerateBlock() {
	//fmt.Println("GenerateBlock start")
	//fmt.Println("node.NormalGenerateBlockDone", node.NormalGenerateBlockDone)
	if _, ok := node.accountManager.WorkerSet[node.config.MyPubkey]; ok {
		if !node.NormalGenerateBlockDone {

			if utils.GetCurrentTimeMilli()-node.StartTime*1e3 > node.config.GenerateBlockPeriod*1e3 {

				var block MetaData.Block
				block.Height = node.mongo.GetHeight() + 1
				pubkey := node.keyManager.GetPubkey()
				var block_num = node.GetMyWorkerNumber()
				block.BlockNum = block_num
				block.Generator = pubkey
				//设置前一区块组hash值
				//设置
				block.Transactions = node.txPool.GetCurrentTxsListDelete()
				headerBytes, _ := node.mongo.Block.ToHeaderBytes(nil)
				block.PreviousHash = keymanager.GetHash(headerBytes)
				block.MerkleRoot = keymanager.GetHash(block.GetTransactionsBytes())
				block.Timestamp = utils.GetCurrentTime()
				block.Sig, _ = node.keyManager.Sign([]byte(block.MerkleRoot))

				header, msg := node.msgManager.CreatePublishBlockMsg(block, 0)
				node.SendMessage(header, &msg)
				node.NormalGenerateBlockDone = true
			}
		}
	}
}

func (node *Node) GetMyWorkerNumber() uint32 {
	pubkey := node.keyManager.GetPubkey()
	var block_num uint32
	var find = false
	for block_num = 0; block_num < uint32(len(node.accountManager.WorkerNumberSet)); block_num++ {
		key, ok := node.accountManager.WorkerNumberSet[block_num]
		if ok && key == pubkey {
			find = true
			break
		}
	}
	if !find {
		fmt.Println(node.accountManager.WorkerNumberSet)
		fmt.Println(node.accountManager.WorkerSet)
		fmt.Println("GenerateBlock--找不到记账节点编号")
	}
	return block_num
}

func (node *Node) GetMyVoterNumber() uint32 {
	pubkey := node.config.MyPubkey
	var block_num uint32
	var find = false
	for block_num = 0; block_num < uint32(len(node.accountManager.VoterNumberSet)); block_num++ {
		key, ok := node.accountManager.VoterNumberSet[block_num]
		if ok && key == pubkey {
			find = true
			break
		}
	}
	if !find {
		fmt.Println("GenerateBlock--找不到投票节点编号")
	}
	return block_num
}

// 轮值记账节点生成区块组头部并发布
func (node *Node) GenerateBlockGroupHeader() {
	if !node.NormalGenerateBlockGroupHeaderDone {
		value, ok := node.BlockGroups.Load(node.mongo.GetHeight() + 1)
		if ok {
			item := value.(MetaData.BlockGroup)
			var count int = 0
			for _, y := range item.VoteTickets {
				if y.BlockHashes != nil && y.VoteResult != nil {
					count += 1
				}
			}
			if !node.isTimeout2 && count < len(node.accountManager.VoterSet) {
				return
			}
			new_item, is_complete := node.VotingStatistics(item)
			if !is_complete {
				return
			}

			new_item.NextDutyWorker = (node.dutyWorkerNumber + 1) % uint32(node.config.WorkerNum)
			new_item.Height = node.mongo.GetHeight() + 1
			new_item.Generator = node.config.MyPubkey
			headerBytes, _ := node.mongo.Block.ToHeaderBytes(nil)
			new_item.PreviousHash = keymanager.GetHash(headerBytes)
			new_item.Timestamp = utils.GetCurrentTime()
			var temp_header MetaData.BlockGroup
			temp_header = new_item
			temp_header.VoteTickets = nil
			tempHeaderBytes, _ := temp_header.ToHeaderBytes(nil)
			new_item.Sig, _ = node.keyManager.Sign(tempHeaderBytes)

			if node.config.SendMsgToBCMgmt {
				node.SendTransactionMsgToManagementServer(new_item) //Front end
			}

			var blockgroupheadermsg Message.BlockGroupHeader
			blockgroupheadermsg.Data, _ = new_item.ToHeaderBytes(nil)

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = node.network.MyNodeInfo.ID
			msgheader.Receiver = 0
			msgheader.Pubkey = node.config.MyPubkey
			msgheader.MsgType = Message.BlockGroupHeaderMsg
			node.SendMessage(msgheader, &blockgroupheadermsg)
			node.NormalGenerateBlockGroupHeaderDone = true
			//fmt.Println("GenerateBlockGroupHeader",new_item.Height,"done!")
		}
	}
}

// 移除节点请求
func (node *Node) RequestRemoveNode(pubkey string, nodeId uint64) {
	fmt.Println("Remove node")
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.RequestRemoveNode
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.RequestRemoveNodeMsg //消息体
	gm.Pubkey = pubkey
	gm.Type = "IamBack"
	gm.NodeID = nodeId
	node.SendMessage(mh, &gm)
}

func (node *Node) SendRequestQuitNode() {
	fmt.Println("节点发生退出交易请求")
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.RequestQuitNode
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.RequestQuitNodeMsg //消息体
	gm.Type = "IamBack"
	gm.Pubkey = node.config.MyPubkey
	gm.NodeID = node.network.MyNodeInfo.ID
	node.SendMessage(mh, &gm)
}

func (node *Node) SendMessage(header Message.MessageHeader, messageInterface Message.MessageInterface) {
	data, err := messageInterface.ToByteArray()
	if err != nil {
		fmt.Println("SendMessage：messageInterface.ToByteArray()错误！")
	}
	header.Data = data
	b, err := header.MarshalMsg(nil)
	if err != nil {
		fmt.Println("SendMessage：header.MarshalMsg()错误！")
	}
	ciphertext, err := keymanager.SM4Encrypt("x385HfQqGCYWlb4W", b)
	if err != nil {
		fmt.Println(err)
		fmt.Println("SM4加密错误！")
		return
	}
	rawbyte := base64.RawURLEncoding.EncodeToString(ciphertext)
	node.network.SendMessage([]byte(rawbyte), header.Receiver)
}

// HandleBCMessage 区块链节点之间通信处理函数
func (node *Node) HandleBCMessage(rawstr []byte) {
	// base64解码 base64可以把所有字符转换为可见字符 防止传输过程中解码错误
	ciphertext, err := base64.RawURLEncoding.DecodeString(string(rawstr))
	if err != nil {
		fmt.Println("base64解密错误！err:", err)
		return
	}
	// 国密算法解码
	data, err := keymanager.SM4Decrypt("x385HfQqGCYWlb4W", ciphertext)
	if err != nil {
		fmt.Println("SM4解密错误！err:", err)
		return
	}
	var header Message.MessageHeader
	// 解码data到header
	data, err = header.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("HandleBCMessage:header.UnmarshalMsg(data)错误！")
		return
	}


	data = header.Data
	// 跟据类型解码数据，并做出对应的处理
	switch header.MsgType {
	case Message.Zero:
		var msg Message.ZeroMsg
		data, err = msg.UnmarshalMsg(data)
	case Message.QueryPubkey:
		var msg Message.QueryPubkeyMsg
		data, err = msg.UnmarshalMsg(data)
		node.HandleQueryPubMessage(msg, header.Sender)
	case Message.GenesisBlock:
		var msg Message.BlockMsg
		data, err = msg.UnmarshalMsg(data)
		node.HandleGenesisBlockPublishMessage(msg, header)
	case Message.TransactionMsg:
		var msg Message.TransactionMessage
		data, err = msg.UnmarshalMsg(data)
		node.HandleTransactionMessage(msg, header)
	case Message.NormalPublishBlock:
		var msg Message.PublishBlockMsg
		err = msg.FromByteArray(data)
		node.HandlePublishBlockMessage(msg, header)
	case Message.NormalBlockVoteMsg:
		var msg Message.NormalBlocksVoteMsg
		err = msg.FromByteArray(data)
		node.HandleVoteTicketMessage(msg, header)
	case Message.RequestHeight:
		var msg Message.RequestHeightMsg
		err = msg.FromByteArray(data)
		node.HandleRequestHeightMessage(msg, header)
	case Message.RespondHeight:
		var msg Message.RequestHeightMsg
		err = msg.FromByteArray(data)
		node.HandleRespondHeightMessage(msg, header)
	case Message.RequestBlockGroup:
		var msg Message.RequestBlockGroupMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockGroupMessage(msg, header)
	case Message.RespondBlockGroup:
		var msg Message.RespondBlockGroupMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockGroupMessage(msg, header)
	case Message.RequestBlockGroupHeader:
		var msg Message.RequestBlockGroupHeaderMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockGroupHeaderMessage(msg, header)
	case Message.ResponseBlockGroupHeader:
		var msg Message.RespondBlockGroupHeaderMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockGroupHeaderMessage(msg, header)
	case Message.RequestBlock:
		var msg Message.RequestBlockMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockMessage(msg, header)
	case Message.ResponseBlock:
		var msg Message.RespondBlockMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockMessage(msg, header)
	case Message.ElectNewWorker:
		var msg Message.ElectNewWorkerMsg
		err = msg.FromByteArray(data)
		node.HandleElectNewWorkerMessage(msg, header)
	case Message.BlockGroupHeaderMsg:
		var msg Message.BlockGroupHeader
		err = msg.FromByteArray(data)
		node.HandleBlockGroupHeaderMessage(msg, header)
	case Message.RequestAddNewNode:
		var msg Message.RequestAddNewNodeMsg
		err = msg.FromByteArray(data)
		node.HandleAddNewNodeRequest(msg, header)
	case Message.QueryAddNodeVote:
		var msg Message.QueryAddNodeVoteMsg
		err = msg.FromByteArray(data)
		node.HandleQueryAddNewNodeVote(msg, header)
	case Message.ResponseQueryAddNodeVote:
		var msg Message.QueryAddNodeVoteMsg
		err = msg.FromByteArray(data)
		node.HandleResponseQueryAddNodeVote(msg, header)
	case Message.PublishIamOk:
		var msg Message.PublishIamOkMsg
		err = msg.FromByteArray(data)
		node.HandlePublishIamOkMsg(msg, header)
	case Message.RequestQuitNode:
		var msg Message.RequestQuitNodeMsg
		err = msg.FromByteArray(data)
		node.HandleRequestQuitNodeMsg(msg, header)
	case Message.RequestRemoveNode:
		var msg Message.RequestRemoveNodeMsg
		err = msg.FromByteArray(data)
		node.HandleRequestRemoveNode(msg, header)
	case Message.ResponseRemoveNode:
		var msg Message.ResponseRemoveNodeMsg
		err = msg.FromByteArray(data)
		node.HandleResponseRemoveNode(msg, header)
	case Message.NodeStatus:
		var msg Message.NodeStatusMsg
		err = msg.FromByteArray(data)
		node.HandleNodeStatusMsgFromBC(msg, header)
	case Message.TransactionStatistics:
		var msg Message.TransactionStatisticsMsgs
		err = msg.FromByteArray(data)
		node.HandleTransactionMsgFromBC(msg, header)
	}
}

func (node *Node) HandleQueryPubMessage(msg Message.QueryPubkeyMsg, pre_sender Network.NodeID) {
	switch msg.Type {
	case 100:
		// 回送QueryPubKey信息 类型101
		var mh Message.MessageHeader
		mh.MsgType = Message.QueryPubkey
		mh.Sender = node.network.MyNodeInfo.ID
		mh.Receiver = pre_sender
		var gm Message.QueryPubkeyMsg
		gm.Type = 101
		gm.Information = node.config.MyPubkey
		gm.NodeID = node.network.MyNodeInfo.ID
		node.SendMessage(mh, &gm)

	case 101:
		node.RWLock.Lock()
		temp, ok := node.network.NodeList[pre_sender]
		if !ok {
			fmt.Println("sender not exist")
			node.RWLock.Unlock()
			return
		}
		for _, x := range node.config.WorkerList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.WorkerPubList[msg.Information]
				if !ok {
					node.WorkerPubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		for _, x := range node.config.WorkerCandidateList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.WorkerCandidatePubList[msg.Information]
				if !ok {
					node.WorkerCandidatePubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		for _, x := range node.config.VoterList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.VoterPubList[msg.Information]
				if !ok {
					node.VoterPubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		node.RWLock.Unlock()
	}
}

func (node *Node) HandleGenesisBlockPublishMessage(msg Message.BlockMsg, header Message.MessageHeader) {
	var bg MetaData.BlockGroup
	_, err := bg.FromBytes(msg.Data)

	if err != nil || utils.GetCurrentTime()-bg.Timestamp > 30 || utils.GetCurrentTime()-bg.Timestamp < -30 {
		log.Println("时间戳错误")
		return
	}

	sig := bg.Sig
	bg.Sig = ""
	temp, _ := bg.ToHeaderBytes(nil)
	ok, _ := node.keyManager.Verify(temp, sig, bg.Generator)
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		node.mongo.PushbackBlockToDatabase(bg)
		node.UpdateGenesisBlockVariables(&bg)
	}
}

func (node *Node) HandleTransactionMessage(msg Message.TransactionMessage, header Message.MessageHeader) {
	head, inter := MetaData.DecodeTransaction(msg.Data)
	node.txPool.PushbackTransaction(head, inter)
}

func (node *Node) HandlePublishBlockMessage(msg Message.PublishBlockMsg, header Message.MessageHeader) {
	height := msg.GetHeight()
	blockNum := msg.GetBlockNum()
	if height <= node.mongo.GetHeight() {
		return
	}
	//增加稳定性，防止程序崩溃
	if int(blockNum) >= node.config.WorkerNum {
		return
	}
	value, _ := node.BlockGroups.LoadOrStore(height, node.CreateBlockGroup())
	block := msg.GetBlock()

	if utils.GetCurrentTime()-block.Timestamp > 30 || utils.GetCurrentTime()-block.Timestamp < -30 {
		log.Println("时间戳错误")
		return
	}

	merkleRoot := keymanager.GetHash(block.GetTransactionsBytes())
	if merkleRoot != block.MerkleRoot {
		fmt.Println("MerkleRoot错误")
		return
	}
	ok, _ := node.keyManager.Verify([]byte(block.MerkleRoot), block.Sig, block.Generator)
	if ok {
		block.IsSet = true

		item := value.(MetaData.BlockGroup)
		item.Blocks[blockNum] = block
		if !node.ValidateTransactions(&block.Transactions) {
			item.CheckTransactions[blockNum] = -1

		} else {
			item.CheckTransactions[blockNum] = 1
		}

		node.BlockGroups.Store(height, item)
	}
}

func (node *Node) HandleVoteTicketMessage(msg Message.NormalBlocksVoteMsg, header Message.MessageHeader) {
	height := msg.Height
	blockNum := msg.BlockNum

	if height <= node.mongo.GetHeight() {
		return
	}
	//增加稳定性，防止程序崩溃
	if int(blockNum) >= len(node.accountManager.VoterSet) {
		return
	}
	//如果不存在BlockGroup，则创建一个
	value, _ := node.BlockGroups.LoadOrStore(height, node.CreateBlockGroup())

	if utils.GetCurrentTime()-msg.Ticket.Timestamp > 30 || utils.GetCurrentTime()-msg.Ticket.Timestamp < -30 {
		log.Println("时间戳错误")
		return
	}

	item := value.(MetaData.BlockGroup)
	if uint32(len(item.VoteTickets)) <= blockNum {
		tmpTickets := item.VoteTickets
		item.VoteTickets = make([]MetaData.VoteTicket, blockNum+1)
		for k, v := range tmpTickets {
			item.VoteTickets[k] = v
		}

		results := item.VoteResult
		item.VoteResult = make([]int, blockNum+1)
		for k, v := range results {
			item.VoteResult[k] = v
		}

	}
	item.VoteTickets[blockNum] = msg.Ticket
	node.BlockGroups.Store(height, item)
}

func (node *Node) HandleRequestHeightMessage(msg Message.RequestHeightMsg, header Message.MessageHeader) {
	header, msg = node.msgManager.CreateRespondHeightMsg(header.Sender, node.mongo.GetHeight())
	node.SendMessage(header, &msg)
}

func (node *Node) HandleRespondHeightMessage(msg Message.RequestHeightMsg, header Message.MessageHeader) {
	height := msg.Height
	if height > node.mongo.GetHeight() {
		var s Syncer
		s.NodeID = header.Sender
		s.Height = height
		node.syncTool.Syncers <- s
		fmt.Println("收到height=", s.Height, "NodeID=", s.NodeID, "的高度回复")
	}

}

func (node *Node) HandleRequestBlockGroupMessage(msg Message.RequestBlockGroupMsg, header Message.MessageHeader) {
	height := msg.Height
	if node.mongo.GetHeight() < height {
		return
	}
	group := node.mongo.GetBlockFromDatabase(height)
	header, response := node.msgManager.CreateRespondBlockGroupMsg(header.Sender, height, group)
	node.SendMessage(header, &response)
}

func (node *Node) HandleRespondBlockGroupMessage(msg Message.RespondBlockGroupMsg, header Message.MessageHeader) {
	if len(msg.Group.Blocks) <= 0 {
		return
	}
	fmt.Println("接收到区块组回复")
	msg.Group.ReceivedBlockGroupHeader = true
	msg.Group.CheckTransactions = make([]int, node.config.WorkerNum)
	msg.Group.CheckHeader = make([]int, node.config.WorkerNum)
	for k, v := range msg.Group.VoteResult {
		if v == 1 {
			msg.Group.CheckTransactions[k] = 1
			msg.Group.CheckHeader[k] = 1
			msg.Group.Blocks[k].IsSet = true
		}
	}
	node.RWLock.Lock()
	node.BlockGroups.Store(msg.Group.Height, msg.Group)
	node.RWLock.Unlock()
}

func (node *Node) HandleElectNewWorkerMessage(msg Message.ElectNewWorkerMsg, header Message.MessageHeader) {
	node.ElectNewWorkerList = append(node.ElectNewWorkerList, msg)
}

func (node *Node) HandleBlockGroupHeaderMessage(msg Message.BlockGroupHeader, header Message.MessageHeader) {
	var blockgroup_header, newBGHeader MetaData.BlockGroup
	_, err := blockgroup_header.FromHeaderBytes(msg.Data)

	if err != nil || utils.GetCurrentTime()-blockgroup_header.Timestamp > 30 || utils.GetCurrentTime()-blockgroup_header.Timestamp < -30 {
		log.Println("时间戳错误")
		return
	}

	newBGHeader = blockgroup_header
	newBGHeader.VoteTickets = nil
	newBGHeader.Sig = ""
	temp, _ := newBGHeader.ToHeaderBytes(nil)
	ok, err := node.keyManager.Verify(temp, blockgroup_header.Sig, blockgroup_header.Generator)
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		if blockgroup_header.Height >= node.mongo.Height+1 {
			value, _ := node.BlockGroups.LoadOrStore(blockgroup_header.Height, node.CreateBlockGroup())

			item := value.(MetaData.BlockGroup)
			item.Height = blockgroup_header.Height
			item.Generator = blockgroup_header.Generator
			item.PreviousHash = blockgroup_header.PreviousHash
			item.MerkleRoot = blockgroup_header.MerkleRoot
			item.Timestamp = blockgroup_header.Timestamp
			item.Sig = blockgroup_header.Sig
			item.NextDutyWorker = blockgroup_header.NextDutyWorker
			item.BlockHashes = blockgroup_header.BlockHashes
			item.VoteResult = blockgroup_header.VoteResult
			item.VoteTickets = blockgroup_header.VoteTickets
			item.ReceivedBlockGroupHeader = true
			//fmt.Println("HandleBlockGroupHeaderMessage change variable", blockgroup_header.Height)

			node.BlockGroups.Store(blockgroup_header.Height, item)
		}
	}
}

func (node *Node) HandleRequestBlockGroupHeaderMessage(msg Message.RequestBlockGroupHeaderMsg, header Message.MessageHeader) {
	//fmt.Println("接收到高度为",msg.Height,"的区块组头请求")
	if msg.Height >= node.mongo.GetHeight() {
		value, ok := node.BlockGroups.Load(msg.Height)
		if ok {
			group := value.(MetaData.BlockGroup)
			if group.ReceivedBlockGroupHeader {
				fmt.Println(node.network.MyNodeInfo.ID, "在内存中找到并发送高度为", msg.Height, "的区块组头")
				RespHeader, RespMsg := node.msgManager.CreateRespondBlockGroupHeaderMsg(header.Sender, msg.Height, &group)
				node.SendMessage(RespHeader, &RespMsg)
			}
		} else {
			fmt.Println(node.network.MyNodeInfo.ID, "在内存中找不到高度为", msg.Height, "的区块组")
		}

	} else {
		group := node.mongo.GetBlockFromDatabase(msg.Height)
		RespHeader, RespMsg := node.msgManager.CreateRespondBlockGroupHeaderMsg(header.Sender, msg.Height, &group)
		fmt.Println(node.network.MyNodeInfo.ID, "在数据库中找到并发送高度为", msg.Height, "的区块组头")
		node.SendMessage(RespHeader, &RespMsg)
	}
}

func (node *Node) HandleRespondBlockGroupHeaderMessage(msg Message.RespondBlockGroupHeaderMsg, header Message.MessageHeader) {
	var blockgroup_header MetaData.BlockGroup
	blockgroup_header.FromHeaderBytes(msg.BlockGroupHeaderBytes)
	newBGHeader := blockgroup_header
	newBGHeader.VoteTickets = nil
	newBGHeader.Sig = ""
	temp, _ := newBGHeader.ToHeaderBytes(nil)
	ok, _ := node.keyManager.Verify(temp, blockgroup_header.Sig, blockgroup_header.Generator)

	if ok {
		if blockgroup_header.Height >= node.mongo.Height+1 {
			value, _ := node.BlockGroups.LoadOrStore(blockgroup_header.Height, node.CreateBlockGroup())

			item := value.(MetaData.BlockGroup)
			item.Height = blockgroup_header.Height
			item.Generator = blockgroup_header.Generator
			item.PreviousHash = blockgroup_header.PreviousHash
			item.MerkleRoot = blockgroup_header.MerkleRoot
			item.Timestamp = blockgroup_header.Timestamp
			item.Sig = blockgroup_header.Sig
			item.NextDutyWorker = blockgroup_header.NextDutyWorker
			item.BlockHashes = blockgroup_header.BlockHashes
			item.VoteResult = blockgroup_header.VoteResult
			item.VoteTickets = blockgroup_header.VoteTickets
			item.ReceivedBlockGroupHeader = true

			node.BlockGroups.Store(blockgroup_header.Height, item)
		}
	}
}

func (node *Node) HandleRequestBlockMessage(msg Message.RequestBlockMsg, header Message.MessageHeader) {
	if msg.Height >= node.mongo.GetHeight()+1 {
		value, ok := node.BlockGroups.Load(msg.Height)
		if ok {
			group := value.(MetaData.BlockGroup)
			if group.Blocks[msg.BlockNum].IsSet {
				RespHeader, RespMsg := node.msgManager.CreateRequestBlockMsg(header.Sender, msg.Height, msg.BlockNum)
				node.SendMessage(RespHeader, &RespMsg)
			}
		}
	} else {
		group := node.mongo.GetBlockFromDatabase(msg.Height)

		RespHeader, RespMsg := node.msgManager.CreateRespondBlockMsg(header.Sender, msg.Height, msg.BlockNum, group.Blocks[msg.BlockNum])
		node.SendMessage(RespHeader, &RespMsg)
	}
}

func (node *Node) HandleRespondBlockMessage(msg Message.RespondBlockMsg, header Message.MessageHeader) {
	height := msg.Height
	blockNum := msg.BlockNum
	if height <= node.mongo.GetHeight() {
		return
	}
	//增加稳定性，防止程序崩溃
	if int(blockNum) >= node.config.WorkerNum {
		return
	}
	value, _ := node.BlockGroups.LoadOrStore(height, node.CreateBlockGroup())

	block := msg.Block
	block.IsSet = true
	item := value.(MetaData.BlockGroup)
	item.Blocks[blockNum] = block
	if !node.ValidateTransactions(&block.Transactions) {
		item.CheckTransactions[blockNum] = -1
	} else {
		item.CheckTransactions[blockNum] = 1
	}
	node.BlockGroups.Store(height, item)
}


func (node *Node) HandleAddNewNodeRequest(msg Message.RequestAddNewNodeMsg, header Message.MessageHeader) {
	node.AddNodeLock.Lock()
	defer node.AddNodeLock.Unlock()

	fmt.Println("接收到新的节点加入请求")
	_, ok := node.accountManager.VoterSet[msg.Pubkey]
	if ok {
		fmt.Println("申请成为投票节点失败，已经是投票节点")
		return
	}
	_, ok = node.accountManager.WorkerCandidateSet[msg.Pubkey]
	if ok {
		fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
		return
	}

	if node.mongo.HasData("identity", "pubkey", msg.Pubkey) {
		fmt.Println("申请失败,已经查询到该请求")
		return
	}
	var transaction MetaData.IdentityTransformation
	transaction.Type = msg.Type
	transaction.Port = msg.Port
	transaction.Pubkey = msg.Pubkey
	transaction.SetNodeId(msg.NodeID)
	transaction.IPAddr = msg.IPAddr

	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.IdTransformation
	node.txPool.PushbackTransaction(transactionHeader, &transaction)

}

func (node *Node) HandleQueryAddNewNodeVote(msg Message.QueryAddNodeVoteMsg, header Message.MessageHeader) {
	if !node.mongo.HasData("identity", "pubkey", msg.Pubkey) {
		return
	}

	fmt.Println("接收到新的加入投票查询请求")
	var mh Message.MessageHeader
	mh.MsgType = Message.ResponseQueryAddNodeVote
	mh.Receiver = header.Sender
	mh.Sender = node.network.MyNodeInfo.ID

	var gm Message.QueryAddNodeVoteMsg
	gm.Type = msg.Type
	gm.Result = 1
	gm.Pubkey = node.config.MyPubkey
	gm.NodeID = msg.NodeID
	gm.Sign = ""

	node.SendMessage(mh, &gm)
}

func (node *Node) HandleResponseQueryAddNodeVote(msg Message.QueryAddNodeVoteMsg, header Message.MessageHeader) {
	fmt.Println("接收到新的加入投票请求回复", msg.Pubkey, msg.Result)
	node.AddNodeLock.Lock()
	node.AddNewNodeVoteList[msg.Pubkey] = msg.Result
	defer node.AddNodeLock.Unlock()
}

func (node *Node) HandlePublishIamOkMsg(msg Message.PublishIamOkMsg, header Message.MessageHeader) {
	fmt.Println("接收到新节点同步完成消息")
	if !node.mongo.HasData("identity", "pubkey", msg.Pubkey) {
		return
	}
	var transaction MetaData.IdentityTransformation
	transaction.Type = msg.Type
	transaction.Pubkey = msg.Pubkey

	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.IdTransformation
	node.txPool.PushbackTransaction(transactionHeader, &transaction)

}

func (node *Node) HandleRequestQuitNodeMsg(msg Message.RequestQuitNodeMsg, header Message.MessageHeader) {
	fmt.Println("接收到新的节点退出请求")
	if node.mongo.HasData("identity", "pubkey", msg.Pubkey) {
		id := node.mongo.GetOneNodeIdentityTransFromDatabase("identity", "pubkey", msg.Pubkey)
		if id.Type == msg.Type {
			fmt.Println("已经收到该节点的退出请求")
			return
		}
	}

	_, ok := node.accountManager.VoterSet[msg.Pubkey]
	if !ok {
		fmt.Println("申请退出投票节点失败，不是投票节点")
		return
	}
	_, ok = node.accountManager.WorkerCandidateSet[msg.Pubkey]
	if !ok {
		fmt.Println("申请退出候选记账节点失败，不是候选记账节点")
		return
	}

	var transaction MetaData.IdentityTransformation
	transaction.Type = msg.Type
	transaction.Pubkey = msg.Pubkey
	transaction.SetNodeId(msg.NodeID)

	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.IdTransformation
	node.txPool.PushbackTransaction(transactionHeader, &transaction)
}

func (node *Node) HandleRequestRemoveNode(msg Message.RequestRemoveNodeMsg, header Message.MessageHeader) {
	fmt.Println("接收到删除用户请求")
	if node.mongo.HasData("identity", "pubkey", msg.Pubkey) {
		id := node.mongo.GetOneNodeIdentityTransFromDatabase("identity", "pubkey", msg.Pubkey)
		if id.Type == msg.Type {
			fmt.Println("已经收到该节点的退出请求")
			return
		}
	}

	_, ok := node.accountManager.VoterSet[msg.Pubkey]
	if !ok {
		fmt.Println("申请退出投票节点失败，不是投票节点")
		return
	}
	_, ok = node.accountManager.WorkerCandidateSet[msg.Pubkey]
	if !ok {
		fmt.Println("申请退出候选记账节点失败，不是候选记账节点")
		return
	}

	var mh Message.MessageHeader
	mh.MsgType = Message.ResponseRemoveNode
	mh.Receiver = header.Sender
	mh.Sender = node.network.MyNodeInfo.ID

	var gm Message.ResponseRemoveNodeMsg
	gm.Type = msg.Type
	gm.Result = 1
	gm.Pubkey = msg.Pubkey
	gm.NodeID = msg.NodeID
	gm.Sign = ""
	node.SendMessage(mh, &gm)
}

func (node *Node) HandleResponseRemoveNode(msg Message.ResponseRemoveNodeMsg, header Message.MessageHeader) {
	fmt.Println("接收到删除用户请求的回复")
	node.AddNodeLock.Lock()
	node.RemoveNodeVoteList[header.Sender] = msg.Result
	defer node.AddNodeLock.Unlock()
	fmt.Println(node.RemoveNodeVoteList)

	var agree int = 0
	for _, v := range node.RemoveNodeVoteList {
		if v == 1 {
			agree++
		}
	}

	if agree >= len(node.accountManager.VoterSet)/2 {
		//fmt.Println(len(node.accountManager.VoterSet))
		//fmt.Println("进入同步区块阶段")
		fmt.Println("发出用户退出请求")
		var mh Message.MessageHeader //消息头
		mh.MsgType = Message.RequestQuitNode
		mh.Sender = msg.NodeID
		mh.Receiver = 0
		var gm Message.RequestQuitNodeMsg //消息体
		gm.Type = "IamBack"
		gm.Pubkey = msg.Pubkey
		gm.NodeID = msg.NodeID
		node.SendMessage(mh, &gm)
	}
}