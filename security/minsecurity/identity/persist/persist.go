/**
 * @Author: wzx
 * @Description:身份存储接口
 * @Version: 1.0.0
 * @Date: 2021/1/16 上午6:09
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package persist

import (
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/identity"
	sqlite "MIS-BC/security/minsecurity/identity/persist/sqlite/db"
)


//默认是使用sqlite,修改存储引擎需要手动修改
var PersistPlugin = sec.Sqlite

//存储一个Identity,生成identity结构体和存储分开
func PersistIdentity(identity *identity.Identity) (bool,error) {
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.PersistIdentity(identity)
	}
	return false,nil
}

//通过Identity的名称，删除一个Identity
func DeleteIdentityByNameFromStorage(name string) (bool,error)  {
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.DeleteIdentityByName(name)
	}

	return false,nil
}

//通过Identity的名称，获得一个Identity对象
func GetIdentityByNameFromStorage(name string) (*identity.Identity, error) {
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.GetIdentityByNameFromStorage(name)
	}

	return nil, nil
}

//获得所有的身份实例
func GetAllIdentityFromStorage()  ([]*identity.Identity, error){
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.GetAllIdentityFromStorage()
	}

	return nil, nil
}

//通过Identity的名称，设置默认的Identity
func SetDefaultIdentityByNameInStorage(name string) (bool, error){
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.SetDefaultIdentityByNameInStorage(name)
	}

	return false, nil
}

//获得当前系统默认使用的Identity对象
func GetDefaultIdentityFromStorage() (*identity.Identity, error) {
	switch PersistPlugin {
	case sec.Sqlite:
		return sqlite.GetDefaultIdentityFromStorage()
	}

	return nil, nil
}