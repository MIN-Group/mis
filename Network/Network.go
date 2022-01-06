package Network

import (
	"MIS-BC/common"
	"MIS-BC/security"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"sync"
)

type NodeID = uint64

type NodeInfo struct {
	IP          string
	PORT        int
	Prefix      string
	ID          NodeID
	HostName    string
	AreaName    string
	CountryName string
	Longitude   float64 //经度
	Latitude    float64 //纬度
}

//implementation of set which saves space
type Void struct{}

type Network struct {
	MyNodeInfo NodeInfo
	NodeList   map[NodeID]NodeInfo

	CBforBC       func([]byte)

	Blacklist map[NodeID]Void
	Mutex     *sync.RWMutex

	CertificateInquires string
	Keychain            *security.KeyChain
}

func IPToValue(strIP string) uint32 {
	var a [4]uint32
	temp := strings.Split(strIP, ".")
	for i, x := range temp {
		t, err := strconv.Atoi(x)
		if err != nil {
			fmt.Println(err)
		}
		a[i] = uint32(t)
	}
	var ret uint32 = (a[0] << 24) + (a[1] << 16) + (a[2] << 8) + a[3]
	return ret
}

func GetNodeId(ip string, port int, prefix string) NodeID {
	if prefix != "" {
		return NodeID(crc32.ChecksumIEEE([]byte(prefix)))
	} else {
		var id uint32 = IPToValue(ip)
		var nid NodeID = NodeID(id) << 32
		nid += NodeID(port)
		return nid
	}
}

func (network *Network) AddNodeToNodeList(NodeID uint64, IPAddr string, Port int) {
	var nodelist NodeInfo
	nodelist.IP = IPAddr
	nodelist.PORT = Port
	nodelist.ID = NodeID
	network.Mutex.Lock()
	network.NodeList[NodeID] = nodelist
	network.Mutex.Unlock()
}

func (network *Network) RemoveNodeToNodeList(NodeID uint64) {
	network.Mutex.Lock()
	delete(network.NodeList, NodeID)
	network.Mutex.Unlock()
}

func (network *Network) SetConfig(config common.Config) {
	if config.MyAddress.PubIP != "" {
		network.MyNodeInfo.IP = config.MyAddress.PubIP
	} else {
		network.MyNodeInfo.IP = config.MyAddress.IP
	}

	// index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(config.MyPubkey))))

	network.MyNodeInfo.PORT = config.MyAddress.Port
	network.MyNodeInfo.Prefix = config.MyAddress.Prefix
	network.MyNodeInfo.ID = GetNodeId(network.MyNodeInfo.IP, network.MyNodeInfo.PORT, network.MyNodeInfo.Prefix)
	network.MyNodeInfo.HostName = config.HostName
	network.MyNodeInfo.AreaName = config.AreaName
	network.MyNodeInfo.CountryName = config.CountryName
	network.MyNodeInfo.Longitude = config.Longitude
	network.MyNodeInfo.Latitude = config.Latitude

	network.NodeList = make(map[NodeID]NodeInfo)

	network.Keychain.InitialKeyChainByPath(config.SqlitePath + config.HostName + "/")

	for _, x := range config.WorkerList {
		temp := GetNodeId(x.IP, x.Port, x.Prefix)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.Prefix = x.Prefix
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
	for _, x := range config.WorkerCandidateList {
		temp := GetNodeId(x.IP, x.Port, x.Prefix)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.Prefix = x.Prefix
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
	for _, x := range config.VoterList {
		temp := GetNodeId(x.IP, x.Port, x.Prefix)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.Prefix = x.Prefix
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
}

func (network *Network) SetCB(cbforbc func([]byte)) {
	network.CBforBC = cbforbc
}
