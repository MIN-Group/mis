package Node

import (
	"math/rand"
	"sort"
	"time"

	"MIS-BC/Message"
	"MIS-BC/MetaData"
)

type CountElectNewWorker struct {
	Pubkey  string
	VoteNum int
}

type CountElectNewWorkerSlice []CountElectNewWorker

func (a CountElectNewWorkerSlice) Len() int { // 重写 Len() 方法
	return len(a)
}

func (a CountElectNewWorkerSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (a CountElectNewWorkerSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].VoteNum < a[i].VoteNum
}

func (node *Node) ElectNewWorker() {
	if node.mongo.Height%node.config.BlockGroupPerCycle == 0 {
		rand.Seed(time.Now().Unix())
		var electNewWorkerMsg Message.ElectNewWorkerMsg
		electNewWorkerMsg.MyPubkey = node.config.MyPubkey
		for i := 0; i < node.config.WorkerNum; i++ {
			rand_id := rand.Intn(len(node.accountManager.WorkerCandidateList))
			electNewWorkerMsg.NewWorker = append(electNewWorkerMsg.NewWorker, node.accountManager.WorkerCandidateList[rand_id])
		}

		var header Message.MessageHeader
		header.Sender = node.network.MyNodeInfo.ID
		header.Receiver = node.accountManager.VoterSet[node.accountManager.WorkerNumberSet[node.dutyWorkerNumber]]
		header.Pubkey = node.config.MyPubkey
		header.MsgType = Message.ElectNewWorker

		node.SendMessage(header, &electNewWorkerMsg)
	}
}

func (node *Node) CountElectNewWorker() {
	node.ElectNewWorkerList = make([]Message.ElectNewWorkerMsg, 1000)
	if node.mongo.Height%node.config.BlockGroupPerCycle == 0 && node.accountManager.WorkerNumberSet[node.dutyWorkerNumber] == node.config.MyPubkey {
		for {
			if !(len(node.ElectNewWorkerList) == node.config.VotedNum) {
				continue
			}
			counter_map := make(map[string]int)
			for _, x := range node.ElectNewWorkerList {
				for _, y := range x.NewWorker {
					temp, ok := counter_map[y]
					if ok {
						counter_map[y] = temp + 1
					} else {
						counter_map[y] = 1
					}
				}
			}
			var result []CountElectNewWorker
			for k, v := range counter_map {
				result = append(result, CountElectNewWorker{k, v})
			}
			sort.Sort(CountElectNewWorkerSlice(result))

			var trans MetaData.ElectNewWorkerTeam
			trans.WorkerPubList = make(map[string]uint64)
			trans.ElectResult = counter_map
			for i := 0; i < node.config.WorkerNum; i++ {
				trans.WorkerPubList[result[i].Pubkey] = node.accountManager.WorkerCandidateSet[result[i].Pubkey]
			}

			var transactionHeader MetaData.TransactionHeader
			transactionHeader.TXType = MetaData.ElectNewWorker
			node.txPool.PushbackTransaction(transactionHeader, &trans)

			break
		}
	}
}
