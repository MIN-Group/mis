/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/24 下午11:03
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package db

import (
	"encoding/base64"
	"errors"
	"fmt"
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto"
	"MIS-BC/security/minsecurity/identity"
	"MIS-BC/security/minsecurity/identity/persist/sqlite"

	"MIS-BC/security/minsecurity/crypto/cert"
	"reflect"
	"time"
)

func PersistIdentity(identity *identity.Identity) (bool, error) {
	id, err := GetIdentityByNameFromStorage(identity.Name)
	if id != nil {
		return false, errors.New("The name has existed!!!")
	}

	conn, err := sqlite.GetConn()
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
		certStr, err = identity.Cert.ToPem(nil, sec.SM4ECB)
		if err != nil {
			return false, err
		}
	}

	if identity.PrikeyRawByte != nil{
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

func GetIdentityByNameFromStorage(name string) (*identity.Identity, error) {
	db, err := sqlite.GetConn()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	return getIdentityByNameFromStorage(name, db)
}

func GetAllIdentityFromStorage() ([]*identity.Identity, error) {
	db, err := sqlite.GetConn()
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

func SetDefaultIdentityByNameInStorage(name string) (bool, error) {
	for i := 0 ; i <4 ; i++{
		res, err := setDefaultIdentityByNameInStorage(name)
		if i == 3{
			return res, err
		}

		if res == true{
			return res, err
		}else{
			if err.Error() == "database is locked"{
				time.Sleep(time.Millisecond*50)
			}else{
				return res, err
			}
		}
	}

	return false, nil
}

func setDefaultIdentityByNameInStorage(name string) (bool, error) {
	db, err := sqlite.GetConn()
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

func GetDefaultIdentityFromStorage() (*identity.Identity, error) {
	db, err := sqlite.GetConn()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * from identityinfo where is_default= ?", 1)
	return getDefaultIdentityFromStorage(row)
}

func DeleteIdentityByName(name string) (bool, error) {
	db, err := sqlite.GetConn()
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
