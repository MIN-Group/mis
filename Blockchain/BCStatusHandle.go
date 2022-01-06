/**
 * @Author: xzw
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/6/5 下午3:00
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package Node

import (
	"MIS-BC/Message"
	"MIS-BC/MetaData"
	"MIS-BC/Network/network/encoding"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Response struct {
	StatusCode  int
	MessageType string
	Message     string
}

type TPSCensus struct {
	NodeNum    int
	Period_TPS int
	Total_tps  int
	Period     int
	Comm_num   int
	Butler_num int
	Total_time int
}

type BCStatus struct {
	Mutex    *sync.RWMutex
	Nodes    []MetaData.BCNode
	Agree    float64
	NoState  float64
	Disagree float64
}

func (node *Node) SendNodeStatusMsgToBCManagementServer() {
	fmt.Println("发送区块链状态信息消息")
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.NodeStatus
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0

	if node.mongo.Height%20 == 0 || node.mongo.Height < 100 && node.mongo.Height%5 == 0 {
		var msg Message.NodeStatusMsg //消息体
		msg.Type = "HandleNodeStatusMsgFromBC"
		msg.Height = node.mongo.Height
		msg.Pubkey = node.config.MyPubkey
		msg.HostName = node.config.HostName
		msg.IP = node.network.MyNodeInfo.IP
		msg.HostName = node.network.MyNodeInfo.HostName
		msg.AreaName = node.network.MyNodeInfo.AreaName
		msg.CountryName = node.network.MyNodeInfo.CountryName
		msg.Longitude = node.network.MyNodeInfo.Longitude
		msg.Latitude = node.network.MyNodeInfo.Latitude

		_, msg.Is_butler = node.accountManager.WorkerSet[msg.Pubkey]
		_, msg.Is_butler_candidate = node.accountManager.WorkerCandidateSet[msg.Pubkey]
		_, msg.Is_commissioner = node.accountManager.VoterSet[msg.Pubkey]

		node.SendMessage(mh, &msg)
	}
}

func (node *Node) SendTransactionMsgToManagementServer(bg MetaData.BlockGroup) {
	fmt.Println("发送区块链交易统计信息消息")
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.TransactionStatistics
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var msg Message.TransactionStatisticsMsgs
	msg.Type = "HandleTransactionMsgFromBC"
	var num_of_trans = 0
	for _, eachBlock := range bg.Blocks {
		num_of_trans += len(eachBlock.Transactions)
	}

	for i, eachticket := range bg.VoteTickets {
		if i >= len(bg.VoteResult) {
			return
		}
		var one_msg Message.TransactionStatisticsMsg
		one_msg.Type = 0
		one_msg.Pubkey = eachticket.Voter
		one_msg.Height = bg.Height
		one_msg.Agreement = bg.VoteResult[i]
		one_msg.Txs_num = num_of_trans
		msg.Msg = append(msg.Msg, one_msg)
	}

	node.SendMessage(mh, &msg)
}

func (node *Node) HandleNodeStatusMsgFromBC(msg Message.NodeStatusMsg, header Message.MessageHeader) {
	node.BCStatus.Mutex.Lock()
	for i, item := range node.BCStatus.Nodes {
		if item.Pubkey == msg.Pubkey {
			node.BCStatus.Nodes[i].HostName = msg.HostName
			node.BCStatus.Nodes[i].AreaName = msg.AreaName
			node.BCStatus.Nodes[i].CountryName = msg.CountryName
			node.BCStatus.Nodes[i].Longitude = msg.Longitude
			node.BCStatus.Nodes[i].Latitude = msg.Latitude

			if msg.Height > node.BCStatus.Nodes[i].Height {
				node.BCStatus.Nodes[i].Height = msg.Height
			}
			node.BCStatus.Nodes[i].IP = msg.IP
			node.BCStatus.Nodes[i].Is_butler_candidate = msg.Is_butler_candidate
			node.BCStatus.Nodes[i].Is_butler = msg.Is_butler
			node.BCStatus.Nodes[i].Is_commissioner = msg.Is_commissioner
			node.BCStatus.Mutex.Unlock()
			return
		}
	}
	var new_item MetaData.BCNode
	new_item.Number = len(node.BCStatus.Nodes)
	new_item.Pubkey = msg.Pubkey
	new_item.HostName = msg.HostName
	new_item.IP = msg.IP
	new_item.Is_butler_candidate = msg.Is_butler_candidate
	new_item.Is_butler = msg.Is_butler
	new_item.Is_commissioner = msg.Is_commissioner
	new_item.CountryName = msg.CountryName
	new_item.AreaName = msg.AreaName
	new_item.HostName = msg.HostName
	new_item.Longitude = msg.Longitude
	new_item.Latitude = msg.Latitude
	node.BCStatus.Nodes = append(node.BCStatus.Nodes, new_item)

	node.BCStatus.Mutex.Unlock()
}

func (node *Node) HandleTransactionMsgFromBC(gm Message.TransactionStatisticsMsgs, header Message.MessageHeader) {
	node.BCStatus.Mutex.Lock()
	node.BCStatus.Agree = 0
	node.BCStatus.NoState = 0
	node.BCStatus.Disagree = 0
	for _, item := range gm.Msg {
		if item.Type != 0 {
			continue
		}
		for i, each_node := range node.BCStatus.Nodes {
			if item.Pubkey == each_node.Pubkey {
				if item.Height > node.BCStatus.Nodes[i].Height {
					node.BCStatus.Nodes[i].Height = item.Height
				}
				node.BCStatus.Nodes[i].Agreement = item.Agreement
				node.BCStatus.Nodes[i].Txs_num = item.Txs_num
				switch item.Agreement {
				case 1:
					node.BCStatus.Agree += 1
				case -1:
					node.BCStatus.Disagree += 1
				}
				break
			}
		}
	}
	if len(node.BCStatus.Nodes) > 0 {
		node.BCStatus.Agree /= float64(len(node.BCStatus.Nodes))
		node.BCStatus.Disagree /= float64(len(node.BCStatus.Nodes))
		node.BCStatus.NoState = 1 - node.BCStatus.Agree - node.BCStatus.Disagree
	}

	if len(node.BCStatus.Nodes) == node.config.ServerNum {
		var info MetaData.BCStatus
		info.Agree = node.BCStatus.Agree
		info.NoState = node.BCStatus.NoState
		info.Disagree = node.BCStatus.Disagree
		info.Nodeinfo = node.BCStatus.Nodes
		info.Timestamp = time.Now().Format("2006-01-02 15:04:05")

		node.mongo.SaveBCStatusToDatabase(info)
	}
	node.BCStatus.Mutex.Unlock()
}

func (node *Node) AgreeAddNewNode(res map[string]interface{}, conn net.Conn) {
	if res["Pubkey"] == nil {
		return
	}
	pubkey := res["Pubkey"].(string)

	node.NewNodeList.Set(pubkey, 1, 3600*time.Second)
	var response Response
	response.StatusCode = 200
	response.MessageType = "string"
	response.Message = "注册成功"

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	new_data, err := encoding.Encode(data)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(new_data)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (node *Node) RemoveNodeApply(res map[string]interface{}, conn net.Conn) {
	if res["Pubkey"] == nil {
		return
	}
	var response Response
	pubkey := res["Pubkey"].(string)
	if nodeid, ok := node.accountManager.VoterSet[pubkey]; ok {
		node.RequestRemoveNode(pubkey, nodeid)
		response.StatusCode = 200
		response.MessageType = "string"
		response.Message = "删除申请成功"
	} else {
		response.StatusCode = 400
		response.MessageType = "string"
		response.Message = "该节点不是投票节点"
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	new_data, err := encoding.Encode(data)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(new_data)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (node *Node) QuitMyself(conn net.Conn) {
	node.SendRequestQuitNode()
	var response Response
	response.StatusCode = 200
	response.MessageType = "string"
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	new_data, err := encoding.Encode(data)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(new_data)
	if err != nil {
		fmt.Println(err)
	}
}

func (node *Node) GetNodeListPrint(conn net.Conn) {

	fmt.Println("WorkerNumberSet", node.accountManager.WorkerNumberSet)
	fmt.Println("WorkerSet", node.accountManager.WorkerSet)
	fmt.Println("VoterSet", node.accountManager.VoterSet)
	fmt.Println("VoterNumberSet", node.accountManager.VoterNumberSet)
	fmt.Println("NetworkNodeList", node.network.NodeList)

	var response Response
	response.StatusCode = 200
	response.MessageType = "string"
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	new_data, err := encoding.Encode(data)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(new_data)
	if err != nil {
		fmt.Println(err)
	}
	return
}
