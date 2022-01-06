package Network

import (
	"MIS-BC/Network/network/encoding"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func SendPacket(message []byte, ip string, port int) {
	port_s := strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", ip+":"+port_s, 3*time.Second)
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	data, err := encoding.Encode(message)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}

}

func SendPacketAndGetAns(message []byte, ip string, port int) []byte {
	port_s := strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", ip+":"+port_s, 3*time.Second)
	if err != nil {
		fmt.Println("dial failed, err", err)
		return nil
	}
	defer conn.Close()
	data, err := encoding.Encode(message)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return nil
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return nil
	}

	msg, err := encoding.DecodeTcp(conn)
	if err == io.EOF {
		fmt.Println("IO errror, err:", err)
		return nil
	}
	if err != nil {
		fmt.Println("decode msg failed, err:", err)
		return nil
	}
	return msg
}

func (network *Network) SendToAll(message []byte) {
	network.Mutex.RLock()
	tmpList := network.NodeList
	network.Mutex.RUnlock()
	for _, x := range tmpList {
		SendPacket(message, x.IP, x.PORT)
	}
}

func (network *Network) SendToNeighbor(message []byte) {
	network.Mutex.RLock()
	tmpList := network.NodeList
	temp := network.MyNodeInfo.ID
	network.Mutex.RUnlock()

	for _, x := range tmpList {
		if x.ID == temp {
			continue
		}
		SendPacket(message, x.IP, x.PORT)
	}
}

func (network *Network) SendToOne(message []byte, receiver NodeID) {
	network.Mutex.RLock()
	ip := network.NodeList[receiver].IP
	port := network.NodeList[receiver].PORT
	network.Mutex.RUnlock()
	SendPacket(message, ip, port)
}

func (network *Network) SendMessage(message []byte, receiver NodeID) {
	if receiver == 0 {
		network.SendToAll(message)
	} else if receiver == 1 {
		network.SendToNeighbor(message)
	} else {
		network.SendToOne(message, receiver)
	}
}
