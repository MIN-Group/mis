/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/17 下午7:12
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package identity

import (
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/cert"
	_ "MIS-BC/security/minsecurity/crypto/md5"
	sm2 "MIS-BC/security/minsecurity/crypto/sm2"
	_ "MIS-BC/security/minsecurity/crypto/sm3"
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

func TestSM3(t *testing.T) {
	h := sec.SM3.New()
	h.Write([]byte("123456"))
	res := h.Sum(nil)
	fmt.Println("hash result=", res)
}

func TestCreateIdentity(t *testing.T) {
	id := CreateIdentity("!", KeyParam{PublicKeyAlgorithm: sec.SM2, SignatureAlgorithm: sec.SM2WithSM3})
	fmt.Println(string(id.Prikey.GetBytes()))
	fmt.Println(string(id.Pubkey.GetBytes()))
}

func TestIdentity(t *testing.T) {
	pri, pub := sm2.GenKeyPair()
	id := Identity{Name: "root", Prikey: pri, Pubkey: pub, Passwd: "123456", KeyParam: KeyParam{PublicKeyAlgorithm: sec.SM2, SignatureAlgorithm: sec.SM2WithSM3}}

	cert := cert.Certificate{}
	cert.Version = 1
	cert.SerialNumber = 1
	cert.PublicKey = pub
	cert.SignatureAlgorithm = sec.SM2WithSM3
	cert.PublicKeyAlgorithm = sec.SM2
	cert.IssueTo = "root"
	cert.Issuer = "root"
	cert.NotAfter = time.Now().Unix()
	cert.NotBefore = time.Now().Unix() + 1000
	cert.Timestamp = time.Now().Unix()
	cert.KeyUsage = sec.CertSign
	cert.IsCA = true

	err := cert.SignCert(pri)
	if err != nil {
		fmt.Println(err)
	}
	id.Cert = cert

	sign, err := id.Sign(rand.Reader, []byte("something"), nil)
	if err != nil {
		fmt.Println(err)
	}

	flag, err := id.Verify([]byte("something"), sign, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("verify result=", flag)

	enc, err := id.AsymEncrypt(rand.Reader, []byte("something"), nil)
	dec, err := id.AsymDecrypt(rand.Reader, enc, nil)
	fmt.Println("dev msg", string(dec))

	dump, err := id.Dump(id.Passwd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("dump result=", string(dump))

	id0 := Identity{}
	err = id0.Load(dump, "123456")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id0)
}
