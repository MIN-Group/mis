/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/17 下午7:57
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

import (
	km "MIS-BC/security/keymanager"
	"fmt"
	"testing"
)

func TestTCPServer(t *testing.T) {
	fmt.Println("tcp server start...")
	net := TCPNet{}
	key := km.KeyManager{}
	key.Init()
	key.SetPriKey("c907424efdce6903ff380a351d7f2cce3ff06d58512792873abd090bf1c9a84b")
	key.SetPubkey("04e067d23d3fd3ba4ead731b78346cde1084a837760f234621e36a93632b6f6418d4d6b735cd0ee19a8707ceca12ee4ff8106196ad8122ab2c89a7e11ba00d35b7")
	net.Listens("12345", "54321", &key)

	for {
		conn, err := net.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		fmt.Println("Received client message...")
		buf, err := conn.Read()
		if err != nil {
			fmt.Println("Error Reading", err.Error())
			return // 终止程序
		}
		fmt.Printf("Received data: %v", string(buf))
		fmt.Println()

		buf = []byte(string("hello client"))
		err = conn.Write(buf)
		if err != nil {
			fmt.Println("Error writing", err.Error())
			return // 终止程序
		}
		fmt.Println("The client send ", string(buf))
		conn.Close()
		fmt.Println("Finish the connection...")
		fmt.Println()
	}

}

func TestTCPClient(t *testing.T) {
	net := TCPNet{}
	key := km.KeyManager{}
	key.GenKeyPair()
	conn, err := net.Dials("localhost:12345", "localhost:54321", &key, []byte("04e067d23d3fd3ba4ead731b78346cde1084a837760f234621e36a93632b6f6418d4d6b735cd0ee19a8707ceca12ee4ff8106196ad8122ab2c89a7e11ba00d35b7"))
	//conn, err := net.Dial("localhost:10010")
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}

	b := []byte("hello server")
	err = conn.Write(b)
	if err != nil {
		fmt.Println("Error send message", err.Error())
		return // 终止程序
	}
	fmt.Println("The client send ", string(b))

	fmt.Println("Received server message...")
	buf, err := conn.Read()
	if err != nil {
		fmt.Println("Error Reading", err.Error())
		return // 终止程序
	}
	fmt.Printf("Received data: %v", string(buf))
	conn.Close()
}
