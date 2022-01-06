package Node

import (
	"MIS-BC/MetaData"
	"math/rand"
	"time"
)

func (node *Node) TransactionGenerator(n int) {
	for {
		var trans MetaData.ZeroTransaction

		trans.Content = make([]byte, n)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < len(trans.Content); i++ {
			trans.Content[i] = byte(rand.Intn(256))
		}
		var transactionHeader MetaData.TransactionHeader
		transactionHeader.TXType = MetaData.Zero
		output := node.txPool.PushbackTransaction(transactionHeader, &trans)
		if output == -1 {
			break
		}
	}
}
