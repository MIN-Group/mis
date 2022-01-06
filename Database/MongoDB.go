/**
 * @Author: xzw
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/6/8 晚上8:00
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package MongoDB

import (
	"MIS-BC/MetaData"
	"MIS-BC/common"
	"MIS-BC/security/code"
	"MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"MIS-BC/security/minsecurity/identity/persist/sqlite"
	"MIS-BC/security/minsecurity/identity/persist/sqlite/db"
	"fmt"
	"github.com/JodeZer/mgop"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"log"
	"strconv"
	"time"
)

type Mongo struct {
	Pubkey      string
	Height      int
	Block       MetaData.BlockGroup
	pool        mgop.SessionPool
	CacheNumber int              // 缓存的状态信息条数
}

var table_create string = `CREATE TABLE IF NOT EXISTS "identityinfo" (
	"name"	varchar(255) NOT NULL,
	"pubkey"	varchar(255),
	"prikey"	varchar(255),
	"pubkey_algo"	int,
	"signature_algo"int,
	"pass"	varchar(255),
	"cert"	TEXT,
	"is_default"	int DEFAULT 0,
	"prikey_raw_byte" varchar(255),
	PRIMARY KEY("name")
);`

func (pl *Mongo) SetConfig(config common.Config) {
	var err error
	pl.Pubkey = config.MyPubkey
	pl.Height = -1
	pl.pool, err = mgop.DialStrongPool("mongodb://mis:mis20201001@localhost/blockchain", 5)
	if err != nil {
		common.Logger.Error("err !!%s", err)
		return
	}
	if config.DropDatabase {
		pl.deleteDB()
	}
	pl.CacheNumber = config.CacheTime * (60 / int(config.GenerateBlockPeriod))
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "bcstatus"
	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	count, err := c.Find(nil).Count()
	if count >= pl.CacheNumber {
		var bs MetaData.BCStatus
		c.Find(nil).Sort("timestamp").Skip(pl.CacheNumber).Limit(1).One(&bs)
		c.RemoveAll(bson.M{"timestamp": bson.M{"$lte": bs.Timestamp}})
	}

	// index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	// Open("/home/min/identity/"+index+"/", dbFile)
	sqlite.Open(config.SqlitePath + config.HostName + "/")

	if pl.HasIdentityData("identityidentifier", "root") {
		item := pl.GetOneIdentityMapFromDatabase("root")
		pl.ResetIdentityTimeStamp(item)

		id, err := db.GetIdentityByNameFromStorage("/root")
		if err != nil {
			fmt.Println(err)
		}
		if id == nil {
			identity := pl.GetOneIdentityFromDatabase("identityidentifier", "root")
			pri := sm2.Sm2PrivateKey{}
			pri.SetBytes([]byte(config.MyPrikey))
			var prikey minsecurity.PrivateKey = &pri

			i := identity.ParseBCIdentityToCommon()
			i.Prikey = prikey
			// PersistIdentity(&i, "/home/min/identity/"+index+"/", dbFile)
			if _, err := db.PersistIdentity(&i); err != nil {
				common.Logger.Error(err)
			}
			// SetDefaultIdentityByNameInStorage("/root", "/home/min/identity/"+index+"/", dbFile)
			_, err = db.SetDefaultIdentityByNameInStorage("/root")
			if err != nil {
				common.Logger.Error(err)
			}
		}
	} else {
		// 填充证书内容
		pub := sm2.Sm2PublicKey{}
		pub.SetBytes([]byte(config.MyPubkey))
		var pubkey minsecurity.PublicKey = &pub
		cert := cert.Certificate{}
		cert.Version = 0
		cert.SerialNumber = 1
		cert.PublicKey = pubkey
		cert.SignatureAlgorithm = 0
		cert.PublicKeyAlgorithm = 0
		cert.IssueTo = "/root"
		cert.Issuer = "/root"
		cert.NotBefore = time.Now().Unix()
		cert.NotAfter = time.Now().AddDate(5, 0, 0).Unix()
		cert.KeyUsage = minsecurity.CertSign
		cert.IsCA = true
		cert.Timestamp = time.Now().Unix()

		pri := sm2.Sm2PrivateKey{}
		pri.SetBytes([]byte(config.MyPrikey))
		var prikey minsecurity.PrivateKey = &pri
		err := cert.SignCert(prikey)
		if err != nil {
			common.Logger.Error(err)
		}

		c, err := cert.ToPem([]byte("Pkusz112233"), 0)
		if err != nil {
			fmt.Printf("err !!%s", err)
			return
		}

		var t []MetaData.ModifyRecord
		root := MetaData.Identity{IsValid: code.VALID, IdentityIdentifier: "root", Passwd: "Pkusz112233", Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Pubkey: config.MyPrikey, Cert: c, Type: "manager-act", Command: "CreatManager", KeyParam: MetaData.KeyParam{0, 0},
			ModifyRecords: append(t, MetaData.ModifyRecord{"manager-act", "CreatManager", time.Now().Format("2006-01-02 15:04:05")})}
		common.Logger.Info("管理员账号：", root)

		croot := root.ParseBCIdentityToCommon()
		croot.Prikey = prikey
		// PersistIdentity(&croot, "/home/min/identity/"+index+"/", dbFile)
		if _, err := db.PersistIdentity(&croot); err != nil {
			common.Logger.Error(err)
		}
		// SetDefaultIdentityByNameInStorage("/root", "/home/min/identity/"+index+"/", dbFile)
		_, err = db.SetDefaultIdentityByNameInStorage("/root")
		if err != nil {
			common.Logger.Error(err)
		}

		pl.SaveIdentityToDatabase(root)
	}
	identities := pl.GetAllIdentityFromDatabase()
	// fmt.Println("ids:", identities)
	for i := 0; i < len(identities); i++ {
		if identities[i].IdentityIdentifier == "root" {
			continue
		}
		id := identities[i].ParseBCIdentityToCommon()
		identity, err := db.GetIdentityByNameFromStorage(id.Name)
		// fmt.Println("id:", id, "identity:", identity)
		if err != nil {
			common.Logger.Error(err)
		} else if identity == nil {
			if _, err := db.PersistIdentity(&id); err != nil {
				common.Logger.Error(err)
			}
		} else {
			continue
		}
	}
}

func (pl *Mongo) InsertToMogoBG(bg MetaData.BlockGroup, index string) {
	session := pl.pool.AcquireSession()
	defer session.Release()

	index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index_mongo + "-block" + "-" + strconv.Itoa(bg.Height/(10*10000)))
	for _, v := range bg.Blocks {
		if bg.Height == v.Height {
			err := c.Insert(&v)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	bg.Blocks = bg.Blocks[0:0] //清空Blocks

	c1 := session.DB("blockchain").C(index_mongo + "-blockgroup" + "-" + strconv.Itoa(bg.Height/(10*10000)))
	err := c1.Insert(&bg)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoBlockstate(item MetaData.BlockState, index string) {
	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoAccount(item MetaData.Account, index string) {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoIdentity(item MetaData.Identity, index string) {
	session := pl.pool.AcquireSession()
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoBCStatus(item MetaData.BCStatus, index string) {
	session := pl.pool.AcquireSession()
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoNodeIdentityTrans(item MetaData.IdentityTransformation, index string) {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoNodelist(item MetaData.NodeList, index string) {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) InsertToMogoIdentityTransList(item MetaData.IdentityTransList, index string) {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		common.Logger.Error(err)
	}
}

func (pl *Mongo) GetHeight() int {
	return pl.Height
}

func (pl *Mongo) deleteDB() {
	session := pl.pool.AcquireSession()
	defer session.Release()

	tmp := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	index1 := tmp + "-" + "Identity"
	index2 := tmp + "-" + "NDN"
	index3 := tmp + "-" + "UserLog"
	index4 := tmp + "-" + "white-paper"
	index5 := tmp + "-" + "machine"
	index6 := tmp + "-" + "domain"
	index7 := tmp + "-" + "management"
	index8 := tmp + "-" + "identity"
	index9 := tmp + "-" + "account"
	index10 := tmp + "-" + "nodelist"
	index11 := tmp + "-" + "identityTrans"
	index12 := tmp + "-" + "blockstate"
	index13 := tmp + "-" + "contract_account"
	index14 := tmp + "-" + "account-data"

	height := pl.GetHeight()
	for i := 0; i <= height/(10*10000); i++ {
		_ = session.DB("blockchain").C(tmp + "-block" + "-" + strconv.Itoa(i)).DropCollection()
		_ = session.DB("blockchain").C(tmp + "-blockgroup" + "-" + strconv.Itoa(i)).DropCollection()
	}

	_ = session.DB("blockchain").C(index1).DropCollection()
	_ = session.DB("blockchain").C(index2).DropCollection()
	_ = session.DB("blockchain").C(index3).DropCollection()
	_ = session.DB("blockchain").C(index4).DropCollection()
	_ = session.DB("blockchain").C(index5).DropCollection()
	_ = session.DB("blockchain").C(index6).DropCollection()
	_ = session.DB("blockchain").C(index7).DropCollection()
	_ = session.DB("blockchain").C(index8).DropCollection()
	_ = session.DB("blockchain").C(index9).DropCollection()
	_ = session.DB("blockchain").C(index10).DropCollection()
	_ = session.DB("blockchain").C(index11).DropCollection()
	_ = session.DB("blockchain").C(index12).DropCollection()
	_ = session.DB("blockchain").C(index13).DropCollection()
	_ = session.DB("blockchain").C(index14).DropCollection()
	_ = session.DB("blockchain").C(tmp + "-block").DropCollection()
	_ = session.DB("blockchain").C(tmp + "-blockgroup").DropCollection()
	_ = session.DB("blockchain").C(tmp).DropCollection()
}

func (pl *Mongo) HasData(typ, key, value string) bool {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	if len(items) > 0 {
		return true
	}
	return false
}

func (pl *Mongo) SaveDataToDatabase(typ string, item MetaData.Identity) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	pl.InsertToMogoIdentity(item, subname)
}

func (pl *Mongo) GetOneDataFromDatabase(typ, key, value string) MetaData.Identity {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var item MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value}).One(&item)
	if err != nil {
		common.Logger.Error(err)
	}
	return item
}

func (pl *Mongo) HasDataByTwoKey(typ, key1, value1, key2, value2 string) bool {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key1: value1, key2: value2}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	if len(items) > 0 {
		return true
	}
	return false
}

func (pl *Mongo) GetMultipleDataFromDatabase(typ, key, value string) []MetaData.Identity {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)

	err := c.Find(bson.M{key: value}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return items
}

func (pl *Mongo) GetRegexMultipleDataFromDatabase(typ, key, pattern, option string) []MetaData.Identity {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)

	err := c.Find(bson.M{key: bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: option}}}).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return items
}

func (pl *Mongo) GetAllDataFromDatabase(typ string) []MetaData.Identity {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return items
}

func (pl *Mongo) GetCountFromDatabase(typ string) int {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var total int
	c := session.DB("blockchain").C(subname)
	total, err := c.Find(nil).Count()
	if err != nil {
		common.Logger.Error(err)
	}
	return total
}

func (pl *Mongo) GetPageDataFromDatabase(typ string, skip, limit int) []MetaData.Identity {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	var items []MetaData.Identity
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).Sort("Username").Skip(skip).Limit(limit).All(&items)
	if err != nil {
		common.Logger.Error(err)
	}
	return items
}

func (pl *Mongo) DeleteData(typ, key, value string) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	_, err := c.RemoveAll(bson.M{key: value})
	if err != nil {
		common.Logger.Error(err)
	}
}
