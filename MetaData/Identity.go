package MetaData

import (
	"MIS-BC/common"
	"MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"MIS-BC/security/minsecurity/identity"
	"fmt"
)

//公钥算法
type PublicKeyAlgorithm int

const (
	SM2 PublicKeyAlgorithm = iota
)

//签名算法
type SignatureAlgorithm int

const (
	SM2WithSM3 SignatureAlgorithm = iota
)

//身份采用的密码学算法
type KeyParam struct {
	PublicKeyAlgorithm PublicKeyAlgorithm
	SignatureAlgorithm SignatureAlgorithm
}

type ModifyRecord struct {
	Type      string
	Command   string
	Timestamp string
}

type DecrypterOpts interface{}

//go:generate msgp
type Identity struct {
	// 不可变部分
	KeyParam           KeyParam `msg:"keyparam"`           //身份采用的密码学算法
	IdentityIdentifier string   `msg:"identityidentifier"` //身份标识
	Pubkey             string   `msg:"pubkey"`             //公钥
	Cert               string   `msg:"cert"`               //用户证书
	Timestamp          string   `msg:"timestamp"`          //注册时间

	// 可变部分
	Type          string         `msg:"type"`         //身份操作类型
	Command       string         `msg:"command"`      //身份操作命令
	Passwd        string         `msg:"passwd"`       //密码，用于给私钥加密
	IPIdentifier  string         `msg:"ipidentifier"` //IP标识
	ModifyRecords []ModifyRecord `msg:"modifyrecord"` //修改记录
	IsValid       int            `msg:"isvalid"`      //身份有效性
}

func (id Identity) ToByteArray() []byte {
	data, _ := id.MarshalMsg(nil)
	return data
}

func (id *Identity) FromByteArray(data []byte) {
	_, err := id.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}

func (id *Identity) ParseBCIdentityToCommon() identity.Identity {
	cidentity := identity.Identity{}

	if len(id.IdentityIdentifier) <= 0 {
		common.Logger.Error("身份标识为空")
	} else if id.IdentityIdentifier[0] != '/' {
		cidentity.Name = "/" + id.IdentityIdentifier
		common.Logger.Warn("Invalid identifierString \"%s\", require start with / ", id.IdentityIdentifier)
	} else {
		cidentity.Name = id.IdentityIdentifier
	}
	cidentity.KeyParam = identity.KeyParam{minsecurity.PublicKeyAlgorithm(id.KeyParam.PublicKeyAlgorithm), minsecurity.SignatureAlgorithm(id.KeyParam.SignatureAlgorithm)}

	Ca := cert.Certificate{}
	err := Ca.FromPem(id.Cert, []byte(id.Passwd), minsecurity.SM4ECB)
	if err != err {
		common.Logger.Error(err)
	}
	cidentity.Cert = Ca

	//pub := new(sm2.Sm2PublicKey)
	//pub.SetBytes([]byte(id.Pubkey))
	//var pubkey minsecurity.PublicKey = pub
	cidentity.Pubkey = Ca.PublicKey

	cidentity.Prikey = nil

	cidentity.PrikeyRawByte = []byte("")

	cidentity.Passwd = id.Passwd

	return cidentity
}

func ParseCommonToBCIdentity(cid *identity.Identity) Identity {
	id := Identity{}

	id.IdentityIdentifier = cid.Name
	id.KeyParam = KeyParam{PublicKeyAlgorithm(cid.KeyParam.PublicKeyAlgorithm), SignatureAlgorithm(cid.KeyParam.SignatureAlgorithm)}

	c, err := cid.Cert.ToPem([]byte(cid.Passwd), 0)
	if err != err {
		common.Logger.Error(err)
	}
	id.Cert = c

	pub := sm2.Sm2PublicKey{}
	pub.SetBytes([]byte(id.Pubkey))
	var pubkey minsecurity.PublicKey = &pub
	cid.Pubkey = pubkey

	id.Pubkey = string(cid.Pubkey.GetBytes())

	id.Passwd = cid.Passwd

	return id
}
