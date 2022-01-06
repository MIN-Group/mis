/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午4:45
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package cert

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/sm2"
	_ "MIS-BC/security/minsecurity/crypto/sm3"
	"MIS-BC/security/minsecurity/crypto/sm4"
	"strconv"
	"time"
)

//定义pem和证书结构体之间相互的转换
//func PEMtoCerficate(pem []byte, pass []byte, mode sec.SymmetricAlgorithm) (Certificate, error) {
//
//}
//
//func CertificateToPem(certificate Certificate, pass []byte, mode sec.SymmetricAlgorithm) ([]byte, error) {
//
//}

type CertVersion int

const (
	VERSION1 = iota+1
)

type certificate struct {
	TBSCertificate     tbsCertificate
	SignatureAlgorithm sec.SignatureAlgorithm
	SignatureValue     []byte
}

type tbsCertificate struct {
	Version             int
	SerialNumber        int64
	PublicKey           []byte //公钥
	SignatureAlgorithm  sec.SignatureAlgorithm
	PublicKeyAlgorithm  sec.PublicKeyAlgorithm //公钥算法
	IssueTo             string                 //被签发者
	Issuer              string                 //签发者
	NotBefore, NotAfter int64                  //有效期
	KeyUsage            sec.KeyUsage           //用途
	IsCA                bool                   //是否自签
	Timestamp           int64                  //时间戳

}

type Certificate struct {
	Version             int                    //版本号
	SerialNumber        int64	               //序列号
	PublicKey           sec.PublicKey            //公钥
	Signature           []byte                 //签名
	SignatureAlgorithm  sec.SignatureAlgorithm //签名算法
	PublicKeyAlgorithm  sec.PublicKeyAlgorithm //公钥算法
	IssueTo             string                 //被签发者
	Issuer              string                 //签发者
	NotBefore, NotAfter int64                  //有效期
	KeyUsage            sec.KeyUsage           //用途
	IsCA                bool                   //是否自签
	Timestamp           int64                  //时间戳
}

//
// 证书结构体转换为内部证书结构体
//
// @Description:
// @param cert
// @return error
//
func parseCertToInnerCert(cert *Certificate) certificate {
	tbsCert := tbsCertificate{}
	tbsCert.Version = cert.Version
	tbsCert.SerialNumber = cert.SerialNumber
	tbsCert.PublicKey = cert.PublicKey.GetBytes()
	tbsCert.SignatureAlgorithm = cert.SignatureAlgorithm
	tbsCert.PublicKeyAlgorithm = cert.PublicKeyAlgorithm
	tbsCert.IssueTo = cert.IssueTo
	tbsCert.Issuer = cert.Issuer
	tbsCert.NotBefore = cert.NotBefore
	tbsCert.NotAfter = cert.NotAfter
	tbsCert.KeyUsage = cert.KeyUsage
	tbsCert.IsCA = cert.IsCA
	tbsCert.Timestamp = cert.Timestamp

	certInner := certificate{}
	certInner.SignatureAlgorithm = cert.SignatureAlgorithm
	certInner.SignatureValue = cert.Signature
	certInner.TBSCertificate = tbsCert

	return certInner
}

//
// 内部证书结构体转换为证书结构体
//
// @Description:
// @param cert
// @param c
// @return error
//
func parseInnerCertToCert(cert *certificate, c *Certificate) (error){
	c.Version = cert.TBSCertificate.Version
	c.SerialNumber = cert.TBSCertificate.SerialNumber
	publicKey,err := unMarshalPublicKey(cert.TBSCertificate.PublicKey, cert.TBSCertificate.PublicKeyAlgorithm)
	if err != nil{
		return err
	}
	c.PublicKey = publicKey
	c.Signature = cert.SignatureValue
	c.SignatureAlgorithm = cert.TBSCertificate.SignatureAlgorithm
	c.PublicKeyAlgorithm = cert.TBSCertificate.PublicKeyAlgorithm
	c.IssueTo = cert.TBSCertificate.IssueTo
	c.Issuer = cert.TBSCertificate.Issuer
	c.NotBefore = cert.TBSCertificate.NotBefore
	c.NotAfter = cert.TBSCertificate.NotAfter
	c.KeyUsage = cert.TBSCertificate.KeyUsage
	c.IsCA = cert.TBSCertificate.IsCA
	c.Timestamp = cert.TBSCertificate.Timestamp

	return nil
}

//
// 公钥解码
//
// @Description:
// @param pubByte
// @param algo
// @return PublicKey
// @return error
//
func unMarshalPublicKey(pubByte []byte, algo sec.PublicKeyAlgorithm) (sec.PublicKey, error){
	if pubByte == nil{
		return nil,nil
	}
	switch algo {
	case sec.SM2:
		p := sm2.Sm2PublicKey{}
		err := p.SetBytes(pubByte)
		if err != nil{return nil, err}
		var res sec.PublicKey = &p
		return res, nil
	}
	return nil, nil
}

//
// 使用私钥priv对证书cert进行签名
//
// @Description:
// @receiver cert
// @param priv
// @return error
//
func (cert *Certificate) SignCert(priv sec.PrivateKey) error {
	innercert := parseCertToInnerCert(cert)
	b, err := json.Marshal(innercert.TBSCertificate)
	if err != nil {
		return err
	}
	switch cert.SignatureAlgorithm {
	case sec.SM2WithSM3:
		hashFunc := sec.SM3.New()
		hashFunc.Write(b)
		digest := hashFunc.Sum(nil)
		sig, err := priv.Sign(rand.Reader,digest,nil)
		if err != nil {
			return err
		}
		cert.Signature = sig
		return nil
	default:
		return errors.New("Unknown SignatureAlgorithm")
	}
}


