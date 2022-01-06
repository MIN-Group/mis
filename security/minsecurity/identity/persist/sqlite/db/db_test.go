package db

import (
	"fmt"
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"MIS-BC/security/minsecurity/identity"
	"MIS-BC/security/minsecurity/identity/persist/sqlite"
	"strconv"
	"sync"

	"testing"
	"time"
)

//func TestDatabase(t *testing.T) {
//	//测试PersistIdentity
//	pri, pub := sm2_general.GenKeyPair()
//	sqlite.OpenDefault()
//	id := identity.Identity{}
//	id.Name = "xzw1"
//	id.KeyParam.PublicKeyAlgorithm = 0
//	id.KeyParam.SignatureAlgorithm = 0
//	id.Prikey = pri
//	id.Pubkey = pub
//	id.Passwd = "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
//	cert := cert.Certificate{}
//	cert.Version = 1
//	cert.SerialNumber = 1
//	cert.PublicKey = pub
//	cert.SignatureAlgorithm = sec.SM2WithSM3
//	cert.PublicKeyAlgorithm = sec.SM2
//	cert.IssueTo = "xzw1"
//	cert.Issuer = "root"
//	cert.NotAfter = time.Now().Unix()
//	cert.NotBefore = time.Now().Unix() + 1000
//	cert.Timestamp = time.Now().Unix()
//	cert.KeyUsage = sec.CertSign
//	cert.IsCA = true
//
//	err := cert.SignCert(pri)
//	if err != nil {
//		fmt.Println(err)
//	}
//	id.Cert = cert
//	res, err := PersistIdentity(&id)
//	if res != true {
//		t.Log(err)
//	}
//	//测试GetAllIdentityFromStorage
//	result1, err := GetAllIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印所有身份：%v \n", result1)
//	//测试GetIdentityByNameFromStorag
//	result2, err := GetIdentityByNameFromStorage("xzw1")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印指定身份：%v \n", result2)
//	//测试SetDefaultIdentityFromStorage
//	result3, err := SetDefaultIdentityByNameInStorage("xzw1")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印default设置结果：%v \n", result3)
//	//测试etDefaultIdentityFromStorage
//	result4, err := GetDefaultIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印default身份：%v \n", result4)
//	//测试DeleteIdentityByName
//	result5, err := DeleteIdentityByName("xzw1")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印删除结果：%v \n", result5)
//	result6, err := GetAllIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印所有身份：%v \n", result6)
//}

//func TestDatabaseWithHighLoad(t *testing.T) {
//	sqlite.OpenDefault()
//	//测试PersistIdentity
//	for i := 1; i <= 300; i++ {
//		pri, pub := sm2_general.GenKeyPair()
//		id := identity.Identity{}
//		id.Name = "xzw" + strconv.Itoa(i)
//		id.KeyParam.PublicKeyAlgorithm = 0
//		id.KeyParam.SignatureAlgorithm = 0
//		id.Prikey = pri
//		id.Pubkey = pub
//		id.Passwd = "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
//		cert := cert.Certificate{}
//		cert.Version = 1
//		cert.SerialNumber = 1
//		cert.PublicKey = pub
//		cert.SignatureAlgorithm = sec.SM2WithSM3
//		cert.PublicKeyAlgorithm = sec.SM2
//		cert.IssueTo = "root"
//		cert.Issuer = "root"
//		cert.NotAfter = time.Now().Unix()
//		cert.NotBefore = time.Now().Unix() + 1000
//		cert.Timestamp = time.Now().Unix()
//		cert.KeyUsage = sec.CertSign
//		cert.IsCA = true
//
//		err := cert.SignCert(pri)
//		if err != nil {
//			fmt.Println(err)
//		}
//		id.Cert = cert
//		res, err := PersistIdentity(&id)
//		if res != true {
//			t.Log(err)
//		}
//	}
//
//	//测试GetAllIdentityFromStorage
//	result1, err := GetAllIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印所有身份：%v \n", result1)
//	//测试GetIdentityByNameFromStorag
//	result2, err := GetIdentityByNameFromStorage("xzw1")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印指定身份：%v \n", result2)
//	//测试SetDefaultIdentityFromStorage
//	result3, err := SetDefaultIdentityByNameInStorage("xzw1")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印default设置结果：%v \n", result3)
//	//测试etDefaultIdentityFromStorage
//	result4, err := GetDefaultIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印default身份：%v \n", result4)
//	//测试DeleteIdentityByName
//	result5, err := DeleteIdentityByName("xzw2")
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印删除结果：%v \n", result5)
//	result6, err := GetAllIdentityFromStorage()
//	if err != nil {
//		t.Log(err)
//	}
//	fmt.Printf("打印所有身份：%v \n", result6)
//}

func TestDatabaseWithHighConcurrency(t *testing.T) {
	sqlite.OpenDefault()

	// 这个 WaitGroup 被用于等待该函数开启的所有协程。
	var wg sync.WaitGroup

	// 开启几个协程，并为其递增 WaitGroup 的计数器。
	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go operation(i, &wg)
	}

	// 阻塞，直到 WaitGroup 计数器恢复为 0，即所有协程的工作都已经完成。
	wg.Wait()
	result1, err := GetAllIdentityFromStorage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("打印所有身份：%v \n", string(result1[0].PrikeyRawByte))
	result1[0].UnLock("0123456789abcdef", sec.SM4ECB)
	fmt.Printf("打印所有身份：%v \n", string(result1[0].Prikey.GetBytes()))
}

// 每个协程都会运行该函数。
// 注意，WaitGroup 必须通过指针传递给函数。
func operation(i int, wg *sync.WaitGroup) {
	fmt.Printf("Operation %d starting\n", i)

	// 测试PersistIdentity
	pri, pub := sm2.GenKeyPair()
	id := identity.Identity{}
	id.Name = "xzw" + strconv.Itoa(i)
	id.KeyParam.PublicKeyAlgorithm = 0
	id.KeyParam.SignatureAlgorithm = 0
	id.Prikey = pri
	id.Pubkey = pub
	id.Passwd = "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
	id.Lock("0123456789abcdef", sec.SM4ECB)
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
	res, err := PersistIdentity(&id)
	if res != true {
		fmt.Println(err)
	}

	//测试GetIdentityByNameFromStorag
	result2, err := GetIdentityByNameFromStorage("xzw" + strconv.Itoa(i))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("插入身份：%v \n", result2)
	fmt.Printf("开始设置 xzw%v 为default身份\n", strconv.Itoa(i))

	//测试SetDefaultIdentityFromStorage
	result3, err := SetDefaultIdentityByNameInStorage("xzw" + strconv.Itoa(i))
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("xzw%v default设置结果：%v \n", strconv.Itoa(i), result3)

	//测试etDefaultIdentityFromStorage
	result4, err := GetDefaultIdentityFromStorage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("当前default身份：%v \n", result4)

	////测试DeleteIdentityByName
	//result5, err := DeleteIdentityByName("xzw" + strconv.Itoa(i))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("删除xzw%v结果：%v \n", strconv.Itoa(i), result5)

	// 通知 WaitGroup，当前协程的工作已经完成。
	wg.Done()
}
