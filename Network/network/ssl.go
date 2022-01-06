/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/18 上午2:26
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

import (
	"MIS-BC/common"
	km "MIS-BC/security/keymanager"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var secKeyMap *sync.Map

func init() {
	secKeyMap = new(sync.Map)
}

func startSSLerver(sslAddressOrPrefix string, key *km.KeyManager) {
	net := TCPNet{}
	net.listens(sslAddressOrPrefix, key)
	common.Logger.Info("The ssl server start success")
	for {
		conn, err := net.AcceptTCP()
		if err != nil {
			common.Logger.Fatal("The ssl server accepts failed, because the reason is ", err.Error())
		}
		_, err = conn.readSSL()
		if err != nil {
			common.Logger.Error("The ssl server handles setup message failed...")
		}
	}
}

func sslSetup(connect Connect, pubkey string, key *km.KeyManager, serverKey *km.KeyManager) (string, error) {
	secretKey := fmt.Sprintf("%16v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1e16))
	req := setupMessage{Nonce: pubkey, Secretkey: secretKey}
	data, _ := json.Marshal(req)

	encData, err := serverKey.Encrypt(string(data))
	if err != nil {
		errMsg := "SSL client encrypt setup message failed, the reason is" + err.Error()
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}
	mes := message{IsEncrypted: encrypted, MType: setup, Data: []byte(encData)}
	mesData, _ := json.Marshal(mes)

	err = connect.write(mesData)
	if err != nil {
		errMsg := "SSL client send message failed, the reason is" + err.Error()
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	resData, err := connect.readSSL()
	if err != nil {
		errMsg := "SSL client received message failed..."
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	decRes, err := km.SM4Decrypt(secretKey, resData)
	if err != nil {
		errMsg := "SSL client decrypt message failed, the reason is" + err.Error()
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}
	var resp respMessage
	err = json.Unmarshal(decRes, &resp)
	if err != nil {
		errMsg := "SSl setup failed, the server response is illegal, which is " + err.Error()
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	if resp.Code != success {
		errMsg := "SSL setup failed, the server doesn't accept the request, and the status code is" + strconv.Itoa(int(resp.Code))
		common.Logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	return secretKey, nil
}

func handleSetup(req message, key *km.KeyManager, connInfo string) ([]byte, error) {
	// TODO 如果不是加密接口 或者 不是建立加密的类型 或者 数据空 那么有错误
	if req.IsEncrypted == not_encrypted || req.MType != setup || req.Data == nil {
		resInfo, _ := json.Marshal(req)
		errMsg := "SSl setup failed, the client request is illegal, which is " + string(resInfo)
		common.Logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	// 对Data 进行解密
	decData, err := key.Decrypt(string(req.Data))
	// TODO 如果解密失败 那么返回BAD——REQUEST
	if err != nil {
		var resp respMessage
		resp.Code = bad_request
		respData, _ := json.Marshal(resp)
		resInfo, _ := json.Marshal(req)
		errMsg := "SSL setup failed, the client request can't be decrpted, which is " + string(resInfo)
		return respData, errors.New(errMsg)
	}
	common.Logger.DebugWithConn(connInfo, "SSL server received setup:", string(decData))
	var setup setupMessage
	// 解密后的数据 json-->struct
	err = json.Unmarshal(decData, &setup)
	if err != nil {
		errMsg := "SSl setup failed, the client request is illegal, which is " + err.Error()
		var resp respMessage
		resp.Code = bad_request
		respData, _ := json.Marshal(resp)
		return respData, errors.New(errMsg)
	}
	if setup.Secretkey == "" || setup.Nonce == "" {
		errMsg := "SSl setup failed, the client request is illegal, which is " + "the parameter is nil"
		var resp respMessage
		resp.Code = bad_request
		respData, _ := json.Marshal(resp)
		return respData, errors.New(errMsg)
	}

	secKeyMap.Store(setup.Nonce, setup.Secretkey)

	var resp respMessage
	resp.Code = success
	respData, _ := json.Marshal(resp)

	response, err := km.SM4Encrypt(setup.Secretkey, respData)
	if err != nil {
		errMsg := "SSl setup failed, the client secretkey is illegal, which is " + setup.Secretkey
		return nil, errors.New(errMsg)
	}
	return response, nil
}

func sslPacket(code int, data []byte, isEnc enc, mtype messagetype) ([]byte, error) {
	var mes message
	mes.MType = mtype
	mes.IsEncrypted = isEnc
	mes.Data = data
	mes.Code = code

	res, _ := json.Marshal(mes)
	return res, nil
}

func sslDePacket(data []byte, role int, secretKey string) (string, enc, int, []byte, error) {
	var res message
	// json-->struct
	err := json.Unmarshal(data, &res)
	// TODO 解析错误 返回400
	if err != nil {
		errMsg := "SSl depacket failed, the client request is illegal, which is " + err.Error()
		common.Logger.Error(errMsg)
		return "", not_encrypted, 400, nil, errors.New(errMsg)
	}
	// 没有加密直接返回 200 和源数据res.Data
	if res.IsEncrypted == not_encrypted {
		return "", not_encrypted, 200, res.Data, nil
	} else if res.IsEncrypted == encrypted {
		// 加密 判断role
		if role == 0 {
			// TODO 跟据解析出来的数据中的公钥ID从本地MAP获取公钥 对数据SM4解密
			if res.Nonce == "" {
				errMsg := "SSL depacket failed, the Nonce doesn't exist"
				common.Logger.Error(errMsg)
				return "", encrypted, 400, nil, errors.New(errMsg)
			}
			// 取出密钥
			if _, flag := secKeyMap.Load(res.Nonce); flag == false {
				common.Logger.Error("SSL depacket failed, the secKey of Nonce doesn't exist")
				return "", encrypted, 400, nil, errors.New("SSL depacket failed, the pubKey doesn't exist")
			}
			secretKey, _ := secKeyMap.Load(res.Nonce)
			// SM4解密
			ans, err := km.SM4Decrypt(secretKey.(string), res.Data)
			if err != nil {
				errMsg := "SSL depacket failed, the secretKey can't use to decrypt the data, the err is " + err.Error()
				common.Logger.Error(errMsg)
				return "", encrypted, 400, nil, errors.New(errMsg)
			}
			// 返回密钥和200
			return secretKey.(string), encrypted, 200, ans, nil

		} else {
			// TODO 直接使用传入的密钥 对数据SM4解密
			if secretKey == "" {
				errMsg := "SSL depacket failed, the secretkey doesn't exist"
				common.Logger.Error(errMsg)
				return "", encrypted, 400, nil, errors.New(errMsg)
			}
			ans, err := km.SM4Decrypt(secretKey, res.Data)
			if err != nil {
				errMSg := "SSL depacket failed, the secretKey can't use to decrypt the data, the err is " + err.Error()
				common.Logger.Error(errMSg)
				return "", encrypted, 400, nil, errors.New(errMSg)
			}
			return secretKey, encrypted, 200, ans, nil

		}

	} else {
		return "", not_encrypted, 400, nil, nil
	}

}
