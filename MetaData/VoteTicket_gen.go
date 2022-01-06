package MetaData

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *VoteTicket) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "VoteResult":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "VoteResult")
				return
			}
			if cap(z.VoteResult) >= int(zb0002) {
				z.VoteResult = (z.VoteResult)[:zb0002]
			} else {
				z.VoteResult = make([]int, zb0002)
			}
			for za0001 := range z.VoteResult {
				z.VoteResult[za0001], err = dc.ReadInt()
				if err != nil {
					err = msgp.WrapError(err, "VoteResult", za0001)
					return
				}
			}
		case "hashes":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "BlockHashes")
				return
			}
			if cap(z.BlockHashes) >= int(zb0003) {
				z.BlockHashes = (z.BlockHashes)[:zb0003]
			} else {
				z.BlockHashes = make([]string, zb0003)
			}
			for za0002 := range z.BlockHashes {
				z.BlockHashes[za0002], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "BlockHashes", za0002)
					return
				}
			}
		case "Voter":
			z.Voter, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Voter")
				return
			}
		case "Timestamp":
			z.Timestamp, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "Sig":
			z.Sig, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Sig")
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
func (z *VoteTicket) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "VoteResult"
	err = en.Append(0x85, 0xaa, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.VoteResult)))
	if err != nil {
		err = msgp.WrapError(err, "VoteResult")
		return
	}
	for za0001 := range z.VoteResult {
		err = en.WriteInt(z.VoteResult[za0001])
		if err != nil {
			err = msgp.WrapError(err, "VoteResult", za0001)
			return
		}
	}
	// write "hashes"
	err = en.Append(0xa6, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.BlockHashes)))
	if err != nil {
		err = msgp.WrapError(err, "BlockHashes")
		return
	}
	for za0002 := range z.BlockHashes {
		err = en.WriteString(z.BlockHashes[za0002])
		if err != nil {
			err = msgp.WrapError(err, "BlockHashes", za0002)
			return
		}
	}
	// write "Voter"
	err = en.Append(0xa5, 0x56, 0x6f, 0x74, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Voter)
	if err != nil {
		err = msgp.WrapError(err, "Voter")
		return
	}
	// write "Timestamp"
	err = en.Append(0xa9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Timestamp)
	if err != nil {
		err = msgp.WrapError(err, "Timestamp")
		return
	}
	// write "Sig"
	err = en.Append(0xa3, 0x53, 0x69, 0x67)
	if err != nil {
		return
	}
	err = en.WriteString(z.Sig)
	if err != nil {
		err = msgp.WrapError(err, "Sig")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *VoteTicket) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "VoteResult"
	o = append(o, 0x85, 0xaa, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.VoteResult)))
	for za0001 := range z.VoteResult {
		o = msgp.AppendInt(o, z.VoteResult[za0001])
	}
	// string "hashes"
	o = append(o, 0xa6, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.BlockHashes)))
	for za0002 := range z.BlockHashes {
		o = msgp.AppendString(o, z.BlockHashes[za0002])
	}
	// string "Voter"
	o = append(o, 0xa5, 0x56, 0x6f, 0x74, 0x65, 0x72)
	o = msgp.AppendString(o, z.Voter)
	// string "Timestamp"
	o = append(o, 0xa9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o = msgp.AppendFloat64(o, z.Timestamp)
	// string "Sig"
	o = append(o, 0xa3, 0x53, 0x69, 0x67)
	o = msgp.AppendString(o, z.Sig)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VoteTicket) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "VoteResult":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "VoteResult")
				return
			}
			if cap(z.VoteResult) >= int(zb0002) {
				z.VoteResult = (z.VoteResult)[:zb0002]
			} else {
				z.VoteResult = make([]int, zb0002)
			}
			for za0001 := range z.VoteResult {
				z.VoteResult[za0001], bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "VoteResult", za0001)
					return
				}
			}
		case "hashes":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlockHashes")
				return
			}
			if cap(z.BlockHashes) >= int(zb0003) {
				z.BlockHashes = (z.BlockHashes)[:zb0003]
			} else {
				z.BlockHashes = make([]string, zb0003)
			}
			for za0002 := range z.BlockHashes {
				z.BlockHashes[za0002], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "BlockHashes", za0002)
					return
				}
			}
		case "Voter":
			z.Voter, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Voter")
				return
			}
		case "Timestamp":
			z.Timestamp, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "Sig":
			z.Sig, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Sig")
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
func (z *VoteTicket) Msgsize() (s int) {
	s = 1 + 11 + msgp.ArrayHeaderSize + (len(z.VoteResult) * (msgp.IntSize)) + 7 + msgp.ArrayHeaderSize
	for za0002 := range z.BlockHashes {
		s += msgp.StringPrefixSize + len(z.BlockHashes[za0002])
	}
	s += 6 + msgp.StringPrefixSize + len(z.Voter) + 10 + msgp.Float64Size + 4 + msgp.StringPrefixSize + len(z.Sig)
	return
}
