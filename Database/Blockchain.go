/**
 * @Author: xzw
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/6/15 下午4:00
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package MongoDB

import (
	"MIS-BC/MetaData"
	"MIS-BC/common"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"log"
	"strconv"
)

func (pl *Mongo) QueryHeight() int {
	var height = -1
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "blockstate"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var item MetaData.BlockState
	c := session.DB("blockchain").C(subname)
	count, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		c := session.DB("blockchain").C(subname)
		err := c.Find(nil).One(&item)
		if err != nil {
			log.Println(err)
		}
		height = item.Height
	}

	return height
}

func (pl *Mongo) SetHeight(height int) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "blockstate"
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	count, err := c.Find(nil).Count()
	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		item := make(map[string]interface{})
		c := session.DB("blockchain").C(subname)
		err := c.Find(nil).One(item)
		if err != nil {
			log.Println(err)
		}
		selector := bson.M{"_id": item["_id"]}
		data := bson.M{"$set": bson.M{"height": height}}

		err = c.Update(selector, data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		blockstate := MetaData.BlockState{Height: height}
		pl.saveBlockstateToDatabase(blockstate)
	}
}

func (pl *Mongo) saveBlockstateToDatabase(item MetaData.BlockState) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "blockstate"
	pl.InsertToMogoBlockstate(item, subname)
}

func (pl *Mongo) SaveBCStatusToDatabase(item MetaData.BCStatus) {
	typ := "bcstatus"
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ

	session := pl.pool.AcquireSession()
	defer session.Release()

	c := session.DB("blockchain").C(subname)
	count, err := c.Find(nil).Count()
	if err != nil {
		common.Logger.Error(err)
	}
	if count > pl.CacheNumber {
		var bs MetaData.BCStatus
		c.Find(nil).Sort("timestamp").Limit(1).One(&bs)
		c.RemoveAll(bson.M{"timestamp": bs.Timestamp})
	}
	pl.InsertToMogoBCStatus(item, subname)
}

func (pl *Mongo) GetBCStatusFromDatabase() MetaData.BCStatus {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + "bcstatus"
	session := pl.pool.AcquireSession()
	defer session.Release()

	var item MetaData.BCStatus
	c := session.DB("blockchain").C(subname)
	count, err := c.Find(nil).Count()
	if count > pl.CacheNumber {
		var bs MetaData.BCStatus
		c.Find(nil).Sort("timestamp").Limit(1).One(&bs)
		c.RemoveAll(bson.M{"timestamp": bs.Timestamp})
		err = c.Find(nil).Sort("-timestamp").Skip(pl.CacheNumber * 9 / 10).Limit(1).One(&item)
		if err != nil {
			common.Logger.Error(err)
		}
	} else {
		err = c.Find(nil).Sort("timestamp").Limit(1).One(&item)
		if err != nil {
			common.Logger.Error(err)
		}
		_, err = c.RemoveAll(bson.M{"timestamp": item.Timestamp})
		if err != nil {
			common.Logger.Error(err)
		}
	}

	return item
}

func (pl *Mongo) DeleteBCStatus(key, value string) {
	typ := "bcstatus"
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

func (pl *Mongo) GetAmount() int {
	return pl.QueryHeight() + 1
}

func (pl *Mongo) PushbackBlockToDatabase(block MetaData.BlockGroup) {
	if block.Height == 0 {
		block.CheckHeader = []int{1}
	}
	pl.InsertToMogoBG(block, pl.Pubkey)
	pl.Block = block
	pl.Height = block.Height
	pl.SetHeight(block.Height)
}

func (pl *Mongo) GetBlockFromDatabase(height int) MetaData.BlockGroup {
	session := pl.pool.AcquireSession()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Release()

	var blockgroup MetaData.BlockGroup
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index + "-blockgroup" + "-" + strconv.Itoa(height/(10*10000)))
	err := c.Find(bson.M{"height": height}).One(&blockgroup)
	if err != nil {
		log.Println(err)
	}

	var blocks []MetaData.Block
	c1 := session.DB("blockchain").C(index + "-block" + "-" + strconv.Itoa(height/(10*10000)))
	err = c1.Find(bson.M{"height": height}).All(&blocks)
	if err != nil {
		log.Println(err)
	}

	true_blocks := make([]MetaData.Block, len(blockgroup.CheckHeader))
	for _, v := range blocks {
		true_blocks[v.BlockNum] = v
	}

	blockgroup.Blocks = true_blocks
	return blockgroup
}
