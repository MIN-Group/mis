package MongoDB

import (
	"MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/identity"
	"MIS-BC/security/minsecurity/identity/persist/sqlite"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"
)

var dbFile = "identity.db"
var passwd = "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"

//自定义文件地址,注意要用/结尾,比如:/.mir/identity/
func Open(filePath string, index string) {
	dbPath := filePath
	isExits, err := sqlite.PathExists(dbPath)
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
	realFile := dbPath + index
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

func PersistIdentity(identity *identity.Identity, Path, File string) (bool, error) {
	id, err := GetIdentityByNameFromStorage(identity.Name, Path, File)
	if id != nil {
		return false, errors.New("The name has existed!!!")
	}

	conn, err := GetConn(Path, File)
	defer conn.Close()
	if err != nil {
		return false, err
	}
	pubStr := ""
	priStr := ""
	algo := 0
	certStr := ""
	prikeyRawByte := ""
	if identity.Prikey != nil {
		algorithm, priByte := crypto.MarshalPrivateKey(identity.Prikey)
		priStr = base64.StdEncoding.EncodeToString(priByte)
		algo = int(algorithm)
	}
	if identity.Pubkey != nil {
		algorithm, pubByte := crypto.MarshalPublicKey(identity.Pubkey)
		pubStr = base64.StdEncoding.EncodeToString(pubByte)
		algo = int(algorithm)
	}
	algo = int(identity.KeyParam.PublicKeyAlgorithm)
	sign := int(identity.KeyParam.SignatureAlgorithm)

	if reflect.DeepEqual(identity.Cert, cert.Certificate{}) {
		certStr = ""
	} else {
		certStr, err = identity.Cert.ToPem(nil, minsecurity.SM4ECB)
		if err != nil {
			return false, err
		}
	}

	if identity.PrikeyRawByte != nil {
		prikeyRawByte = base64.StdEncoding.EncodeToString(identity.PrikeyRawByte)
	}

	stmt, err := conn.Prepare("INSERT INTO identityinfo(name, pubkey, prikey, pubkey_algo, signature_algo, pass, cert,prikey_raw_byte) values(?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	_, err = stmt.Exec(identity.Name, pubStr, priStr, algo, sign, "", certStr, prikeyRawByte)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func GetIdentityByNameFromStorage(name string, Path, File string) (*identity.Identity, error) {
	db, err := GetConn(Path, File)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	return getIdentityByNameFromStorage(name, db)
}

func GetAllIdentityFromStorage(Path, File string) ([]*identity.Identity, error) {
	db, err := GetConn(Path, File)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * from identityinfo")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var res []*identity.Identity
	for rows.Next() {
		id, err := getIdentityFromSqlRows(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	return res, nil
}

func SetDefaultIdentityByNameInStorage(name string, Path, File string) (bool, error) {
	for i := 0; i < 4; i++ {
		res, err := setDefaultIdentityByNameInStorage(name, Path, File)
		if i == 3 {
			return res, err
		}

		if res == true {
			return res, err
		} else {
			if err.Error() == "database is locked" {
				time.Sleep(time.Millisecond * 50)
			} else {
				return res, err
			}
		}
	}

	return false, nil
}

func setDefaultIdentityByNameInStorage(name string, Path, File string) (bool, error) {
	db, err := GetConn(Path, File)
	defer db.Close()
	if err != nil {
		return false, err
	}

	tx, err := db.Begin()
	row := tx.QueryRow("select * from identityinfo where name = ?", name)
	id, err := getDefaultIdentityFromStorage(row)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if id == nil {
		tx.Rollback()
		return false, nil
	}

	row = tx.QueryRow("SELECT * from identityinfo where is_default = 1")
	id, err = getDefaultIdentityFromStorage(row)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if id != nil {
		flag, err := cancelDefaultIdentityFromStorage(id.Name, tx)
		if flag == false {
			tx.Rollback()
			return false, err
		}
	}

	flag, err := setDefaultIdentityFromStorage(name, tx)
	if flag == false {
		tx.Rollback()
		return false, err
	}
	tx.Commit()
	return true, nil

}

func GetDefaultIdentityFromStorage(Path, File string) (*identity.Identity, error) {
	db, err := GetConn(Path, File)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * from identityinfo where is_default= ?", 1)
	return getDefaultIdentityFromStorage(row)
}

func DeleteIdentityByName(name string, Path, File string) (bool, error) {
	db, err := GetConn(Path, File)
	defer db.Close()
	if err != nil {
		return false, err
	}
	stmt, err := db.Prepare("delete from identityinfo where name=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name)
	if err != nil {
		return false, err
	}
	return true, nil
}

//注意该函数没有对name的唯一性检查
func setDefaultIdentityFromStorage(name string, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("update identityinfo set is_default = 1 where name= ?")
	if stmt == nil || err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name)
	if err != nil {
		return false, err
	}

	return true, nil
}

//注意该函数没有实现name检查
func cancelDefaultIdentityFromStorage(name string, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("update identityinfo set is_default= 0 where name= ?")
	if stmt == nil || err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getDefaultIdentityFromStorage(row *sql.Row) (*identity.Identity, error) {
	pubStr := ""
	priStr := ""
	algo := 0
	sign := 0
	pass := ""
	certStr := ""
	name := ""
	def := 0
	prikeyRawByte := ""

	err := row.Scan(&name, &pubStr, &priStr, &algo, &sign, &pass, &certStr, &def, &prikeyRawByte)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}
	pubByte, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, minsecurity.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, minsecurity.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != "" {
		err = cert.FromPem(certStr, nil, minsecurity.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: minsecurity.PublicKeyAlgorithm(algo), SignatureAlgorithm: minsecurity.SignatureAlgorithm(sign)}}

	return id, nil
}

func getIdentityByNameFromStorage(name string, db *sql.DB) (*identity.Identity, error) {

	row := db.QueryRow("SELECT * from identityinfo where name = ?", name)

	return getIdentityFromSqlRow(row)
}

func getIdentityFromSqlRows(row *sql.Rows) (*identity.Identity, error) {
	pubStr := ""
	priStr := ""
	algo := 0
	sign := 0
	pass := ""
	certStr := ""
	name := ""
	def := 0
	prikeyRawByte := ""

	err := row.Scan(&name, &pubStr, &priStr, &algo, &sign, &pass, &certStr, &def, &prikeyRawByte)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}
	pubByte, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, minsecurity.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, minsecurity.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != "" {
		err = cert.FromPem(certStr, nil, minsecurity.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: minsecurity.PublicKeyAlgorithm(algo), SignatureAlgorithm: minsecurity.SignatureAlgorithm(sign)}}

	return id, nil
}

func getIdentityFromSqlRow(row *sql.Row) (*identity.Identity, error) {
	pubStr := ""
	priStr := ""
	algo := 0
	sign := 0
	pass := ""
	certStr := ""
	name := ""
	def := 0
	prikeyRawByte := ""

	err := row.Scan(&name, &pubStr, &priStr, &algo, &sign, &pass, &certStr, &def, &prikeyRawByte)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}
	pubByte, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, minsecurity.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, minsecurity.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != "" {
		err = cert.FromPem(certStr, nil, minsecurity.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: minsecurity.PublicKeyAlgorithm(algo), SignatureAlgorithm: minsecurity.SignatureAlgorithm(sign)}}

	return id, nil
}

func GetConn(Path, File string) (*sql.DB, error) {
	realFile := Path + File
	dbName := fmt.Sprintf("%s?_pragma_key=x'%s'&_busy_timeout=900000000&_pragma_cipher_page_size=4096", realFile, passwd)
	//dbName := fmt.Sprintf("%s?_busy_timeout=9999999&cache=shared&mode=rwc", realFile)
	db, err := sql.Open("sqlite3", dbName)
	db.SetMaxOpenConns(1)
	return db, err
}
