/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/18 下午10:50
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package db

import (
	"database/sql"
	"encoding/base64"
	sec "MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/identity"
)

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

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, sec.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, sec.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != "" {
		err = cert.FromPem(certStr, nil, sec.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: sec.PublicKeyAlgorithm(algo), SignatureAlgorithm: sec.SignatureAlgorithm(sign)}}

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

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, sec.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, sec.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != "" {
		err = cert.FromPem(certStr, nil, sec.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: sec.PublicKeyAlgorithm(algo), SignatureAlgorithm: sec.SignatureAlgorithm(sign)}}

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

	pubKey, err := crypto.UnMarshalPublicKey(pubByte, sec.PublicKeyAlgorithm(algo))
	if err != nil {
		return nil, err
	}

	priByte, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}
	priKey, err := crypto.UnMarshalPrivateKey(priByte, sec.PublicKeyAlgorithm(algo))

	cert := cert.Certificate{}
	if certStr != ""{
		err = cert.FromPem(certStr, nil, sec.SM4ECB)
		if err != nil {
			return nil, err
		}
	}

	priKeyByte, err := base64.StdEncoding.DecodeString(prikeyRawByte)
	if err != nil {
		return nil, err
	}

	id := &identity.Identity{Name: name, Pubkey: pubKey, Prikey: priKey, Passwd: pass, Cert: cert, PrikeyRawByte: priKeyByte, KeyParam: identity.KeyParam{PublicKeyAlgorithm: sec.PublicKeyAlgorithm(algo), SignatureAlgorithm: sec.SignatureAlgorithm(sign)}}

	return id, nil
}
