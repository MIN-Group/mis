package Message

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *PublishIamOkMsg) DecodeMsg(dc *msgp.Reader) (err error) {
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
			z.Type, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Type")
				return
			}
		case "pubkey":
			z.Pubkey, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Pubkey")
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
func (z PublishIamOkMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "type"
	err = en.Append(0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Type)
	if err != nil {
		err = msgp.WrapError(err, "Type")
		return
	}
	// write "pubkey"
	err = en.Append(0xa6, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.Pubkey)
	if err != nil {
		err = msgp.WrapError(err, "Pubkey")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PublishIamOkMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "type"
	o = append(o, 0x82, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	// string "pubkey"
	o = append(o, 0xa6, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	o = msgp.AppendString(o, z.Pubkey)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PublishIamOkMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Type")
				return
			}
		case "pubkey":
			z.Pubkey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Pubkey")
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
func (z PublishIamOkMsg) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Type) + 7 + msgp.StringPrefixSize + len(z.Pubkey)
	return
}