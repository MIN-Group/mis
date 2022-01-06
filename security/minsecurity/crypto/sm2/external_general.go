// +build !amd64

/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午4:56
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sm2

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	sec "MIS-BC/security/minsecurity"
)

//对空值和曲线参数匹配进行检测
func IsOnCurve(priv *Sm2PrivateKey, pub *Sm2PublicKey) bool {
	if pub == nil || priv == nil || pub.PublicKey == nil || pub.PublicKey.Curve == nil || pub.PublicKey.X == nil || pub.PublicKey.Y == nil ||
		priv.PrivateKey == nil || priv.PrivateKey.D == nil {
		return false
	}

	return priv.PrivateKey.IsOnCurve(pub.PublicKey.X, pub.PublicKey.Y)
}

//随机生成一对公私钥
func GenKeyPair() (*Sm2PrivateKey, *Sm2PublicKey) {
	prk, _ := generateKey(rand.Reader)
	pub := &prk.PublicKey

	prik := &Sm2PrivateKey{prk}
	pubk := &Sm2PublicKey{pub}
	for len(pubk.GetBytes()) != 130 {
		prk, _ := generateKey(rand.Reader)
		pub := &prk.PublicKey

		prik = &Sm2PrivateKey{prk}
		pubk = &Sm2PublicKey{pub}
	}
	return prik, pubk
}

//二次封装的sm2私钥
type Sm2PrivateKey struct {
	PrivateKey *PrivateKey
}

func (pri *Sm2PrivateKey) GetBytes() []byte {
	return []byte(hex.EncodeToString(pri.PrivateKey.D.Bytes()))
}

func (pri *Sm2PrivateKey) SetBytes(priByte []byte) error {
	if len(string(priByte)) != 64 {
		return errors.New("设置私钥长度错误")
	}
	P256Sm2()
	//priByte = []byte("549c0acaaa48db91bc962513e5bf373364b8e64e21b4ed75a6d2d049975e8c18")
	sk, flag := new(big.Int).SetString(string(priByte), 16)
	if flag != true {
		return errors.New("解析字节数组错误")
	}

	pkx, pky := sm2P256.ScalarBaseMult(sk.Bytes())

	priv := PrivateKey{PublicKey{sm2P256, pkx, pky}, sk}

	pri.PrivateKey = &priv
	return nil
}

//封装的sm2公钥
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
	P256Sm2()
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
		Curve: sm2P256,
	}
	pubK.X = x
	pubK.Y = y
	pub.PublicKey = pubK
	return nil
}

func (pub Sm2PublicKey) Encrypt(rand io.Reader, msg []byte, opts sec.DecrypterOpts) (encryptedtext []byte, err error) {
	return Encrypt(pub.PublicKey, msg, rand)
}

func (pri *Sm2PrivateKey) Decrypt(rand io.Reader, msg []byte, opts sec.DecrypterOpts) (plaintext []byte, err error) {
	return pri.PrivateKey.Decrypt(rand, msg, opts)
}

func (pri *Sm2PrivateKey) Sign(rand io.Reader, digest []byte, opts sec.SignerOpts) (signature []byte, err error) {
	return pri.PrivateKey.Sign(rand, digest, nil)

}

func (pub Sm2PublicKey) Verify(msg []byte, sign []byte, opts sec.SignerOpts) (bool, error) {
	return pub.PublicKey.Verify(msg, sign), nil
}
