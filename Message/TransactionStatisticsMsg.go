package Message

import "MIS-BC/common"

//go:generate msgp
type TransactionStatisticsMsg struct {
	Type      int    `msg:"type"`
	Pubkey    string `msg:"publickey"`
	Height    int    `msg:"height"`
	Agreement int    `msg:"isagree"`
	Txs_num   int    `msg:"txs_num"`
}

type TransactionStatisticsMsgs struct {
	Type string
	Msg  []TransactionStatisticsMsg
}

func (gm TransactionStatisticsMsgs) ToByteArray() ([]byte, error) {
	data, _ := gm.MarshalMsg(nil)
	return data, nil
}

func (gm *TransactionStatisticsMsgs) FromByteArray(data []byte) error {
	_, err := gm.UnmarshalMsg(data)
	if err != nil {
		common.Logger.Fatal(err)
	}
	return err
}
