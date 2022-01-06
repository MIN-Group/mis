package main

import (
	"MIS-BC/security/keymanager"
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:5020", "server address")
	d, _ = client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	opt = client.DefaultOption
)

type Args struct {
	IdentityIdentifier string
	Pubkey             string
	Passwd             string
}

type CommonResponse struct {
	Code    int
	Message string
	Data    interface{}
}


func inquire(args Args) {
	xclient := client.NewXClient("Registry", client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()
	reply := &CommonResponse{}
	time1 := time.Now().UnixNano()
	err := xclient.Call(context.Background(), "GetOneIdentityInfByIdentityIdentifierforTest", args, reply)
	time2 := time.Now().UnixNano()
	fmt.Println("time2-time1",time2-time1)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(reply)
}

func registry(args Args){
	xclient := client.NewXClient("Registry", client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()
	// TODO 调用keygen生成公钥
	var keyManager keymanager.KeyManager
	keyManager.Init()
	keyManager.GenKeyPair()
	args.Pubkey = keyManager.GetPubkey()

	reply := &CommonResponse{}
	err := xclient.Call(context.Background(), "IdentityRegistryforTest", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Println(reply)
}

func destroy(args Args){
	xclient := client.NewXClient("Registry", client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()
	reply := &CommonResponse{}
	err := xclient.Call(context.Background(), "DestroyByIdentityIdentifierforTest", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Println(reply)
}

func resetPasswd(args Args){
	xclient := client.NewXClient("Registry", client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()
	reply := &CommonResponse{}
	err := xclient.Call(context.Background(), "ResetPasswordforTest", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(reply)
}

func main() {
	flag.Parse()
	opt.SerializeType = protocol.JSON
	var choice int
	fmt.Print("请选择需要的函数: \n 1.注册标识\n 2.查询标识\n 3.修改密码\n 4.删除标识\n")
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		var args Args
		fmt.Println("请输入身份标识名:")
		fmt.Scanln(&args.IdentityIdentifier)
		fmt.Println("请输入密码:")
		fmt.Scanln(&args.Passwd)
		registry(args)

	case 2:
		var args Args
		fmt.Println("请输入身份标识名:")
		fmt.Scanln(&args.IdentityIdentifier)
		inquire(args)

	case 3:
		var args Args
		fmt.Println("请输入身份标识名:")
		fmt.Scanln(&args.IdentityIdentifier)
		fmt.Println("请输入新密码:")
		fmt.Scanln(&args.Passwd)
		resetPasswd(args)

	case 4:
		var args Args
		fmt.Println("请输入身份标识名:")
		fmt.Scanln(&args.IdentityIdentifier)
		destroy(args)
	}
}
