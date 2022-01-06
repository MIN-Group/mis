/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/17 下午7:02
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"os"
)

var dbFile = "identity.db"
var dbPath = ""
var passwd = "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"

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

/*
* 进行数据库文件的连接测试，如果表不存在则新建
* 数据库文件默认写入地址为~//identity/identity.db
*
 */
func OpenDefault() {
	homePath, err := Home()

	if err != nil {
		fmt.Println(err)
		panic("The project can't get the system home path!!!")
	}

	homePath = homePath + "/.mir/identity/"
	dbPath = homePath

	isExits, err := PathExists(homePath)
	if err != nil {
		fmt.Println(err)
		panic("The sql db file can't be created!!!")
	}
	if !isExits {
		err := os.MkdirAll(homePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			panic("The sql db file can't be created!!!")
		}
	}

	realFile := homePath + dbFile
	dbName := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096&_busy_timeout=9999999", realFile, passwd)
	//dbName := fmt.Sprintf("%s", realFile)

	db, err := sql.Open("sqlite3", dbName)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	if err != nil {
		fmt.Println(err)
		panic("The sql file can't be opened")
	}

	_, err = db.Exec(table_create)
	if err != nil {
		fmt.Println(err)
		panic("Create table fails!!!")
	}

}

//自定义文件地址,注意要用/结尾,比如:/.mir/identity/
func Open(filePath string) {
	dbPath = filePath
	isExits, err := PathExists(dbPath)
	if err != nil {
		fmt.Println(err)
		panic("The sql db file can't be created!!!")
	}
	if !isExits {
		err := os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			panic("The sql db file can't be created!!!")
		}
	}
	realFile := dbPath + dbFile
	dbName := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096&_busy_timeout=9999999", realFile, passwd)
	//dbName := fmt.Sprintf("%s?_busy_timeout=9999999&cache=shared&mode=rwc", realFile)

	db, err := sql.Open("sqlite3", dbName)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic("The sql file can't be opened")
	}

	_, err = db.Exec(table_create)
	if err != nil {
		fmt.Println(err)
		panic("Create table fails!!!")
	}
}

//设置数据库连接密码
func SetPasswd(pass string) error {
	if len(pass) != 64 {
		return errors.New("Invalid length")
	}
	passwd = pass
	return nil
}

func GetConn() (*sql.DB, error) {
	realFile := dbPath + dbFile
	dbName := fmt.Sprintf("%s?_pragma_key=x'%s'&_busy_timeout=900000000&_pragma_cipher_page_size=4096", realFile, passwd)
	//dbName := fmt.Sprintf("%s?_busy_timeout=9999999&cache=shared&mode=rwc", realFile)
	db, err := sql.Open("sqlite3", dbName)
	db.SetMaxOpenConns(1)
	return db, err
}
