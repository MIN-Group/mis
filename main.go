package main

import (
	"MIS-BC/Blockchain"
	"MIS-BC/common"
	_ "encoding/binary"
	"flag"
	"fmt"
	_ "github.com/tinylib/msgp/msgp"
	"strconv"
)

// @Title Run
// @Description 区块链启动函数
// @Param 无
// @Return 无
func Run() {
	conf := common.ParseConfig(file) //读取配置
	common.Logger.Init(&conf)
	var nodes []Node.Node //节点数组

	//创建节点并设置参数
	for i := 0; i < conf.SingleServerNodeNum; i++ {
		var parcel = conf
		var node Node.Node
		parcel.DropDatabase = dropDatabase
		parcel.SendMsgToBCMgmt = sendMsgToMgmt
		parcel.MyPubkey = conf.PubkeyList[i]
		parcel.MyPrikey = conf.PrikeyList[i]
		parcel.MyAddress.Port += i
		if conf.SingleServerNodeNum > 1 {
			parcel.HostName = parcel.HostName + strconv.Itoa(i+1)
		}

		node.Init()
		node.SetConfig(parcel)
		nodes = append(nodes, node)
	}
	if len(nodes) == 1 {
		// 如果只有一个节点 不需要进行共识
		fmt.Println("node 0 start")
		// 启动该节点
		nodes[0].Start()
	} else {
		//初始化并启动节点
		for i := 0; i < len(nodes); i++ {
			go nodes[i].Start()
			fmt.Println("node", i, "start")
		}
		// 阻塞
		select {}
	}
}

var (
	file          string //配置文件路径
	dropDatabase  bool   //是否清空数据库
	sendMsgToMgmt bool   //是否向前端发送状态信息
	//addr		  string //rpc地址
)

func init() {
	flag.StringVar(&file, "f", "default", "config file")              //配置文件路径
	flag.BoolVar(&dropDatabase, "d", false, "delete database")        //是否清空数据库
	flag.BoolVar(&sendMsgToMgmt, "s", true, "send msg to management") //是否向前端发送状态信息
}

func main() {
	flag.Parse()
	fmt.Println("欢迎使用PPoV区块链!")
	Run()
}
