package MetaData

import (
	"MIS-BC/Network"
	"strconv"
)

type NodeInfoStr struct {
	IP     string
	PORT   int
	ID     string
	Prefix string
}

//go:generate msgp
type NodeList struct {
	NodeList map[string]NodeInfoStr `msg:"NodeList"`
}

/*
func (min NodeList) ToByteArray() []byte {
	data, _ := min.MarshalMsg(nil)
	return data
}

func (min *NodeList) FromByteArray(data []byte) {
	_, err := min.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
*/
func (n *NodeList) SetNodeList(NodeList map[Network.NodeID]Network.NodeInfo) {
	nl := make(map[string]NodeInfoStr)
	for k, v := range NodeList {
		var ntmp NodeInfoStr
		ntmp.IP = v.IP
		ntmp.PORT = v.PORT
		ntmp.Prefix = v.Prefix
		ntmp.ID = strconv.FormatUint(k, 10)
		nl[strconv.FormatUint(k, 10)] = ntmp
	}
	n.NodeList = nl
}

func (n *NodeList) GetNodeList() map[Network.NodeID]Network.NodeInfo {
	nl := make(map[Network.NodeID]Network.NodeInfo)
	for k, v := range n.NodeList {
		tmp, _ := strconv.ParseUint(k, 0, 64)
		var ntmp Network.NodeInfo
		ntmp.ID, _ = strconv.ParseUint(k, 0, 64)
		ntmp.PORT = v.PORT
		ntmp.IP = v.IP
		ntmp.Prefix = v.Prefix
		nl[Network.NodeID(tmp)] = ntmp
	}
	return nl
}
