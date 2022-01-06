// +build amd64

/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午4:56
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sm2

import (
	sec "MIS-BC/security/minsecurity"
	"crypto/rand"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
)

//对空值和曲线参数匹配进行检测
func IsOnCurve(priv *Sm2PrivateKey, pub *Sm2PublicKey) bool {
	if pub == nil || priv == nil || pub.PublicKey == nil || pub.PublicKey.Curve == nil || pub.PublicKey.X == nil || pub.PublicKey.Y == nil ||
		priv.PrivateKey == nil || priv.PrivateKey.D == nil {
		return false
	}

	return priv.PrivateKey.IsOnCurve(pub.PublicKey.X, pub.PublicKey.Y)
}

func GenKeyPair() (*Sm2PrivateKey, *Sm2PublicKey) {
	prk, _ := generateKey(rand.Reader)
	pub := &prk.PublicKey

	prik := &Sm2PrivateKey{prk}
	pubk := &Sm2PublicKey{pub}
	for len(pubk.GetBytes()) != 130 || len(prik.GetBytes()) != 64 {
		prk, _ := generateKey(rand.Reader)
		pub := &prk.PublicKey

		prik = &Sm2PrivateKey{prk}
		pubk = &Sm2PublicKey{pub}
	}
	return prik, pubk
}

type Sm2PrivateKey struct {
	PrivateKey *PrivateKey
}

type sm2Signature struct {
	R, S *big.Int
}

func (pri *Sm2PrivateKey) GetBytes() []byte {
	return []byte(hex.EncodeToString(pri.PrivateKey.D.Bytes()))
}
func (pri *Sm2PrivateKey) SetBytes(priByte []byte) error {
	if len(string(priByte)) != 64 {
		return errors.New("设置私钥长度错误")
	}

	sk, flag := new(big.Int).SetString(string(priByte), 16)
	if flag != true {
		return errors.New("解析字节数组错误")
	}

	pkx, pky := P256().ScalarBaseMult(sk.Bytes())

	priv := PrivateKey{PublicKey{P256(), pkx, pky, nil}, sk, nil}

	pri.PrivateKey = &priv
	return nil
}
func (pri *Sm2PrivateKey) Sign(rand io.Reader, digest []byte, opts sec.SignerOpts) (signature []byte, err error) {
	r, s, err := Sign(rand, pri.PrivateKey, digest)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(sm2Signature{r, s})
}
func (pri *Sm2PrivateKey) Decrypt(rand io.Reader, msg []byte, opts sec.DecrypterOpts) (plaintext []byte, err error) {
	return pri.PrivateKey.Decrypt(rand, msg, opts)
}

type Sm2PublicKey struct {
	PublicKey *PublicKey
}

func (pub Sm2PublicKey) GetBytes() []byte {
	xStr := hex.EncodeToString(pub.PublicKey.X.Bytes())
	yStr := hex.EncodeToString(pub.PublicKey.Y.Bytes())
	return []byte("04" + xStr + yStr)
}
func (pub *Sm2PublicKey) SetBytes(pubBytes []byte) error {
	if len(string(pubBytes)) != 130 {
		return errors.New("公钥长度错误")
	}

	data := string(pubBytes)
	x, flag := new(big.Int).SetString(data[2:66], 16)

	if flag == false {
		errors.New("设置公钥错误")
	}

	y, flag := new(big.Int).SetString(data[66:130], 16)
	if flag != true {
		errors.New("设置公钥错误")
	}

	pubK := &PublicKey{
		Curve: P256(),
	}
	pubK.X = x
	pubK.Y = y
	pub.PublicKey = pubK
	return nil
}
func (pub Sm2PublicKey) Encrypt(rand io.Reader, msg []byte, opts sec.DecrypterOpts) (encryptedtext []byte, err error) {
	return Encrypt(rand, pub.PublicKey, msg)
}
func (pub Sm2PublicKey) Verify(msg []byte, sign []byte, opts sec.SignerOpts) (bool, error) {
	var sm2Sign sm2Signature
	_, err := asn1.Unmarshal(sign, &sm2Sign)
	if err != nil {
		return false, err
	}
	return Verify(pub.PublicKey, msg, sm2Sign.R, sm2Sign.S), nil
}
