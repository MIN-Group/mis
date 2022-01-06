/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/18 上午4:05
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package crypto

import (
	"crypto/rand"
	"fmt"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"testing"
)

func TestKeyGenPair(t *testing.T) {
	priv, pub := sm2.GenKeyPair()

	priv0 := sm2.Sm2PrivateKey{}
	err0 := priv0.SetBytes([]byte("549c0acaaa48db91bc962513e5bf373364b8e64e21b4ed75a6d2d049975e8c18"))
	fmt.Println(string(priv0.GetBytes()))
	fmt.Println(err0)

	sig, err := priv0.Sign(rand.Reader, []byte("xzw"), nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := pub.Verify([]byte("xzw"), sig, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	enc, err := pub.Encrypt(rand.Reader, []byte("xzw"), nil)
	result, err := priv.Decrypt(rand.Reader, enc, nil)
	fmt.Println(string(result))
	pub = nil
	fmt.Println(sm2.IsOnCurve(priv, pub))
}
