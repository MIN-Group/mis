/**
 * @Author: xzw
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/6/2 晚上9:00
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package Network

import (
	"MIS-BC/Network/network/encoding"
	"fmt"
	"io"
	"log"

	"net"
	"strconv"
)

// for blockchain over ip
func (network *Network) HandleConnection(conn net.Conn) {
	defer conn.Close()
	// 从连接中解码消息 前四个字节指定大小 防止粘包
	msg, err := encoding.DecodeTcp(conn)
	if err == io.EOF {
		return
	}
	if err != nil {
		fmt.Println("decode msg failed, err:", err)
		return
	}
	network.CBforBC(msg)
}

func (network *Network) Start() {
	// 监听本节点端口
	port := strconv.Itoa(network.MyNodeInfo.PORT)
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+port)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	// 死循环
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// 开启协程单独处理请求
		go network.HandleConnection(conn)
	}
}