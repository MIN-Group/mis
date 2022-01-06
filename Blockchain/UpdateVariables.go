package Node

import (
	"MIS-BC/MetaData"
	"MIS-BC/Network"
	"MIS-BC/common"
	"fmt"
	"time"
)

func (node *Node) UpdateVariables(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		node.StartTime = bg.Timestamp

		for i, eachBlock := range bg.Blocks {
			if len(bg.VoteResult) <= i {
				continue
			}
			if bg.VoteResult[i] != 1 {
				continue
			}
			for _, eachTransaction := range eachBlock.Transactions {
				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdentityAction:
					node.UpdateIdentityVariables(transactionInterface)
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				}
			}
		}
	} else {
		fmt.Println("更新变量错误")
	}

	node.UpdateIdTransOk()
}

func (node *Node) UpdateIdentityVariables(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.Identity); ok {
		switch transaction.Command {
		case "Registry":
			node.mongo.SaveIdentityToDatabase(*transaction)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "申请注册成功")
		case "DestroyByIdentityIdentifier":
			node.mongo.DeleteIdentity("identityidentifier", transaction.IdentityIdentifier)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "删除成功")
		case "ResetPassword":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.ResetIdentityPassword(transaction.Passwd, item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "密码修改成功")
		case "ResetIPIdentifier":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.ResetIdentityIPIdentifier(transaction.IPIdentifier, item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "IP更新成功")
		case "EnableIdentity":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.EnableIdentity(transaction.IsValid, transaction.Cert, item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "审核成功")
		case "ResetValidation":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.ResetIdentityValidation(transaction.IsValid, item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "有效性变更成功")
		case "CertRevocation":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.CertRevocation(item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "证书撤销成功")
		case "CertReissue":
			item := node.mongo.GetOneIdentityMapFromDatabase(transaction.IdentityIdentifier)
			node.mongo.CertReissue(transaction.Cert, item)
			common.Logger.Info("身份", transaction.IdentityIdentifier, "证书重新颁发成功")
		}
	}

}

func (node *Node) UpdateIdTransformationVaribles(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.IdentityTransformation); ok {
		switch transaction.Type {

		case "ApplyNode": //apply for voter and worker
			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if ok {
				fmt.Println("申请成为投票节点失败，已经是投票节点")
				return
			}
			_, ok = node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if ok {
				fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
				return
			}
			ok = node.mongo.HasData("identity", "pubkey", transaction.Pubkey)
			if ok {
				fmt.Println("已经接收到该用户请求")
				return
			}
			node.network.AddNodeToNodeList(transaction.GetNodeId(), transaction.IPAddr, transaction.Port)
			node.mongo.SaveNodeIdentityTransToDatabase(*transaction)
			fmt.Println(transaction.GetNodeId(), "申请成为投票节点成功")

		case "IamOk":
			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if ok {
				fmt.Println("申请成为投票节点失败，已经是投票节点")
				return
			}
			_, ok = node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if ok {
				fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
				return
			}
			if !node.mongo.HasData("identity", "pubkey", transaction.Pubkey) {
				fmt.Println("无法查到该节点的申请信息")
				return
			}

			identity := node.mongo.GetOneNodeIdentityTransFromDatabase("identity", "pubkey", transaction.Pubkey)
			identity.Type = "IamOk"
			identity.Timestamp = time.Now().Unix()

			/*_, ok = node.IdentityTransList[identity]
			if ok{
				fmt.Println("已经接收到该节点加入请求")
				return
			}*/

			node.IdentityTransList[identity] = node.mongo.GetHeight() + 3
			fmt.Println("新节点申请加入成功")
			var nodelist MetaData.NodeList
			nodelist.SetNodeList(node.network.NodeList)
			node.mongo.InsertOrUpdateNodeList(nodelist)

			var identityTransList MetaData.IdentityTransList
			identityTransList.SetIdentityTransList(node.IdentityTransList)
			node.mongo.InsertOrUpdateIdentityTransList(identityTransList)
		case "IamBack":
			if node.mongo.HasData("identity", "pubkey", transaction.Pubkey) {
				id := node.mongo.GetOneNodeIdentityTransFromDatabase("identity", "pubkey", transaction.Pubkey)
				if id.Type == transaction.Type {
					fmt.Println("已经收到该节点的退出请求")
					return
				}
			}

			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if !ok {
				fmt.Println("申请退出投票节点失败，不是投票节点")
				return
			}
			_, ok = node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if !ok {
				fmt.Println("申请退出候选记账节点失败，不是候选记账节点")
				return
			}
			node.mongo.SaveNodeIdentityTransToDatabase(*transaction)
			node.IdentityTransList[*transaction] = node.mongo.GetHeight() + 3
			fmt.Println("新节点申请退出成功")

		case "ApplyForVoter":
			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if !ok {
				node.accountManager.VoterSet[transaction.Pubkey] = transaction.GetNodeId()
			} else {
				fmt.Println("申请成为投票节点失败，已经是投票节点")
			}
			_, ok = node.network.NodeList[transaction.GetNodeId()]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.GetNodeId()
				node.network.NodeList[transaction.GetNodeId()] = nodelist
			}
		case "ApplyForWorkerCandidate":
			_, ok := node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if !ok {
				node.accountManager.WorkerCandidateSet[transaction.Pubkey] = transaction.GetNodeId()
			} else {
				fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
			}
			_, ok = node.network.NodeList[transaction.GetNodeId()]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.GetNodeId()
				node.network.NodeList[transaction.GetNodeId()] = nodelist
			}
		case "QuitVoter":
			delete(node.accountManager.VoterSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.GetNodeId())
			fmt.Println("退出投票节点成功")
		case "QuitWorkerCandidate":
			delete(node.accountManager.WorkerCandidateSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.GetNodeId())
			fmt.Println("退出候选记账节点成功")
		}
	}
}

