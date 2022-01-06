// Package security
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/2/1 9:29 上午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package security

import (
	"MIS-BC/security/minsecurity"
	cert2 "MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"MIS-BC/security/minsecurity/identity"
	"MIS-BC/security/minsecurity/identity/persist"
	"MIS-BC/security/minsecurity/identity/persist/sqlite"
	"fmt"
	"time"
)

const DefaultIdentityDBPath = "/usr/local/.mir/identity/"

type IdentityManager struct {
	// 一个map，存储了身份名字和网络身份实体的映射
	identities map[string]*identity.Identity

	// 一个map，存储了身份名字和对应身份的版本号
	// 1. 初始加载到内存中时，所有身份的版本号均为0
	// 2. 接着每次对该网络身份进行了任何的修改，对应版本号都会++
	versionMap map[string]uint64

	// 默认网络身份
	defaultIdentity *identity.Identity

	// 对秘钥加密所使用的加密算法
	privateKeyEncryptionAlgorithm minsecurity.SymmetricAlgorithm

	// 版本号，每次创建一个对象，version从0开始，对身份的任何增删改都会导致版本号的增加
	version uint64
}

// Init 默认初始化，sqlite数据库文件存储到默认路径
//
// @Description:
// @receiver i
// @return error
//
func (i *IdentityManager) Init() error {
	return i.InitByPath(DefaultIdentityDBPath)
}

// InitByPath
// 初始化
//
// @Description:
//	1. 尝试从本地持久化存储中拉取所有的网络身份信息；
//	2. 尝试获取本地持久化记录当中本系统的默认网络身份信息；
// @receiver i
//
func (i *IdentityManager) InitByPath(path string) error {
	sqlite.Open(path)
	// 首先尝试从本地持久化存储中拉取所有的网络身份信息
	if err := i.loadAllIdentifies(); err != nil {
		return err
	}

	// 获取默认的网络身份
	defaultIdentity, err := persist.GetDefaultIdentityFromStorage()
	if err != nil {
		return err
	}
	i.defaultIdentity = defaultIdentity

	// 默认使用 SM4_ECB 对秘钥进行加解密
	i.privateKeyEncryptionAlgorithm = minsecurity.SM4ECB
	return nil
}

//
// 从持久化载体中获取所有存储的网络身份信息，放到 identities 当中
//
// @Description:
// @receiver
//
func (i *IdentityManager) loadAllIdentifies() error {
	// 从持久化存储中获取所有的网络身份
	localIdentifies, err := persist.GetAllIdentityFromStorage()
	if err != nil {
		return err
	}

	// 将网络身份存储到Map当中
	i.identities = make(map[string]*identity.Identity)
	i.versionMap = make(map[string]uint64)

	for _, item := range localIdentifies {
		i.identities[item.Name] = item
		i.versionMap[item.Name] = 0
	}
	return nil
}

// GetAllIdentities 获取所有的网络身份
//
// @Description:
// @receiver i
// @return map[string]*identity.Identity
//
func (i *IdentityManager) GetAllIdentities() map[string]*identity.Identity {
	return i.identities
}

// GetIdentityByName
// 通过身份的名字获取网络身份
//
// @Description:
// @receiver i
// @param name
//
func (i *IdentityManager) GetIdentityByName(name string) *identity.Identity {
	return i.identities[name]
}

// DeleteIdentityByName
// 通过网络身份的名字和对应的密码删除一个网络身份
//
// @Description:
//	1. 如果这个网络身份只包含证书、公钥，并不包含私钥
// @receiver i
// @param name
// @param password
//
func (i *IdentityManager) DeleteIdentityByName(name string, password string) (bool, error) {
	// 首先在Map里面删除它
	delete(i.identities, name)
	delete(i.versionMap, name)
	i.version++
	return persist.DeleteIdentityByNameFromStorage(name)
}

// SaveIdentity
// 将一个网络身份保存到本地
//
// @Description:
//	1. 首先检查网络身份在本地是否存在，如果存在，且 force = false，则不覆盖，保存失败
//	2. 如果网络身份已存在，且 force = true，则覆盖原先的网络身份
//	3. 如果网络身份不存在，则直接保存
// @receiver i
// @param newIdentity
// @param force
// @return error
//
func (i *IdentityManager) SaveIdentity(newIdentity *identity.Identity, force bool) error {
	if i.ExistIdentity(newIdentity.Name) && !force {
		// 身份已存在，且不强制覆盖，则返回错误，保存失败
		return IdentityManagerError{msg: fmt.Sprintf(
			"Save Identity failed, identify with name %s already exist!", newIdentity.Name)}
	}

	// 身份不存在，或者存在但是指定强制覆盖，则保存该网络身份
	if force {
		// 如果是强制保存，则先删除之前存储的
		if _, err := persist.DeleteIdentityByNameFromStorage(newIdentity.Name); err != nil {
			return err
		}
	}

	// 将新的网络身份进行持久化存储
	if _, err := persist.PersistIdentity(newIdentity); err != nil {
		return err
	}

	i.identities[newIdentity.Name] = newIdentity
	currentVersion, ok := i.versionMap[newIdentity.Name]
	if !ok {
		currentVersion = 0
	} else {
		currentVersion++
	}
	i.versionMap[newIdentity.Name] = currentVersion
	return nil
}

