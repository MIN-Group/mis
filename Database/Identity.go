/**
 * @Author: xzw
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/6/10 晚上7:00
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package MongoDB

import (
	"MIS-BC/MetaData"
	"MIS-BC/common"
	"MIS-BC/security/code"
	"MIS-BC/utils"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"log"
	"strconv"
	"time"
)

func (pl *Mongo) HasIdentityData(key, value string) bool {
	value = utils.EncryptString(value)
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value}).All(&items)
	if err != nil {
		log.Println(err)
	}
	if len(items) > 0 {
		return true
	}
	return false
}

func (pl *Mongo) GetAllIdentityFromDatabase() []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).All(&items)
	if err != nil {
		log.Println(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetPageIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		log.Println(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetPendingIdentityFromDatabase() []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": code.PENDING_REVIEW}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetPagePendingIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": code.PENDING_REVIEW}).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetPendingIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(bson.M{"isvalid": code.PENDING_REVIEW}).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) GetCheckedIdentityFromDatabase() []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.VALID, code.INVALID}}}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetPageCheckedIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.VALID, code.INVALID}}}).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetCheckedIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.VALID, code.INVALID}}}).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) GetPageDisabledIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.INVALID}}}).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetDisabledIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.INVALID}}}).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) GetPageAbledIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.VALID}}}).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetAbledIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.VALID}}}).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) GetPageWithoutCertIdentityFromDatabase(skip, limit int) []MetaData.Identity {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.WITHOUT_CERT}}}).Sort("IdentityIdentifier").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return decryptAllIdentity(items)
}

func (pl *Mongo) GetWithoutCertIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(bson.M{"isvalid": bson.M{"$in": []int{code.WITHOUT_CERT}}}).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) SaveIdentityToDatabase(item MetaData.Identity) {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	pl.InsertToMogoIdentity(encryptIdentity(item), subname)
}

func (pl *Mongo) DeleteIdentity(key, value string) {
	value = utils.EncryptString(value)
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()
	c := session.DB("blockchain").C(subname)
	_, err := c.RemoveAll(bson.M{key: value})
	if err != nil {
		log.Println(err)
	}
}

func (pl *Mongo) GetOneIdentityFromDatabase(key, value string) MetaData.Identity {
	value = utils.EncryptString(value)
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var item MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value}).One(&item)
	if err != nil {
		log.Println(err)
	}
	return decryptIdentity(item)
}

func (pl *Mongo) GetIdentityCountFromDatabase() int {
	typ := "Identity"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}
	return total
}

func (pl *Mongo) GetOneIdentityMapFromDatabase(IdentityIdentifier string) map[string]interface{} {
	IdentityIdentifier = utils.EncryptString(IdentityIdentifier)
	fmt.Println("身份标识加密后的数据为", IdentityIdentifier)
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	item := make(map[string]interface{})
	c := session.DB("blockchain").C(subname)
	fmt.Println("所在的表名为：", subname)
	err := c.Find(bson.M{"identityidentifier": IdentityIdentifier}).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}

func (pl *Mongo) EnableIdentity(isvalid int, cert string, item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	selector := bson.M{"_id": item["_id"]}

	v := bson.M{"$set": bson.M{"isvalid": isvalid}}
	err := c.Update(selector, v)
	if err != nil {
		common.Logger.Error("EnableIdentity函数启用身份失败：", err)
	}

	ct := bson.M{"$set": bson.M{"cert": utils.EncryptString(cert)}}
	err = c.Update(selector, ct)
	if err != nil {
		common.Logger.Error("EnableIdentity函数启用身份失败：", err)
	}
}

func (pl *Mongo) ResetIdentityValidation(isvalid int, item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	selector := bson.M{"_id": item["_id"]}
	data := bson.M{"$set": bson.M{"isvalid": isvalid}}

	c := session.DB("blockchain").C(subname)
	err := c.Update(selector, data)
	if err != nil {
		common.Logger.Error("ResetIdentityValidation函数更新身份有效性失败：", err)
	}
}

func (pl *Mongo) CertRevocation(item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	selector := bson.M{"_id": item["_id"]}

	c := session.DB("blockchain").C(subname)

	v := bson.M{"$set": bson.M{"isvalid": code.WITHOUT_CERT}}
	err := c.Update(selector, v)
	if err != nil {
		common.Logger.Error("CertRevocation函数重置身份有效性失败：", err)
	}

	ct := bson.M{"$set": bson.M{"cert": utils.EncryptString("")}}
	err = c.Update(selector, ct)
	if err != nil {
		common.Logger.Error("CertRevocation函数重新颁发证书失败：", err)
	}
}

func (pl *Mongo) CertReissue(cert string, item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	selector := bson.M{"_id": item["_id"]}

	v := bson.M{"$set": bson.M{"isvalid": code.VALID}}
	err := c.Update(selector, v)
	if err != nil {
		common.Logger.Error("CertReissue函数重新颁发证书失败：", err)
	}

	ct := bson.M{"$set": bson.M{"cert": utils.EncryptString(cert)}}
	err = c.Update(selector, ct)
	if err != nil {
		common.Logger.Error("CertReissue函数重新颁发证书失败：", err)
	}
}

func (pl *Mongo) ResetIdentityPassword(password string, item map[string]interface{}) {
	password = utils.EncryptString(password)
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	selector := bson.M{"_id": item["_id"]}
	data := bson.M{"$set": bson.M{"passwd": password}}

	c := session.DB("blockchain").C(subname)
	err := c.Update(selector, data)
	if err != nil {
		common.Logger.Error("ResetMINUserPassword函数修改密码失败：", err)
	}
}

func (pl *Mongo) ResetIdentityTimeStamp(item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	selector := bson.M{"_id": item["_id"]}
	t := utils.EncryptString(time.Now().Format("2006-01-02 15:04:05"))
	data := bson.M{"$set": bson.M{"timestamp": t}}

	c := session.DB("blockchain").C(subname)
	err := c.Update(selector, data)
	if err != nil {
		common.Logger.Error("ResetIdentityTimeStamp函数更新时间戳失败：", err)
	}
}

func (pl *Mongo) ResetIdentityIPIdentifier(ip string, item map[string]interface{}) {
	ip = utils.EncryptString(ip)
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	selector := bson.M{"_id": item["_id"]}
	data := bson.M{"$set": bson.M{"ipidentifier": ip}}

	c := session.DB("blockchain").C(subname)
	err := c.Update(selector, data)
	if err != nil {
		common.Logger.Error("ResetIdentityIPIdentifier函数更新IP失败：", err)
	}
}

func (pl *Mongo) UpdateIdentityModifyRecords(res map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	Identity := pl.GetOneIdentityFromDatabase("identityidentifier", res["IdentityIdentifier"].(string))
	Identity.ModifyRecords = append(Identity.ModifyRecords, MetaData.ModifyRecord{Type: res["Type"].(string),
		Command: res["Command"].(string), Timestamp: time.Now().Format("2006-01-02 15:04:05")})
	t := Identity.ModifyRecords

	item := pl.GetOneIdentityMapFromDatabase(res["IdentityIdentifier"].(string))
	selector := bson.M{"_id": item["_id"]}

	c := session.DB("blockchain").C(subname)

	modifyrecords := bson.M{"$set": bson.M{"modifyrecords": t}}
	err := c.Update(selector, modifyrecords)
	if err != nil {
		common.Logger.Error("UpdateIdentityModifyRecords函数更新身份操作记录失败：", err)
	}

	command := bson.M{"$set": bson.M{"command": utils.EncryptString(res["Command"].(string))}}
	err = c.Update(selector, command)
	if err != nil {
		common.Logger.Error("UpdateIdentityModifyRecords函数更新身份操作记录失败：", err)
	}

	typ := bson.M{"$set": bson.M{"type": utils.EncryptString(res["Type"].(string))}}
	err = c.Update(selector, typ)
	if err != nil {
		common.Logger.Error("UpdateIdentityModifyRecords函数更新身份操作记录失败：", err)
	}
}

func (pl *Mongo) UpdateIdentityModifyRecordsforAuthentication(id MetaData.Identity) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "Identity"
	session := pl.pool.AcquireSession()
	defer session.Release()

	Identity := pl.GetOneIdentityFromDatabase("identityidentifier", id.IdentityIdentifier)
	Identity.ModifyRecords = append(Identity.ModifyRecords, MetaData.ModifyRecord{Type: id.Type,
		Command: id.Command, Timestamp: time.Now().Format("2006-01-02 15:04:05")})
	t := Identity.ModifyRecords

	item := pl.GetOneIdentityMapFromDatabase(id.IdentityIdentifier)
	selector := bson.M{"_id": item["_id"]}

	c := session.DB("blockchain").C(subname)

	modifyrecords := bson.M{"$set": bson.M{"modifyrecords": t}}
	err := c.Update(selector, modifyrecords)
	if err != nil {
		fmt.Println("UpdateIdentityModifyRecords函数更新身份操作记录失败")
		// log.Fatal(err)
	}

	command := bson.M{"$set": bson.M{"command": utils.EncryptString(id.Command)}}
	err = c.Update(selector, command)
	if err != nil {
		fmt.Println("UpdateIdentityModifyRecords函数更新身份操作记录失败")
		// log.Fatal(err)
	}

	typ := bson.M{"$set": bson.M{"type": utils.EncryptString(id.Type)}}
	err = c.Update(selector, typ)
	if err != nil {
		fmt.Println("UpdateIdentityModifyRecords函数更新身份操作记录失败")
		// log.Fatal(err)
	}
}


func decryptAllIdentity(sources []MetaData.Identity) []MetaData.Identity {
	var results []MetaData.Identity

	for _, min := range sources {
		results = append(results, decryptIdentity(min))
	}

	return results
}

func decryptIdentity(source MetaData.Identity) MetaData.Identity {
	result := MetaData.Identity{}
	result.Type = utils.DecryptString(source.Type)
	result.Command = utils.DecryptString(source.Command)
	result.Pubkey = utils.DecryptString(source.Pubkey)

	if source.Cert == "" {
		result.Cert = ""
	} else {
		result.Cert = utils.DecryptString(source.Cert)
	}

	result.Timestamp = utils.DecryptString(source.Timestamp)
	result.KeyParam = source.KeyParam
	result.IPIdentifier = utils.DecryptString(source.IPIdentifier)
	result.Passwd = utils.DecryptString(source.Passwd)
	result.ModifyRecords = source.ModifyRecords
	result.IsValid = source.IsValid

	result.IdentityIdentifier = utils.DecryptString(source.IdentityIdentifier)

	return result
}

func encryptIdentity(source MetaData.Identity) MetaData.Identity {
	result := MetaData.Identity{}
	result.Type = utils.EncryptString(source.Type)
	result.Command = utils.EncryptString(source.Command)
	result.Pubkey = utils.EncryptString(source.Pubkey)

	if source.Cert == "" {
		result.Cert = ""
	} else {
		result.Cert = utils.EncryptString(source.Cert)
	}
	result.Timestamp = utils.EncryptString(source.Timestamp)
	result.KeyParam = source.KeyParam
	result.IPIdentifier = utils.EncryptString(source.IPIdentifier)
	result.Passwd = utils.EncryptString(source.Passwd)
	result.ModifyRecords = source.ModifyRecords
	result.IsValid = source.IsValid

	result.IdentityIdentifier = utils.EncryptString(source.IdentityIdentifier)
	fmt.Println("注册的用户身份标识为：", result.IdentityIdentifier)

	return result
}
