package Node

import (
	"MIS-BC/MetaData"
)

func (node *Node) ValidateBlockHeader(b *MetaData.Block) bool {
	_, existed := node.accountManager.WorkerSet[b.Generator]
	if !existed {
		return false
	}
	if node.accountManager.WorkerNumberSet[b.BlockNum] != b.Generator {
		return false
	}
	return true
}

func (node *Node) ValidateTransactions(txs *([][]byte)) bool {
	for i := 0; i < len(*txs); i++ {
		header, transactionInterface := MetaData.DecodeTransaction((*txs)[i])
		switch header.TXType {
		case MetaData.Zero:
			if transaction, ok := transactionInterface.(*MetaData.ZeroTransaction); ok {
				if !node.ValidateZeroTransaction(*transaction) {
					return false
				}
			}
		case MetaData.Genesis:
			if transaction, ok := transactionInterface.(*MetaData.GenesisTransaction); ok {
				if !node.ValidateGenesisTransaction(*transaction) {
					return false
				}
			}
		}
	}
	return true
}

func (node *Node) ValidateZeroTransaction(tx MetaData.ZeroTransaction) bool {
	return true
}

func (node *Node) ValidateGenesisTransaction(tx MetaData.GenesisTransaction) bool {
	return true
}
