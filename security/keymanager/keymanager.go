/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/18 上午1:11
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package keymanager

import (
	"MIS-BC/security"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"MIS-BC/security/minsecurity/crypto/sm4"
	"MIS-BC/utils"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type KeyManager struct {
	prk      *sm2.Sm2PrivateKey
	pub      *sm2.Sm2PublicKey
	keychain *security.KeyChain
}

func (km *KeyManager) Init() {
	km.prk = &sm2.Sm2PrivateKey{}
	km.pub = &sm2.Sm2PublicKey{}
}

func (km *KeyManager) InitFromPem(keyPath string) {
	km.prk, km.pub = sm2.GenKeyPair()
	path_slice := strings.Split(keyPath, ",")
	if len(path_slice) != 2 {
		panic("keyPath格式错误")
	}
	km.ReadKeyFromPem(path_slice[0], path_slice[1])
}

func (km *KeyManager) GenKeyPair() {
	km.prk, km.pub = sm2.GenKeyPair()
}

//得到公钥 类型格式-->字符串格式
func (km *KeyManager) GetPubkey() string {
	return string(km.pub.GetBytes())
}

//得到私钥 类型格式-->字符串格式
func (km *KeyManager) GetPriKey() string {
	return string(km.prk.GetBytes())
}

func (km *KeyManager) SetPubkey(data string) error {
	err := km.pub.SetBytes([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func (km *KeyManager) SetPriKey(data string) error {
	err := km.prk.SetBytes([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func (km *KeyManager) Encrypt(data string) ([]byte, error) {
	return km.pub.Encrypt(rand.Reader, []byte(data), nil)
}

func (km *KeyManager) Decrypt(data string) ([]byte, error) {
	return km.prk.Decrypt(rand.Reader, []byte(data), nil)
}

func (km *KeyManager) IsOnCurve() bool {
	return sm2.IsOnCurve(km.prk, km.pub)
}

//得到数字签名
func (km *KeyManager) Sign(text []byte) (string, error) {
	sig, err := km.prk.Sign(rand.Reader, text, nil)
	return base64.StdEncoding.EncodeToString(sig), err
}

//验证签名
func (km *KeyManager) Verify(text []byte, signature string, pubkey string) (bool, error) {
	t, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("base64编码错误")
		return false, err
	}

	pub := sm2.Sm2PublicKey{}
	err = pub.SetBytes([]byte(pubkey))
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	res, err := pub.Verify(text, t, nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return res, nil
}

func SM4Encrypt(secretkey string, origData []byte) ([]byte, error) {
	return sm4.Sm4Ecb([]byte(secretkey), origData, sm4.ENC)
}
func SM4Decrypt(secretkey string, crypted []byte) ([]byte, error) {
	return sm4.Sm4Ecb([]byte(secretkey), crypted, sm4.DEC)
}

func GetHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	bs := h.Sum(nil)
	return utils.BytesToHex(bs)
}

//读取公私钥
func (km *KeyManager) ReadKeyFromPem(privKeyPath string, pubKeyPath string) {
	prifile, err := os.Open(privKeyPath)
	defer prifile.Close()
	if err != nil {
		fmt.Println("打开私钥文件失败")
	}

	info1, _ := prifile.Stat()
	buf1 := make([]byte, info1.Size())
	prifile.Read(buf1)

	err = km.prk.SetBytes(buf1)
	if err != nil {
		fmt.Println(err)
		panic("读取私钥失败")
	}

	pubfile, err := os.Open(pubKeyPath)
	defer pubfile.Close()
	if err != nil {
		fmt.Println("打开公钥文件失败")
	}

	info2, _ := pubfile.Stat()
	buf2 := make([]byte, info2.Size())
	pubfile.Read(buf2)
	err = km.pub.SetBytes(buf2)
	if err != nil {
		fmt.Println(err)
		panic("读取公钥失败")
	}
}
