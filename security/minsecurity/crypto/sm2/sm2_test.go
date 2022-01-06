package sm2

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func BenchmarkSm2Sign(b *testing.B) {
	pri, _ := GenKeyPair()
	fmt.Println(pri.PrivateKey.X.String())
	fmt.Println(len(pri.PrivateKey.X.String()))
	data := make([]byte, 8800)
	rand.Read(data)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pri.Sign(nil, data[:], nil)
		if err != nil {
			b.Errorf("S M2签名错误")
		}
	}
}

func BenchmarkSm2Verify(b *testing.B) {
	pri, pub := GenKeyPair()
	data := make([]byte, 8800)
	rand.Read(data)
	signResult, err := pri.Sign(rand.Reader, data[:], nil)
	if err != nil {
		b.Errorf("SM签名错误")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := pub.Verify(data[:], signResult[:], nil)
		if res == false || err != nil {
			b.Errorf("SM2签名错误")
		}
	}
}