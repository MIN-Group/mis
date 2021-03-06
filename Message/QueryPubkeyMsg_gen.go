package Message

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *QueryPubkeyMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "type":
			z.Type, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Type")
				return
			}
		case "info":
			z.Information, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Information")
				return
			}
		case "nodeid":
			z.NodeID, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "NodeID")
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
func (z QueryPubkeyMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "type"
	err = en.Append(0x83, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Type)
	if err != nil {
		err = msgp.WrapError(err, "Type")
		return
	}
	// write "info"
	err = en.Append(0xa4, 0x69, 0x6e, 0x66, 0x6f)
	if err != nil {
		return
	}
	err = en.WriteString(z.Information)
	if err != nil {
		err = msgp.WrapError(err, "Information")
		return
	}
	// write "nodeid"
	err = en.Append(0xa6, 0x6e, 0x6f, 0x64, 0x65, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.NodeID)
	if err != nil {
		err = msgp.WrapError(err, "NodeID")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z QueryPubkeyMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "type"
	o = append(o, 0x83, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Type)
	// string "info"
	o = append(o, 0xa4, 0x69, 0x6e, 0x66, 0x6f)
	o = msgp.AppendString(o, z.Information)
	// string "nodeid"
	o = append(o, 0xa6, 0x6e, 0x6f, 0x64, 0x65, 0x69, 0x64)
	o = msgp.AppendUint64(o, z.NodeID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *QueryPubkeyMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "type":
			z.Type, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Type")
				return
			}
		case "info":
			z.Information, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Information")
				return
			}
		case "nodeid":
			z.NodeID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NodeID")
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
func (z QueryPubkeyMsg) Msgsize() (s int) {
	s = 1 + 5 + msgp.IntSize + 5 + msgp.StringPrefixSize + len(z.Information) + 7 + msgp.Uint64Size
	return
}