func (node *Node) UpdateIdTransOk() {
	for k, v := range node.IdentityTransList {
		if v == node.mongo.GetHeight() {
			if k.Type == "IamOk" {
				_, ok := node.accountManager.VoterSet[k.Pubkey]
				flag := true
				if ok {
					flag = false
					fmt.Println(node.accountManager.VoterNumberSet)
					fmt.Println(node.accountManager.VoterSet)
					fmt.Println("申请成为投票节点失败，已经是投票节点")
				}
				_, ok = node.accountManager.WorkerCandidateSet[k.Pubkey]
				if ok {
					flag = false
					fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
				}
				if !node.mongo.HasData("identity", "pubkey", k.Pubkey) {
					flag = false
					fmt.Println("无法查到该节点的申请信息")
				}
				if flag {
					node.accountManager.VoterSet[k.Pubkey] = k.GetNodeId()

					var nums uint32 = 0
					for k, _ := range node.accountManager.VoterNumberSet {
						if nums < k {
							nums = k
						}
					}
					nums++
					node.accountManager.VoterNumberSet[nums] = k.Pubkey
					node.accountManager.WorkerCandidateSet[k.Pubkey] = k.GetNodeId()
					if len(node.accountManager.WorkerSet) < node.config.WorkerNum {
						var i = 0
						for ; i < node.config.WorkerNum; i++ {
							if _, ok := node.accountManager.WorkerNumberSet[uint32(i)]; !ok {
								node.accountManager.WorkerSet[k.Pubkey] = k.GetNodeId()
								node.accountManager.WorkerNumberSet[uint32(i)] = k.Pubkey
								break
							}
						}
					}

					var account MetaData.Account
					account.SetVoterSet(node.accountManager.VoterSet)
					account.SetWorkerSet(node.accountManager.WorkerSet)
					account.SetWorkerCandidateSet(node.accountManager.WorkerCandidateSet)
					account.WorkerCandidateList = node.accountManager.WorkerCandidateList
					account.SetVoterNumberSet(node.accountManager.VoterNumberSet)
					account.SetWorkerNumberSet(node.accountManager.WorkerNumberSet)
					node.mongo.InsertOrUpdateAccount(account)
					node.mongo.DeleteData("identity", "pubkey", k.Pubkey)
					fmt.Println("新节点列表更新成功")
				}
				delete(node.IdentityTransList, k)
			} else if k.Type == "IamBack" {
				flag := true
				_, ok := node.accountManager.VoterSet[k.Pubkey]
				if !ok {
					fmt.Println("申请退出投票节点失败，不是投票节点")
					flag = false
				}
				_, ok = node.accountManager.WorkerCandidateSet[k.Pubkey]
				if !ok {
					fmt.Println("申请退出候选记账节点失败，不是候选记账节点")
					flag = false
				}
				if !node.mongo.HasDataByTwoKey("identity", "pubkey", k.Pubkey, "type", k.Type) {
					fmt.Println("没有查到该节点的退出请求")
					flag = false
				}
				if flag {
					fmt.Println("transaction", k)
					fmt.Println(node.accountManager.WorkerNumberSet)
					delete(node.accountManager.WorkerSet, k.Pubkey)
					delete(node.accountManager.WorkerCandidateSet, k.Pubkey)
					delete(node.accountManager.VoterSet, k.Pubkey)

					var voterNum uint32 = 0
					var max uint32 = 0
					for k1, v1 := range node.accountManager.VoterNumberSet {
						if v1 == k.Pubkey {
							voterNum = k1
							delete(node.accountManager.VoterNumberSet, k1)
						}
						if k1 > max {
							max = k1
						}
					}
					if voterNum != max {
						node.accountManager.VoterNumberSet[voterNum] = node.accountManager.VoterNumberSet[max]
						delete(node.accountManager.VoterNumberSet, max)
					}

					var workerNum uint32 = 0
					var max1 uint32 = uint32(len(node.accountManager.WorkerNumberSet)) - 1
					isDel := false
					for k1, v1 := range node.accountManager.WorkerNumberSet {
						if v1 == k.Pubkey {
							workerNum = k1
							delete(node.accountManager.WorkerNumberSet, k1)
							isDel = true
						}
					}
					if workerNum != max1 && isDel {
						node.accountManager.WorkerNumberSet[workerNum] = node.accountManager.WorkerNumberSet[max1]
						delete(node.accountManager.WorkerNumberSet, max1)
					}

					if len(node.accountManager.WorkerSet) < node.config.WorkerNum && len(node.accountManager.WorkerCandidateSet) >= node.config.WorkerNum {
						for k1, v1 := range node.accountManager.WorkerCandidateSet {
							if _, ok := node.accountManager.WorkerSet[k1]; !ok {
								node.accountManager.WorkerSet[k1] = v1
								node.accountManager.WorkerNumberSet[max1] = k1
								break
							}
						}
					}

					var account MetaData.Account
					account.SetVoterSet(node.accountManager.VoterSet)
					account.SetWorkerSet(node.accountManager.WorkerSet)
					account.SetWorkerCandidateSet(node.accountManager.WorkerCandidateSet)
					account.WorkerCandidateList = node.accountManager.WorkerCandidateList
					account.SetVoterNumberSet(node.accountManager.VoterNumberSet)
					account.SetWorkerNumberSet(node.accountManager.WorkerNumberSet)
					node.mongo.InsertOrUpdateAccount(account)
					node.mongo.DeleteData("identity", "pubkey", k.Pubkey)
					node.network.RemoveNodeToNodeList(k.GetNodeId())
					fmt.Println("新节点列表更新成功")

					if node.config.MyPubkey == k.Pubkey {
						//os.Exit(0)
					}
				}
			}
		}
	}
}

