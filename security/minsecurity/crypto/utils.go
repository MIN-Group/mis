/**
 * @Author: wzx
 * @Description:工具文件
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午6:18
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package crypto

import (
	"MIS-BC/security/minsecurity/crypto/sm2"
)
import sec "MIS-BC/security/minsecurity"

//定义pem格式到公私钥的转换
func FromPemToPubkey() {

}

func PubkeyToPem() {

}

func FromPemToPrikey() {

}

func PrikeyToPem() {

}

//通过给定的公钥结构体返回对应的编码结果和使用的算法
func MarshalPublicKey(pub interface{}) (sec.PublicKeyAlgorithm, []byte) {
	switch pub.(type) {
	case *sm2.Sm2PublicKey:
		p := pub.(*sm2.Sm2PublicKey)
		return sec.SM2, p.GetBytes()
	}

	return sec.SM2, nil
}

//通过给定的私钥结构体返回对应的编码和使用的算法
func MarshalPrivateKey(pri interface{}) (sec.PublicKeyAlgorithm, []byte) {
	switch pri.(type) {
	case *sm2.Sm2PrivateKey:
		p := pri.(*sm2.Sm2PrivateKey)
		return sec.SM2, p.GetBytes()
	}
	return sec.SM2, nil
}

//通过给定的私钥类型和编码返回对应的私钥结构体
func UnMarshalPrivateKey(priByte []byte, algo sec.PublicKeyAlgorithm) (sec.PrivateKey, error) {
	if priByte == nil || len(priByte) == 0 {
		return nil, nil
	}

	switch algo {
	case sec.SM2:
		p := sm2.Sm2PrivateKey{}
		err := p.SetBytes(priByte)
		if err != nil {
			return nil, err
		}
		var res sec.PrivateKey = &p
		return res, nil
	}
	return nil, nil
}

//通过给定的公钥类型和编码返回对应的公钥结构体

func UnMarshalPublicKey(pubByte []byte, algo sec.PublicKeyAlgorithm) (sec.PublicKey, error) {

	if pubByte == nil || len(pubByte) == 0 {
		return nil, nil
	}
	switch algo {
	case sec.SM2:
		p := sm2.Sm2PublicKey{}
		err := p.SetBytes(pubByte)
		if err != nil {
			return nil, err
		}
		var res sec.PublicKey = &p
		return res, nil
	}
	return nil, nil
}