//
// 验证证书合法性，包括任何层面的检查
//
// @Description:
// @param cert
// @param ca
// @return bool
// @return error
//
func Verify(cert, ca Certificate) (bool, error) {
	myInnercert := parseCertToInnerCert(&cert)
	switch myInnercert.TBSCertificate.Version {
	case VERSION1:
		err := CheckDuration(myInnercert)
		if err != nil {
			return false, err
		}
		if cert.IsCA {
			return CheckSign(myInnercert,cert.PublicKey)
		} else {
			return CheckSign(myInnercert,ca.PublicKey)
		}
	default:
		return false, CertificateInvalidError{
			Reason: IncompatibleVersion,
			Detail: "",
		}
	}
}

//
// 证书有效期检查
//
// @Description:
// @param cert
// @return error
//
func CheckDuration (cert certificate) error {
	currTime := time.Now().UTC().Unix()
	if cert.TBSCertificate.NotAfter > currTime {
		return CertificateInvalidError{
			Reason: NotReachEffectiveDate,
			Detail: "the effective date is " + time.Unix(cert.TBSCertificate.NotAfter,0).Format("2006-01-02 15:04:05"),
		}
	}
	if currTime > cert.TBSCertificate.NotBefore {
		return CertificateInvalidError{
			Reason: Expired,
			Detail: "the expiration date is " + time.Unix(cert.TBSCertificate.NotBefore,0).Format("2006-01-02 15:04:05"),
		}
	}
	return nil
}

//
// 证书签名检查
//
// @Description:
// @param cert
// @param pub
// @return bool
// @return error
//
func CheckSign(cert certificate, pub sec.PublicKey) (bool, error) {
	switch cert.TBSCertificate.SignatureAlgorithm {
	case sec.SM2WithSM3:
		b, err := json.Marshal(cert.TBSCertificate)
		if err != nil {
			return false, CertificateInvalidError{
				Reason: MarshalError,
				Detail: err.Error(),
			}
		}
		hashFunc := sec.SM3.New()
		hashFunc.Write(b)
		digest := hashFunc.Sum(nil)
		ret, err := pub.Verify(digest,cert.SignatureValue,nil)
		if err != nil {
			return false, CertificateInvalidError{
				Reason: SignatureError,
				Detail: err.Error(),
			}
		}
		if !ret {
			return false, CertificateInvalidError{
				Reason: CANotAuthorizedForThisName,
				Detail: "",
			}
		}
		return true, nil
	default:
		return false, CertificateInvalidError{
			Reason: UnknownSignatureAlgorithm,
			Detail: strconv.Itoa(int(cert.TBSCertificate.SignatureAlgorithm)),
		}
	}
}

//
// 证书结构体转换为pem
//
// @Description:
// @receiver cert
// @param passwd
// @param mode
// @return string
// @return error
//
func (cert *Certificate) ToPem(passwd []byte, mode sec.SymmetricAlgorithm) (string, error) {
	if cert == nil || cert.Signature == nil {
		return "", errors.New("Wrong certificate")
	}
	b, err := json.Marshal(parseCertToInnerCert(cert))
	if err != nil {
		return "", errors.New("json marshal fails")
	}
	var ret []byte
	if len(passwd) > 0 {
		hashFunc := sec.SM3.New()
		hashFunc.Write(passwd)
		passHash := hashFunc.Sum(nil)
		if len(passHash) == 32 {
			for i := 0; i < 16; i++ {
				passHash[i] += passHash[i+16]
			}
		}
		switch mode {
		case sec.SM4ECB:
			ret, err = sm4.Sm4Ecb(passHash[:16], b, sm4.ENC)
			if err != nil {
				return "", err
			}
		case sec.SM4CBC:
			ret, err = sm4.Sm4Cbc(passHash[:16], b, sm4.ENC)
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("Unknown SymmetricAlgorithm")
		}
	} else {
		ret = b
	}
	return base64.StdEncoding.EncodeToString(ret), nil
}

//
// pem转换为证书结构体
//
// @Description:
// @receiver cert
// @param pemStr
// @param passwd
// @param mode
// @return error
//
func (cert *Certificate) FromPem(pemStr string, passwd []byte, mode sec.SymmetricAlgorithm) error {
	if pemStr == "" {
		return errors.New("Wrong certificate pem")
	}
	pemByte, err := base64.StdEncoding.DecodeString(pemStr)
	if err != nil {
		return err
	}

	var decCert []byte
	if len(passwd) > 0 {
		hashFunc := sec.SM3.New()
		hashFunc.Write(passwd)
		passHash := hashFunc.Sum(nil)
		if len(passHash) == 32 {
			for i := 0; i < 16; i++ {
				passHash[i] += passHash[i+16]
			}
		}
		switch mode {
		case sec.SM4ECB:
			decCert, err = sm4.Sm4Ecb(passHash[:16], pemByte, sm4.DEC)
			if err != nil {
				return  err
			}
		case sec.SM4CBC:
			decCert, err = sm4.Sm4Cbc(passHash[:16], pemByte, sm4.DEC)
			if err != nil {
				return err
			}
		default:
			return errors.New("Unknown SymmetricAlgorithm")
		}
	} else {
		decCert = pemByte
	}

	c := new(certificate)
	err = json.Unmarshal(decCert, c)
	if err != nil {
		return errors.New("wrong password")
	}
	err = parseInnerCertToCert(c, cert)
	if err != nil {
		return err
	}
	return nil
}
