package Node

import (
	"fmt"
	"sync"

	"MIS-BC/utils"
)

type Syncer struct {
	NodeID        uint64
	Height        int
	StartTime     float64
	WaitTime      float64
	Valid         bool
	SyncingHeight int
}

type SynchronizeModule struct {
	Syncers                    chan Syncer //区块同步器队列
	SyncersLock                sync.Mutex  //线程锁
	FetchHeightStartTime       float64     //请求高度的开始时间
	FetchHeightWaitTime        float64     //请求高度的等待时间
	FetchHeightPeriodStartTime float64     //发送高度请求的时间间隔
	FetchHeightIntervalTIme    float64     //发送高度请求的时间间隔
	CurrentSyncer              Syncer      //正在进行同步的区块同步器
	SyncerWaitTime             float64     //同步器失效的等待时间
	phase                      int
}

func (sm *SynchronizeModule) Init() {
	sm.Syncers = make(chan Syncer, 2000)
	sm.FetchHeightStartTime = -1
	sm.FetchHeightWaitTime = 5
	sm.FetchHeightPeriodStartTime = -1
	sm.FetchHeightIntervalTIme = 1
	sm.SyncerWaitTime = 15
	sm.phase = 0
}

func (node *Node) syncBlockGroupSequence() {

	now := utils.GetCurrentTime()

	if node.syncTool.FetchHeightStartTime <= 0 {
		node.syncTool.FetchHeightStartTime = now
		node.syncTool.FetchHeightPeriodStartTime = now
	}
	switch node.syncTool.phase {
	case 0:
		{
			if now-node.syncTool.FetchHeightStartTime > node.syncTool.FetchHeightWaitTime {
				if len(node.syncTool.Syncers) == 0 {
					//fmt.Println("node.syncTool.phase=2")
					node.syncTool.phase = 2
				} else {
					//fmt.Println("node.syncTool.phase=1")
					node.syncTool.phase = 1
				}
				node.state <- Sync
				return
			}
			if now-node.syncTool.FetchHeightPeriodStartTime > node.syncTool.FetchHeightIntervalTIme {
				header, msg := node.msgManager.CreateRequestHeightMsg(0)
				node.SendMessage(header, &msg)
				fmt.Println("请求区块高度！")
				node.syncTool.FetchHeightPeriodStartTime = now
			}
		}
	case 1:
		{
			//fmt.Println("sync case 1:0")
			valid := false
			if now-node.syncTool.CurrentSyncer.StartTime < node.syncTool.CurrentSyncer.WaitTime && node.syncTool.CurrentSyncer.Height > node.mongo.GetHeight() {
				valid = true
			}
			for len(node.syncTool.Syncers) != 0 && !valid {
				node.syncTool.CurrentSyncer = <-node.syncTool.Syncers
				if node.syncTool.CurrentSyncer.Height > node.mongo.GetHeight() {
					valid = true
					node.syncTool.CurrentSyncer.StartTime = now
					node.syncTool.CurrentSyncer.WaitTime = node.syncTool.SyncerWaitTime
					node.syncTool.CurrentSyncer.SyncingHeight = -1
					//receiver := node.syncTool.CurrentSyncer.NodeID
					node.syncTool.Syncers <- node.syncTool.CurrentSyncer
				}
			}
			//fmt.Println("sync case 1:1")
			if valid {
				node.Commit()
				height := node.mongo.GetHeight()

				if height+1 == node.syncTool.CurrentSyncer.SyncingHeight {
					node.state <- Sync
					return
				}
				//else {
				//fmt.Println("node.syncTool.CurrentSyncer.SyncingHeight=",node.syncTool.CurrentSyncer.SyncingHeight)
				//fmt.Println("height+1=",height+1)
				//}
				if height < node.syncTool.CurrentSyncer.Height {

					receiver := node.syncTool.CurrentSyncer.NodeID
					fmt.Println("向节点", receiver, "请求高度为", height+1, "的区块组")
					header, msg := node.msgManager.CreateRequestBlockGroupMsg(receiver, height+1)
					node.SendMessage(header, &msg)
					node.syncTool.CurrentSyncer.SyncingHeight = height + 1
					//node.syncTool.CurrentSyncer.Height++
				}
			} else {
				node.syncTool.phase = 2
			}
			//fmt.Println("sync case 1:1")
		}
	case 2:
		{
			if node.mongo.GetHeight() < 0 {
				node.state <- Genesis
			} else {
				node.state <- Normal
			}
			return
		}
	}
	node.state <- Sync
}

func (node *Node) syncBlockGroupSequence2() {
	now := utils.GetCurrentTime()
	if node.syncTool.FetchHeightStartTime <= 0 {
		node.syncTool.FetchHeightStartTime = now
	}
	valid := false
	if now-node.syncTool.CurrentSyncer.StartTime < node.syncTool.CurrentSyncer.WaitTime && node.syncTool.CurrentSyncer.Height > node.mongo.GetHeight() {
		valid = true
	}
	if len(node.syncTool.Syncers) == 0 && !valid {
		if now-node.syncTool.FetchHeightStartTime > node.syncTool.FetchHeightWaitTime {
			height := node.mongo.GetHeight()
			fmt.Println("同步结束，高度为", height)
			if height == -1 {
				fmt.Println("同步区块：go to genesis")
				fmt.Println("进入Genesis阶段")
				node.state <- Genesis
			} else {
				for i := 0; i < height; i++ {
					var blockgroup = node.mongo.GetBlockFromDatabase(height)
					node.UpdateVariables(&blockgroup)
				}
				fmt.Println("进入Normal阶段")
				node.state <- Normal
			}
			return
		}
	}
	if len(node.syncTool.Syncers) == 0 {
		header, msg := node.msgManager.CreateRequestHeightMsg(0)
		node.SendMessage(header, &msg)
		fmt.Println("请求区块高度！")
	}
	for len(node.syncTool.Syncers) != 0 && !valid {
		node.syncTool.CurrentSyncer = <-node.syncTool.Syncers
		if node.syncTool.CurrentSyncer.Height > node.mongo.GetHeight() {
			valid = true
			node.syncTool.CurrentSyncer.StartTime = now
			node.syncTool.CurrentSyncer.WaitTime = node.syncTool.SyncerWaitTime
			receiver := node.syncTool.CurrentSyncer.NodeID
			header, msg := node.msgManager.CreateRequestBlockGroupMsg(receiver, node.mongo.GetHeight()+1)
			node.SendMessage(header, &msg)
		}
	}
	if valid {
		if node.mongo.GetHeight() < node.syncTool.CurrentSyncer.Height {
			receiver := node.syncTool.CurrentSyncer.NodeID
			header, msg := node.msgManager.CreateRequestBlockGroupMsg(receiver, node.mongo.GetHeight()+1)
			node.SendMessage(header, &msg)
			node.syncTool.CurrentSyncer.Height++
		}
	}
	node.state <- Sync
}
