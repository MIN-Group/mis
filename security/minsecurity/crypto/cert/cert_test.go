/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/19 上午4:33
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package cert

import (
	"fmt"
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"testing"
	"time"
)

func TestCertPem(t *testing.T) {
	pri, pub := sm2.GenKeyPair()

	cert := Certificate{}
	cert.Version = 1
	cert.SerialNumber = 1
	cert.PublicKey = pub
	cert.SignatureAlgorithm = sec.SM2WithSM3
	cert.PublicKeyAlgorithm = sec.SM2
	cert.IssueTo = "test"
	cert.Issuer = "test"
	cert.NotAfter = time.Now().Unix()
	cert.NotBefore = time.Now().Unix() + 100
	cert.Timestamp = time.Now().Unix()
	cert.KeyUsage = sec.CertSign
	cert.IsCA = true

	err := cert.SignCert(pri)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cert)

	res, err := Verify(cert, cert)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res=", res)

	pem, err := cert.ToPem([]byte("123"), sec.SM4ECB) //测试Certificate转Pem接口，过程中间接测试parseCertToInnerCert
	if err != nil {
		fmt.Println(err)
	}
	if pem == "" {
		println("no string")
	}
	fmt.Println(pem)

	Cert := Certificate{}
	err = Cert.FromPem(pem, []byte("123"), sec.SM4ECB) //测试Pem转Certificate接口,过程中间接测试parseInnerCertToCert
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Cert)

	res, err = Verify(Cert, Cert)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res=", res)
}
