package MetaData

const (
	Zero = iota
	Genesis
	IdTransformation
	IdentityAction
	ElectNewWorker
	UserLogOperation
)

type TransactionInterface interface {
	ToByteArray() []byte
	FromByteArray(data []byte)
}

//go:generate msgp
type TransactionHeader struct {
	TXType int    `msg:"tx"`
	Data   []byte `msg:"data"`
}

func EncodeTransaction(header TransactionHeader, transactionInterface TransactionInterface) (data []byte) {
	data = transactionInterface.ToByteArray()
	header.Data = data
	data, _ = header.MarshalMsg(nil)
	return data
}

func DecodeTransaction(data []byte) (header TransactionHeader, transactionInterface TransactionInterface) {
	data, _ = header.UnmarshalMsg(data)
	data = header.Data
	switch header.TXType {
	case Zero:
		var zt ZeroTransaction
		zt.FromByteArray(data)
		transactionInterface = &zt
	case Genesis:
		var gt GenesisTransaction
		gt.FromByteArray(data)
		transactionInterface = &gt
	case IdTransformation:
		var idt IdentityTransformation
		idt.FromByteArray(data)
		transactionInterface = &idt
	case IdentityAction:
		var id Identity
		id.FromByteArray(data)
		transactionInterface = &id
	case ElectNewWorker:
		var emwt ElectNewWorkerTeam
		emwt.FromByteArray(data)
		transactionInterface = &emwt
	}
	return
}