// CreateIdentityByName
// 指定一个身份名字在本地创建一个网络身份
//
// @Description:
// @receiver i
// @param name
// @param password
// @return *identity.Identity
// @return error
//
func (i *IdentityManager) CreateIdentityByName(name string, password string) (*identity.Identity, error) {
	// 默认使用国密算法进行加解密和签名
	return i.CreateIdentityByNameAndKeyParam(name, identity.KeyParam{
		PublicKeyAlgorithm: minsecurity.SM2,
		SignatureAlgorithm: minsecurity.SM2WithSM3,
	}, password)
}

// CreateIdentityByNameAndKeyParam
// 指定一个身份名字和KeyParam，在本地创建一个网络身份
//
// @Description:
//	1. 首先检查指定的身份名字在本地存储中是否已经存在，如果已经存在则创建失败；
//	2. 如果名字不冲突，则检查 password 是否为空字符串：
//		a. 如果 password 为空字符串，则不对身份的秘钥进行加密保护；
//		b. 如果 password 为非空字符串，则认为需要用其对秘钥进行加密保护，调用 Identity.Lock(password)
// @receiver i
// @param name
// @param param
// @param password
// @return *identity.Identity
// @return error
//
func (i *IdentityManager) CreateIdentityByNameAndKeyParam(name string, param identity.KeyParam, password string) (*identity.Identity, error) {
	// 第一步，判断在本地是否有同名的网络身份
	if i.identities[name] != nil {
		return nil, IdentityManagerError{msg: fmt.Sprintf(
			"Identity name => %s, already exists!", name)}
	}

	// 本地生成一堆公私钥
	pri, pub := sm2.GenKeyPair()
	// 实例化一个新的网络身份
	newIdentity := identity.Identity{
		Name:     name,
		Prikey:   pri,
		Pubkey:   pub,
		KeyParam: param,
	}

	// 如果 password 为非空字符串，则认为其需要对秘钥进行加密保护，使用传入的密码对秘钥进行锁定
	if password != "" {
		_, err := newIdentity.Lock(password, minsecurity.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	// 将新生成的网络身份进行持久化存储
	_, err := persist.PersistIdentity(&newIdentity)
	if err != nil {
		return nil, err
	}

	// 持久化存储成功则在内存map中也存储一份
	i.identities[name] = &newIdentity
	i.versionMap[name] = 0

	i.version++
	return &newIdentity, nil
}

// ResetIdentityPasswd 修改网络身份的密码
//
// @Description:
// @param name
// @param oldPasswd
// @param newPasswd
// @return error
//
func (i *IdentityManager) ResetIdentityPasswd(name, oldPasswd, newPasswd string) error {
	if myIdentity := i.GetIdentityByName(name); myIdentity == nil {
		return IdentityManagerError{
			msg: "Idnetity not exists",
		}
	} else {
		if oldPasswd != "" && myIdentity.IsLocked() {
			if ok, err := myIdentity.UnLock(oldPasswd, i.privateKeyEncryptionAlgorithm); err != nil {
				return err
			} else if !ok {
				return KeyChainError{
					msg: "Unlock " + myIdentity.Name + " by " + oldPasswd + " failed!!",
				}
			}
		}

		// 使用新的密码加密
		if _, err := myIdentity.Lock(newPasswd, i.privateKeyEncryptionAlgorithm); err != nil {
			return err
		}

		// 强制保存
		if err := i.SaveIdentity(myIdentity, true); err != nil {
			return err
		}
	}
	return nil
}

// SetDefaultIdentity
// 设置当前系统的默认网络身份
//
// @Description:
// @receiver k
// @param identity
//
func (i *IdentityManager) SetDefaultIdentity(identity *identity.Identity) (bool, error) {
	if identity == nil {
		return false, KeyChainError{msg: fmt.Sprintf(
			"Not allow set nil as default Identity!")}
	}
	i.defaultIdentity = identity
	i.version++
	return persist.SetDefaultIdentityByNameInStorage(identity.Name)
}

// ExistIdentity
// 判断某个网络身份是否存在
//
// @Description:
// @receiver i
// @param name
// @return bool
//
func (i *IdentityManager) ExistIdentity(name string) bool {
	return i.identities[name] != nil
}

// GetDefaultIdentity 获取默认网络身份
//
// @Description:
// @receiver i
// @return *identity.Identity
//
func (i *IdentityManager) GetDefaultIdentity() *identity.Identity {
	return i.defaultIdentity
}

// GetCurrentVersion 获取当前的版本号
//
// @Description:
// @receiver i
// @return uint64
//
func (i *IdentityManager) GetCurrentVersion() uint64 {
	return i.version
}

// GetIdentityVersion 获取某个网络身份的版本号
//
// @Description:
// @receiver i
// @param identityName
// @return uint64
//
func (i *IdentityManager) GetIdentityVersion(identityName string) uint64 {
	version, ok := i.versionMap[identityName]
	if !ok {
		return 0
	} else {
		return version
	}
}

// LoadCert 加载证书
//
// @Description:
//	1. 如果网络身份不存在，则创建一个新的网络身份，并设置名字和公钥
//	2. 如果网络身份存在，则更新 cert 字段
// @receiver i
// @param identityName
// @param cert
//
func (i *IdentityManager) LoadCert(identityName string, cert *cert2.Certificate) error {
	if i.ExistIdentity(identityName) {
		// 已经存在
		id, err := persist.GetIdentityByNameFromStorage(identityName)
		if err != nil {
			return err
		}
		id.Cert = *cert

		// 持久化保存
		if err := i.SaveIdentity(id, true); err != nil {
			return err
		}

		// 更新内存中的对应身份
		i.GetIdentityByName(identityName).Cert = *cert
	} else {
		// 不存在
		id := identity.Identity{}
		id.Name = cert.IssueTo
		id.Cert = *cert
		id.KeyParam.PublicKeyAlgorithm = cert.PublicKeyAlgorithm
		id.KeyParam.SignatureAlgorithm = cert.SignatureAlgorithm
		id.Pubkey = cert.PublicKey
		if err := i.SaveIdentity(&id, true); err != nil {
			return err
		}

		// 将新创建的网络身份保存到内存当中
		i.identities[id.Name] = &id
	}

	return nil
}

// SelfIssue 使用指定的网络身份给自己签发一个自签证书
//
// @Description:
// @receiver i
// @param identityName
// @return error
//
func (i *IdentityManager) SelfIssue(identityName string, passwd string) error {
	// 判断身份是否存在
	if !i.ExistIdentity(identityName) {
		return IdentityManagerError{
			msg: "Target identity not exists!",
		}
	}

	// 从存储中获取对象，不影响内存中的身份
	id, err := persist.GetIdentityByNameFromStorage(identityName)
	if err != nil {
		return err
	}

	// 如果需要，则解锁身份
	if id.IsLocked() {
		if _, err := id.UnLock(passwd, minsecurity.SM4ECB); err != nil {
			return err
		}
	}

	// 填充证书内容
	cert := cert2.Certificate{}
	cert.Version = 0
	cert.SerialNumber = 1
	cert.PublicKey = id.Pubkey
	cert.SignatureAlgorithm = id.KeyParam.SignatureAlgorithm
	cert.PublicKeyAlgorithm = id.KeyParam.PublicKeyAlgorithm
	cert.IssueTo = id.Name
	cert.Issuer = id.Name
	cert.NotBefore = time.Now().Unix()
	cert.NotAfter = time.Now().AddDate(1, 0, 0).Unix()
	cert.KeyUsage = minsecurity.CertSign
	cert.IsCA = true
	cert.Timestamp = time.Now().Unix()

	// 签发证书
	if err := cert.SignCert(id.Prikey); err != nil {
		return err
	}

	// 保存证书
	if err := i.LoadCert(id.Name, &cert); err != nil {
		return err
	}

	return nil
}

// DumpCert 导出指定网络身份的证书
//
// @Description:
// @receiver i
// @param identityName
// @return string
// @return error
//
func (i *IdentityManager) DumpCert(identityName string) (string, error) {
	// 首先判断指定的网络身份是否存在
	targetIdentity := i.GetIdentityByName(identityName)
	if targetIdentity == nil {
		return "", IdentityManagerError{
			msg: "Target identity not exists!",
		}
	}

	// 判断证书是否存在
	if targetIdentity.Cert.Issuer == "" && targetIdentity.Cert.Signature == nil {
		return "", IdentityManagerError{
			msg: "Target identity not exists!",
		}
	}

	// 导出证书
	if str, err := (&targetIdentity.Cert).ToPem([]byte(""), minsecurity.SM4ECB); err != nil {
		return "", err
	} else {
		return str, nil
	}
}

// ImportCert 导入一个证书
//
// @Description:
// @receiver i
// @param certStr
// @return error
//
func (i *IdentityManager) ImportCert(certBytes []byte) error {
	// 解析证书
	cert := cert2.Certificate{}
	if err := cert.FromPem(string(certBytes), nil, minsecurity.SM4ECB); err != nil {
		return err
	}

	// 加载证书
	if err := i.LoadCert(cert.IssueTo, &cert); err != nil {
		return err
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
///// 错误处理
/////////////////////////////////////////////////////////////////////////////////////////////////////////

type IdentityManagerError struct {
	msg string
}

func (i IdentityManagerError) Error() string {
	return fmt.Sprintf("IdentityManagerError: %s", i.msg)
}
