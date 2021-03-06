package Message

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ElectNewWorkerMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "newworker":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "NewWorker")
				return
			}
			if cap(z.NewWorker) >= int(zb0002) {
				z.NewWorker = (z.NewWorker)[:zb0002]
			} else {
				z.NewWorker = make([]string, zb0002)
			}
			for za0001 := range z.NewWorker {
				z.NewWorker[za0001], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "NewWorker", za0001)
					return
				}
			}
		case "mypubkey":
			z.MyPubkey, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "MyPubkey")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ElectNewWorkerMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "newworker"
	err = en.Append(0x82, 0xa9, 0x6e, 0x65, 0x77, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.NewWorker)))
	if err != nil {
		err = msgp.WrapError(err, "NewWorker")
		return
	}
	for za0001 := range z.NewWorker {
		err = en.WriteString(z.NewWorker[za0001])
		if err != nil {
			err = msgp.WrapError(err, "NewWorker", za0001)
			return
		}
	}
	// write "mypubkey"
	err = en.Append(0xa8, 0x6d, 0x79, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.MyPubkey)
	if err != nil {
		err = msgp.WrapError(err, "MyPubkey")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ElectNewWorkerMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "newworker"
	o = append(o, 0x82, 0xa9, 0x6e, 0x65, 0x77, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72)
	o = msgp.AppendArrayHeader(o, uint32(len(z.NewWorker)))
	for za0001 := range z.NewWorker {
		o = msgp.AppendString(o, z.NewWorker[za0001])
	}
	// string "mypubkey"
	o = append(o, 0xa8, 0x6d, 0x79, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	o = msgp.AppendString(o, z.MyPubkey)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ElectNewWorkerMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "newworker":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NewWorker")
				return
			}
			if cap(z.NewWorker) >= int(zb0002) {
				z.NewWorker = (z.NewWorker)[:zb0002]
			} else {
				z.NewWorker = make([]string, zb0002)
			}
			for za0001 := range z.NewWorker {
				z.NewWorker[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "NewWorker", za0001)
					return
				}
			}
		case "mypubkey":
			z.MyPubkey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MyPubkey")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ElectNewWorkerMsg) Msgsize() (s int) {
	s = 1 + 10 + msgp.ArrayHeaderSize
	for za0001 := range z.NewWorker {
		s += msgp.StringPrefixSize + len(z.NewWorker[za0001])
	}
	s += 9 + msgp.StringPrefixSize + len(z.MyPubkey)
	return
}
