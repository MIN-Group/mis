package MetaData

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *GenesisTransaction) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "WorkerNum":
			z.WorkerNum, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "VotedNum":
			z.VotedNum, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "BlockGroupPerCycle":
			z.BlockGroupPerCycle, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Tcut":
			z.Tcut, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "WorkerPubList":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.WorkerPubList == nil && msz > 0 {
				z.WorkerPubList = make(map[string]uint64, msz)
			} else if len(z.WorkerPubList) > 0 {
				for key, _ := range z.WorkerPubList {
					delete(z.WorkerPubList, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg uint64
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				bzg, err = dc.ReadUint64()
				if err != nil {
					return
				}
				z.WorkerPubList[xvk] = bzg
			}
		case "WorkerCandidatePubList":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.WorkerCandidatePubList == nil && msz > 0 {
				z.WorkerCandidatePubList = make(map[string]uint64, msz)
			} else if len(z.WorkerCandidatePubList) > 0 {
				for key, _ := range z.WorkerCandidatePubList {
					delete(z.WorkerCandidatePubList, key)
				}
			}
			for msz > 0 {
				msz--
				var bai string
				var cmr uint64
				bai, err = dc.ReadString()
				if err != nil {
					return
				}
				cmr, err = dc.ReadUint64()
				if err != nil {
					return
				}
				z.WorkerCandidatePubList[bai] = cmr
			}
		case "VoterPubList":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.VoterPubList == nil && msz > 0 {
				z.VoterPubList = make(map[string]uint64, msz)
			} else if len(z.VoterPubList) > 0 {
				for key, _ := range z.VoterPubList {
					delete(z.VoterPubList, key)
				}
			}
			for msz > 0 {
				msz--
				var ajw string
				var wht uint64
				ajw, err = dc.ReadString()
				if err != nil {
					return
				}
				wht, err = dc.ReadUint64()
				if err != nil {
					return
				}
				z.VoterPubList[ajw] = wht
			}
		case "WNS":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.WorkerSet) >= int(xsz) {
				z.WorkerSet = z.WorkerSet[:xsz]
			} else {
				z.WorkerSet = make([]string, xsz)
			}
			for hct := range z.WorkerSet {
				z.WorkerSet[hct], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "VS":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.VoterSet) >= int(xsz) {
				z.VoterSet = z.VoterSet[:xsz]
			} else {
				z.VoterSet = make([]string, xsz)
			}
			for cua := range z.VoterSet {
				z.VoterSet[cua], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *GenesisTransaction) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "WorkerNum"
	err = en.Append(0x89, 0xa9, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x75, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.WorkerNum)
	if err != nil {
		return
	}
	// write "VotedNum"
	err = en.Append(0xa8, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.VotedNum)
	if err != nil {
		return
	}
	// write "BlockGroupPerCycle"
	err = en.Append(0xb2, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x50, 0x65, 0x72, 0x43, 0x79, 0x63, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.BlockGroupPerCycle)
	if err != nil {
		return
	}
	// write "Tcut"
	err = en.Append(0xa4, 0x54, 0x63, 0x75, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Tcut)
	if err != nil {
		return
	}
	// write "WorkerPubList"
	err = en.Append(0xad, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.WorkerPubList)))
	if err != nil {
		return
	}
	for xvk, bzg := range z.WorkerPubList {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = en.WriteUint64(bzg)
		if err != nil {
			return
		}
	}
	// write "WorkerCandidatePubList"
	err = en.Append(0xb6, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.WorkerCandidatePubList)))
	if err != nil {
		return
	}
	for bai, cmr := range z.WorkerCandidatePubList {
		err = en.WriteString(bai)
		if err != nil {
			return
		}
		err = en.WriteUint64(cmr)
		if err != nil {
			return
		}
	}
	// write "VoterPubList"
	err = en.Append(0xac, 0x56, 0x6f, 0x74, 0x65, 0x72, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.VoterPubList)))
	if err != nil {
		return
	}
	for ajw, wht := range z.VoterPubList {
		err = en.WriteString(ajw)
		if err != nil {
			return
		}
		err = en.WriteUint64(wht)
		if err != nil {
			return
		}
	}
	// write "WNS"
	err = en.Append(0xa3, 0x57, 0x4e, 0x53)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.WorkerSet)))
	if err != nil {
		return
	}
	for hct := range z.WorkerSet {
		err = en.WriteString(z.WorkerSet[hct])
		if err != nil {
			return
		}
	}
	// write "VS"
	err = en.Append(0xa2, 0x56, 0x53)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.VoterSet)))
	if err != nil {
		return
	}
	for cua := range z.VoterSet {
		err = en.WriteString(z.VoterSet[cua])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GenesisTransaction) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "WorkerNum"
	o = append(o, 0x89, 0xa9, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.WorkerNum)
	// string "VotedNum"
	o = append(o, 0xa8, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.VotedNum)
	// string "BlockGroupPerCycle"
	o = append(o, 0xb2, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x50, 0x65, 0x72, 0x43, 0x79, 0x63, 0x6c, 0x65)
	o = msgp.AppendInt(o, z.BlockGroupPerCycle)
	// string "Tcut"
	o = append(o, 0xa4, 0x54, 0x63, 0x75, 0x74)
	o = msgp.AppendFloat64(o, z.Tcut)
	// string "WorkerPubList"
	o = append(o, 0xad, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	o = msgp.AppendMapHeader(o, uint32(len(z.WorkerPubList)))
	for xvk, bzg := range z.WorkerPubList {
		o = msgp.AppendString(o, xvk)
		o = msgp.AppendUint64(o, bzg)
	}
	// string "WorkerCandidatePubList"
	o = append(o, 0xb6, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	o = msgp.AppendMapHeader(o, uint32(len(z.WorkerCandidatePubList)))
	for bai, cmr := range z.WorkerCandidatePubList {
		o = msgp.AppendString(o, bai)
		o = msgp.AppendUint64(o, cmr)
	}
	// string "VoterPubList"
	o = append(o, 0xac, 0x56, 0x6f, 0x74, 0x65, 0x72, 0x50, 0x75, 0x62, 0x4c, 0x69, 0x73, 0x74)
	o = msgp.AppendMapHeader(o, uint32(len(z.VoterPubList)))
	for ajw, wht := range z.VoterPubList {
		o = msgp.AppendString(o, ajw)
		o = msgp.AppendUint64(o, wht)
	}
	// string "WNS"
	o = append(o, 0xa3, 0x57, 0x4e, 0x53)
	o = msgp.AppendArrayHeader(o, uint32(len(z.WorkerSet)))
	for hct := range z.WorkerSet {
		o = msgp.AppendString(o, z.WorkerSet[hct])
	}
	// string "VS"
	o = append(o, 0xa2, 0x56, 0x53)
	o = msgp.AppendArrayHeader(o, uint32(len(z.VoterSet)))
	for cua := range z.VoterSet {
		o = msgp.AppendString(o, z.VoterSet[cua])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GenesisTransaction) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "WorkerNum":
			z.WorkerNum, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "VotedNum":
			z.VotedNum, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "BlockGroupPerCycle":
			z.BlockGroupPerCycle, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "Tcut":
			z.Tcut, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "WorkerPubList":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.WorkerPubList == nil && msz > 0 {
				z.WorkerPubList = make(map[string]uint64, msz)
			} else if len(z.WorkerPubList) > 0 {
				for key, _ := range z.WorkerPubList {
					delete(z.WorkerPubList, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg uint64
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bzg, bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					return
				}
				z.WorkerPubList[xvk] = bzg
			}
		case "WorkerCandidatePubList":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.WorkerCandidatePubList == nil && msz > 0 {
				z.WorkerCandidatePubList = make(map[string]uint64, msz)
			} else if len(z.WorkerCandidatePubList) > 0 {
				for key, _ := range z.WorkerCandidatePubList {
					delete(z.WorkerCandidatePubList, key)
				}
			}
			for msz > 0 {
				var bai string
				var cmr uint64
				msz--
				bai, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				cmr, bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					return
				}
				z.WorkerCandidatePubList[bai] = cmr
			}
		case "VoterPubList":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.VoterPubList == nil && msz > 0 {
				z.VoterPubList = make(map[string]uint64, msz)
			} else if len(z.VoterPubList) > 0 {
				for key, _ := range z.VoterPubList {
					delete(z.VoterPubList, key)
				}
			}
			for msz > 0 {
				var ajw string
				var wht uint64
				msz--
				ajw, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				wht, bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					return
				}
				z.VoterPubList[ajw] = wht
			}
		case "WNS":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.WorkerSet) >= int(xsz) {
				z.WorkerSet = z.WorkerSet[:xsz]
			} else {
				z.WorkerSet = make([]string, xsz)
			}
			for hct := range z.WorkerSet {
				z.WorkerSet[hct], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "VS":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.VoterSet) >= int(xsz) {
				z.VoterSet = z.VoterSet[:xsz]
			} else {
				z.VoterSet = make([]string, xsz)
			}
			for cua := range z.VoterSet {
				z.VoterSet[cua], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *GenesisTransaction) Msgsize() (s int) {
	s = 1 + 10 + msgp.IntSize + 9 + msgp.IntSize + 19 + msgp.IntSize + 5 + msgp.Float64Size + 14 + msgp.MapHeaderSize
	if z.WorkerPubList != nil {
		for xvk, bzg := range z.WorkerPubList {
			_ = bzg
			s += msgp.StringPrefixSize + len(xvk) + msgp.Uint64Size
		}
	}
	s += 23 + msgp.MapHeaderSize
	if z.WorkerCandidatePubList != nil {
		for bai, cmr := range z.WorkerCandidatePubList {
			_ = cmr
			s += msgp.StringPrefixSize + len(bai) + msgp.Uint64Size
		}
	}
	s += 13 + msgp.MapHeaderSize
	if z.VoterPubList != nil {
		for ajw, wht := range z.VoterPubList {
			_ = wht
			s += msgp.StringPrefixSize + len(ajw) + msgp.Uint64Size
		}
	}
	s += 4 + msgp.ArrayHeaderSize
	for hct := range z.WorkerSet {
		s += msgp.StringPrefixSize + len(z.WorkerSet[hct])
	}
	s += 3 + msgp.ArrayHeaderSize
	for cua := range z.VoterSet {
		s += msgp.StringPrefixSize + len(z.VoterSet[cua])
	}
	return
}
