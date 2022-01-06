/**
 * @Author: wzx
 * @Description:身份定义
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午5:48
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package identity

import (
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	_ "MIS-BC/security/minsecurity/crypto/sm3"
	"MIS-BC/security/minsecurity/crypto/sm4"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// KeyParam 身份采用的密码学算法
type KeyParam struct {
	PublicKeyAlgorithm sec.PublicKeyAlgorithm
	SignatureAlgorithm sec.SignatureAlgorithm
}

type Identity struct {
	Name          string           //用户名
	KeyParam      KeyParam         //身份采用的密码学算法
	Prikey        sec.PrivateKey   //私钥
	PrikeyRawByte []byte           //加密后的私钥
	Pubkey        sec.PublicKey    //公钥
	Cert          cert.Certificate //证书
	Passwd        string
}

func CreateIdentity(name string, param KeyParam) Identity {
	id := Identity{}
	id.Name = name
	switch param.PublicKeyAlgorithm {
	case sec.SM2:
		id.Prikey, id.Pubkey = sm2.GenKeyPair()
	}
	return id
}

func (i *Identity) String() string {
	return fmt.Sprintf("{Name:%v}", i.Name)
}

//内部转码用的结构体
type innerIdentity struct {
	Name          string
	KeyParam      KeyParam
	Prikey        []byte
	Pubkey        []byte
	Cert          string
	PrikeyRawByte []byte //加密后的私钥
}

//
// 身份结构体转换为内部身份结构体
//
// @Description:
// @param Identity
// @return innerIdentity
// @return error
//
func parseIdentityToInner(id *Identity) (innerIdentity, error) {
	inner := innerIdentity{}
	inner.Name = id.Name
	inner.KeyParam = id.KeyParam
	inner.PrikeyRawByte = id.PrikeyRawByte
	if id.Prikey != nil {
		inner.Prikey = id.Prikey.GetBytes()
	}
	if id.Pubkey != nil {
		inner.Pubkey = id.Pubkey.GetBytes()
	}

	passwd := "yzytql"
	if id.Cert.IssueTo != "" {
		var err error
		inner.Cert, err = id.Cert.ToPem([]byte(passwd), sec.SM4CBC)

		if err != nil {
			return inner, err
		}
	}

	return inner, nil
}

//
// 内部身份结构体转换为身份结构体
//
// @Description:
// @param inner
// @param id
// @return error
//
func parseInnerToIdentity(inner *innerIdentity, id *Identity) error {
	id.Name = inner.Name
	id.KeyParam = inner.KeyParam
	if inner.Pubkey != nil {
		var err error
		id.Pubkey, err = crypto.UnMarshalPublicKey(inner.Pubkey, inner.KeyParam.PublicKeyAlgorithm)
		if err != nil {
			return err
		}
	}

	if inner.Prikey != nil {
		var err error
		id.Prikey, err = crypto.UnMarshalPrivateKey(inner.Prikey, inner.KeyParam.PublicKeyAlgorithm)
		if err != nil {
			return err
		}
	}
	id.PrikeyRawByte = inner.PrikeyRawByte

	passwd := "yzytql"
	if inner.Cert != "" {
		err := id.Cert.FromPem(inner.Cert, []byte(passwd), sec.SM4CBC)
		if err != nil {
			return err
		}
	}

	return nil
}

// Sign
// 使用Identity对象进行签名
//
// @Description:
// @receiver id
// @param rand
// @param digest
// @param opts
// @return []byte
// @return error
//
func (id *Identity) Sign(rand io.Reader, digest []byte, opts sec.SignerOpts) ([]byte, error) {
	if id.Prikey == nil {
		return nil, errors.New("Invalid Prikey")
	}

	switch id.KeyParam.PublicKeyAlgorithm {
	case sec.SM2:
		p, flag := id.Prikey.(*sm2.Sm2PrivateKey)
		if flag != true {
			return nil, errors.New("Invalid Prikey")
		}
		switch id.KeyParam.SignatureAlgorithm {
		case sec.SM2WithSM3:
			return p.Sign(rand, digest, nil)
		default:
			return nil, errors.New("Can't find the sign algorithm")
		}
	}

	return nil, errors.New("Can't find key type")
}

// Verify
// 使用Identity对象进行签名验证
//
// @Description:
// @receiver id
// @param msg
// @param sign
// @param opts
// @return bool
// @return error
//
func (id *Identity) Verify(msg []byte, sign []byte, opts sec.SignerOpts) (bool, error) {
	if id.Pubkey == nil {
		return false, errors.New("Invalid Pubkey")
	}
	switch id.KeyParam.PublicKeyAlgorithm {
	case sec.SM2:
		p, flag := id.Pubkey.(*sm2.Sm2PublicKey)
		if flag != true {
			return false, errors.New("Invalid Pubkey")
		}
		switch id.KeyParam.SignatureAlgorithm {
		case sec.SM2WithSM3:
			return p.Verify(msg, sign, opts)
		default:
			return false, errors.New("Can't find the sign algorithm")
		}
	}

	return false, errors.New("Can't find key type")
}

// AsymDecrypt
// 使用Identity进行私钥解密
//
// @Description:
// @receiver id
// @param rand
// @param msg
// @param opts
// @return []byte
// @return error
//
func (id *Identity) AsymDecrypt(rand io.Reader, msg []byte, opts sec.DecrypterOpts) ([]byte, error) {
	if id.Prikey == nil {
		return nil, errors.New("Invalid Prikey")
	}
	return id.Prikey.Decrypt(rand, msg, opts)
}

// AsymEncrypt
// 使用Identity进行公钥加密
//
// @Description:
// @receiver id
// @param rand
// @param mgmt_data
// @param opts
// @return []byte
// @return error
//
func (id *Identity) AsymEncrypt(rand io.Reader, data []byte, opts sec.DecrypterOpts) ([]byte, error) {
	if id.Pubkey == nil {
		return nil, errors.New("Invalid Pubkey")
	}
	return id.Pubkey.Encrypt(rand, data, opts)
}

// Dump
// 将Identity对象导出到字节数组
//
// @Description:
// @receiver id
// @param passwd
// @return []byte
// @return error
//
func (id *Identity) Dump(passwd string) ([]byte, error) {
	inner, err := parseIdentityToInner(id)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(inner)
	if err != nil {
		return nil, err
	}
	var ret []byte
	if len(passwd) > 0 {
		hashFunc := sec.SM3.New()
		hashFunc.Write([]byte(passwd))
		passHash := hashFunc.Sum(nil)
		if len(passHash) == 32 {
			for i := 0; i < 16; i++ {
				passHash[i] += passHash[i+16]
			}
		}
		ret, err = sm4.Sm4Cbc(passHash[:16], b, sm4.ENC)
		if err != nil {
			return nil, err
		}
	} else {
		ret = b
	}

	return []byte(base64.StdEncoding.EncodeToString(ret)), nil
}

func (id *Identity) GetJsonString() ([]byte, error) {
	inner, err := parseIdentityToInner(id)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(inner)
	return b, nil
}

// Load
// 从字节数组加载Identity对象
//
// @Description:
// @receiver id
// @param mgmt_data
// @param passwd
// @return error
//
func (id *Identity) Load(data []byte, passwd string) error {
	idByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}

	var plaintext []byte
	if len(passwd) > 0 {
		hashFunc := sec.SM3.New()
		hashFunc.Write([]byte(passwd))
		passHash := hashFunc.Sum(nil)
		if len(passHash) == 32 {
			for i := 0; i < 16; i++ {
				passHash[i] += passHash[i+16]
			}
		}
		var err error
		plaintext, err = sm4.Sm4Cbc(passHash[:16], idByte, sm4.DEC)
		if err != nil {
			return err
		}
	} else {
		plaintext = idByte
	}

	inner := innerIdentity{}
	err = json.Unmarshal(plaintext, &inner)
	if err != nil {
		return errors.New("wrong password")
	}

	err = parseInnerToIdentity(&inner, id)
	if err != nil {
		return err
	}

	return nil
}

func (id *Identity) Lock(passwd string, algo sec.SymmetricAlgorithm) (bool, error) {
	switch algo {
	case sec.SM4ECB:
		if len(passwd) == 0 {
			return false, errors.New("Invalid passwd length")
		}
		hashFunc := sec.SM3.New()
		hashFunc.Write([]byte(passwd))
		passHash := hashFunc.Sum(nil)
		if len(passHash) == 32 {
			for i := 0; i < 16; i++ {
				passHash[i] += passHash[i+16]
			}
		}
		_, priByte := crypto.MarshalPrivateKey(id.Prikey)
		encMsg, err := sm4.Sm4Ecb(passHash[:16], priByte, sm4.ENC)
		if err != nil {
			return false, err
		}
		id.PrikeyRawByte = encMsg
		id.Prikey = nil
		return true, nil
	}

	return false, errors.New("Unspported algorithm")
}

func (id *Identity) IsLocked() bool {
	if len(id.PrikeyRawByte) > 0 {
		return true
	}
	return false
}

func (id *Identity) UnLock(passwd string, algo sec.SymmetricAlgorithm) (bool, error) {
	switch algo {
	case sec.SM4ECB:
		if len(passwd) == 0 {
			return false, errors.New("Invalid passwd length")
		}

		if !(id.PrikeyRawByte == nil || len(id.PrikeyRawByte) == 0) {
			hashFunc := sec.SM3.New()
			hashFunc.Write([]byte(passwd))
			passHash := hashFunc.Sum(nil)
			if len(passHash) == 32 {
				for i := 0; i < 16; i++ {
					passHash[i] += passHash[i+16]
				}
			}
			dec, err := sm4.Sm4Ecb(passHash[:16], id.PrikeyRawByte, sm4.DEC)
			if err != nil {
				return false, err
			}
			if dec == nil {
				return false, errors.New("Passwd error, decrypt failed!")
			}
			priKey, err := crypto.UnMarshalPrivateKey(dec, id.KeyParam.PublicKeyAlgorithm)
			if err != nil {
				return false, err
			}
			id.Prikey = priKey
			id.PrikeyRawByte = nil
			return true, nil
		} else {
			priKey, err := crypto.UnMarshalPrivateKey(id.PrikeyRawByte, id.KeyParam.PublicKeyAlgorithm)
			if err != nil {
				return false, err
			}
			id.Prikey = priKey
			id.PrikeyRawByte = nil
			return true, nil
		}
	}
	return false, errors.New("Unspported algorithm")
}

// HashPrivateKey
// 判断当前网络身份是否包含私钥
//
// @Description:
// @receiver id
// @return bool
//
func (id *Identity) HashPrivateKey() bool {
	return id.Prikey != nil || len(id.PrikeyRawByte) > 0
}
