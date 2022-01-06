/**
 * @Author: wzx
 * @Description:通用基础定义
 * @Version: 1.0.0
 * @Date: 2021/1/15 下午11:29
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package minsecurity

//存储引擎,默认是使用sqlite,支持用户显示配置
type PersistPlugin int

const(
	Sqlite PersistPlugin = iota
)

//公钥算法
type PublicKeyAlgorithm int

const (
	SM2 PublicKeyAlgorithm = iota
)

//签名算法
type SignatureAlgorithm int

const (
	SM2WithSM3  SignatureAlgorithm = iota
)

//对称加密算法
type SymmetricAlgorithm int

const(
	SM4ECB	SymmetricAlgorithm = iota
	SM4CBC
)

//证书用途
type KeyUsage int

const (
	ContentCommitment KeyUsage = iota
	DataEncipherment
	CertSign
)




