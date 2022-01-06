package MongoDB

import (
	"MIS-BC/MetaData"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"log"
	"strconv"
)

func (pl *Mongo) GetAccountFromDatabase() MetaData.Account {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "account"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.Account
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}

func (pl *Mongo) InsertOrUpdateAccount(item MetaData.Account) {
	pl.dropAccount()
	pl.saveAccountToDatabase(item)
}

func (pl *Mongo) dropAccount() {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	tmp := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	index9 := tmp + "-" + "account"
	_ = session.DB("blockchain").C(index9).DropCollection()
}

func (pl *Mongo) saveAccountToDatabase(item MetaData.Account) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "account"
	pl.InsertToMogoAccount(item, subname)
}

func (pl *Mongo) HasAccount() bool {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "account"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	nums, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}
	if nums > 0 {
		return true
	}
	return false
}

func (pl *Mongo) GetNodeListFromDatabase() MetaData.NodeList {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "nodelist"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.NodeList
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}

func (pl *Mongo) InsertOrUpdateNodeList(item MetaData.NodeList) {
	pl.dropNodelist()
	pl.saveNodelistToDatabase(item)
}

func (pl *Mongo) dropNodelist() {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	tmp := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	index9 := tmp + "-" + "nodelist"
	_ = session.DB("blockchain").C(index9).DropCollection()
}

func (pl *Mongo) saveNodelistToDatabase(item MetaData.NodeList) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "nodelist"
	pl.InsertToMogoNodelist(item, subname)
}

func (pl *Mongo) HasNodelist() bool {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "nodelist"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	nums, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}
	if nums > 0 {
		return true
	}
	return false
}

func (pl *Mongo) GetIdentityTransListFromDatabase() MetaData.IdentityTransList {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "identityTrans"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.IdentityTransList
	c := session.DB("blockchain").C(subname)
	err := c.Find(nil).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}

func (pl *Mongo) InsertOrUpdateIdentityTransList(item MetaData.IdentityTransList) {
	pl.dropIdentityTransList()
	pl.saveIdentityTransListToDatabase(item)
}

func (pl *Mongo) dropIdentityTransList() {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	tmp := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	index9 := tmp + "-" + "identityTrans"
	_ = session.DB("blockchain").C(index9).DropCollection()
}

func (pl *Mongo) saveIdentityTransListToDatabase(item MetaData.IdentityTransList) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "identityTrans"
	pl.InsertToMogoIdentityTransList(item, subname)
}

func (pl *Mongo) HasIdentityTransList() bool {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "identityTrans"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	nums, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}
	if nums > 0 {
		return true
	}
	return false
}

func (pl *Mongo) SaveNodeIdentityTransToDatabase(item MetaData.IdentityTransformation) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "identity"
	pl.InsertToMogoNodeIdentityTrans(item, subname)
}

func (pl *Mongo) GetOneNodeIdentityTransFromDatabase(typ, key, value string) MetaData.IdentityTransformation {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.IdentityTransformation
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value}).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}

func (pl *Mongo) GetOneNodeIdentityTransByTwoKeyFromDatabase(typ, key, value, key2, value2 string) MetaData.IdentityTransformation {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.IdentityTransformation
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{key: value, key2: value2}).One(&item)
	if err != nil {
		log.Println(err)
	}
	return item
}
