package MetaData

//go:generate msgp
type ByteArrayHelper struct {
	Data []byte `msg:"data"`
}

type DoubleByteArrayHelper struct {
	Data [][]byte `msg:"data"`
}
