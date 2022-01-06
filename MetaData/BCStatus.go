package MetaData

import (
	"fmt"
)

type BCNode struct {
	Number              int     `json:"number"`
	Pubkey              string  `json:"publickey"`
	IP                  string  `json:"ip"`
	Is_butler_candidate bool    `json:"isbutlernext"`
	Is_butler           bool    `json:"isbulter"`
	Is_commissioner     bool    `json:"iscom"`
	Height              int     `json:"height"`
	Agreement           int     `json:"isagree"`
	Txs_num             int     `json:"txs_num"`
	HostName            string  `json:"hostname"`
	AreaName            string  `json:"areaname"`
	CountryName         string  `json:"countryname"`
	Longitude           float64 `json:"longitude"`
	Latitude            float64 `json:"latitude"`
}

//go:generate msgp
type BCStatus struct {
	Agree     float64  `msg:"agree"`
	NoState   float64  `msg:"no_state"`
	Disagree  float64  `msg:"disagree"`
	Nodeinfo  []BCNode `msg:"nodeinfo"`
	Timestamp string   `msg:"timestamp"` //注册时间
}

func (bs BCStatus) ToByteArray() []byte {
	data, _ := bs.MarshalMsg(nil)
	return data
}

func (bs *BCStatus) FromByteArray(data []byte) {
	_, err := bs.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
