// +build !amd64

/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/19 上午12:48
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sm2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	priv, pub := GenKeyPair()
	msg := []byte("TestEncryptAndDecrypt : success!")
	enc, err := pub.Encrypt(rand.Reader, msg, nil)
	if err != nil {
		t.Fatalf("encrypt failed:%s", err)
	}
	fmt.Println(hex.EncodeToString(enc))
	dec, err := priv.Decrypt(rand.Reader, enc, nil)
	if err != nil {
		t.Fatalf("decrypt failed:%s", err)
	}
	fmt.Println(dec)
}
