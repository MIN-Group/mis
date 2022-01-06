package Message

import "fmt"

//go:generate msgp
type NodeStatusMsg struct {
	Type                string  `msg:"type"`
	NodeID              uint64  `msg:"number"`
	HostName            string  `msg:"hostname"`
	AreaName            string  `msg:"areaname"`
	CountryName         string  `msg:"countryname"`
	Longitude           float64 `msg:"longitude"`
	Latitude            float64 `msg:"latitude"`
	Pubkey              string  `msg:"publickey"`
	IP                  string  `msg:"ip"`
	Is_butler_candidate bool    `msg:"isbutlernext"`
	Is_butler           bool    `msg:"isbulter"`
	Is_commissioner     bool    `msg:"iscom"`
	Height              int     `msg:"height"`
	Agreement           int     `msg:"isagree"`
	Txs_num             int     `msg:"txs_num"`
}

func (gm NodeStatusMsg) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *NodeStatusMsg) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return err
}