func (node *Node) UpdateGenesisBlockVariables(bg *MetaData.BlockGroup) {
	if bg.Height == 0 { //genesis blockgroup
		node.dutyWorkerNumber = 0
		node.StartTime = bg.Timestamp
		if bg.Blocks[0].Height == 0 {
			transactionHeader, transactionInterface := MetaData.DecodeTransaction(bg.Blocks[0].Transactions[0])
			if transactionHeader.TXType == MetaData.Genesis {
				node.UpdateGenesisVaribles(transactionInterface)
			}
		}
		_ = node.state
		//node.state <- Normal
		node.state <- Sync
		time.Sleep(time.Second)
	} else {
		fmt.Println("更新变量错误")
	}
}

func (node *Node) UpdateGenesisVaribles(transactionInterface MetaData.TransactionInterface) {
	if genesisTransaction, ok := transactionInterface.(*MetaData.GenesisTransaction); ok {
		node.config.WorkerNum = genesisTransaction.WorkerNum
		node.config.VotedNum = genesisTransaction.VotedNum
		node.config.BlockGroupPerCycle = genesisTransaction.BlockGroupPerCycle
		node.config.Tcut = genesisTransaction.Tcut
		node.accountManager.WorkerSet = genesisTransaction.WorkerPubList
		node.accountManager.WorkerCandidateSet = genesisTransaction.WorkerCandidatePubList
		node.accountManager.VoterSet = genesisTransaction.VoterPubList
		var index uint32 = 0
		for _, key := range genesisTransaction.WorkerSet {
			node.accountManager.WorkerNumberSet[index] = key
			index = index + 1
		}
		index = 0
		for _, key1 := range genesisTransaction.VoterSet {
			node.accountManager.VoterNumberSet[index] = key1
			index = index + 1
		}
		for key2, _ := range genesisTransaction.WorkerCandidatePubList {
			node.accountManager.WorkerCandidateList = append(node.accountManager.WorkerCandidateList, key2)
		}
	}
}

func (node *Node) UpdateVariablesFromDisk(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		node.StartTime = bg.Timestamp
	} else {
		fmt.Println("更新变量错误")
	}

}
