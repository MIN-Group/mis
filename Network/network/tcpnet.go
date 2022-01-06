/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/17 下午7:27
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

import (
	"MIS-BC/Network/network/encoding"
	"MIS-BC/common"
	"MIS-BC/security/keymanager"
	"encoding/json"
	"errors"
	"net"
)

// TCPNet TCP连接结构体
type TCPNet struct {
	listener  *net.TCPListener       // 网络监听器
	key       *keymanager.KeyManager // 密钥管理
	serverKey *keymanager.KeyManager // just pubkey can be used
}

func (tnet *TCPNet) Listen(address string) error {
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	tnet.listener, err = net.ListenTCP("tcp", tcpAddr)
	//tnet.listener, err = net.Listen("tcp", "localhost:"+port)
	return err
}

func (tnet *TCPNet) listens(sslAddress string, key *keymanager.KeyManager) error {
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", sslAddress)
	tnet.listener, err = net.ListenTCP("tcp", tcpAddr)
	//tnet.listener, err = net.Listen("tcp", sslAddress)
	tnet.key = key
	return err
}

func (tnet *TCPNet) Listens(address string, sslAddress string, key *keymanager.KeyManager) error {
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	tnet.listener, err = net.ListenTCP("tcp", tcpAddr)
	//tnet.listener, err = net.Listen("tcp", address)
	if err != nil {
		common.Logger.Error("The server listen failed, because the reason is ", err.Error())
	}
	go startSSLerver(sslAddress, key)
	return err
}

// @Title    		Accept
// @Description   	阻塞接收加密连接
// @Return			connect 连接接口
// @Return			error   错误
func (tnet *TCPNet) AcceptTCP() (Connect, error) {
	// 得到连接
	conn, err := tnet.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	// 建立TCP连接
	tcpConn := TCPConnect{connect: conn, key: tnet.key, role: 0}
	return &tcpConn, nil
}

func (tnet *TCPNet) Dial(host string) (Connect, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	res := TCPConnect{connect: conn, key: tnet.key, role: 1}
	return &res, nil
}

func (tnet *TCPNet) Dials(host string, sslHost string, key *keymanager.KeyManager, serverPubKey []byte) (Connect, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	tnet.key = key
	serverkey := keymanager.KeyManager{}
	serverkey.Init()
	err = serverkey.SetPubkey(string(serverPubKey))
	if err != nil {
		return nil, err
	}

	connector, err := tnet.Dial(sslHost)
	if err != nil {
		return nil, err
	}
	secretKey, err := sslSetup(connector, key.GetPubkey(), key, &serverkey)
	if err != nil {
		return nil, err
	}

	res := TCPConnect{connect: conn, key: tnet.key, serverKey: &serverkey, isEncrypt: true, secretKey: secretKey, role: 1}
	return &res, nil
}

const (
	ShortConn = "short"
	LongConn  = "long"
)

type TCPConnect struct {
	connect   *net.TCPConn           // 网络连接
	isEncrypt bool                   // 是否加密
	key       *keymanager.KeyManager // 密钥管理结构体
	serverKey *keymanager.KeyManager // just pubkey can be used
	secretKey string                 // 密钥
	role      int                    // 指示是使用传入数据中的公钥ID从本地取出密钥解密 还是直接传入secretKey
}

func (conn *TCPConnect) Read() ([]byte, error) {
	// 四个字节存储大小 防粘包
	data, err := encoding.DecodeTcp(conn.connect)
	if err != nil {
		common.Logger.Error("Error read normal message failed...")
		return nil, err
	}
	// 使用密钥解密数据
	secretKey, encry, _, ans, err := sslDePacket(data, conn.role, conn.secretKey)
	if err != nil {
		common.Logger.Error("Error depacket message failed...")
		return nil, err
	}
	// 如果role == 0说明是非对称加密 得到密钥 加密方式变成encrypted
	if conn.role == 0 {
		conn.isEncrypt = encry == encrypted
		conn.secretKey = secretKey
	}
	// 返回得到的数据
	return ans, nil
}

func (conn *TCPConnect) readSSL() ([]byte, error) {
	data, err := encoding.DecodeTcp(conn.connect)
	if err != nil {
		common.Logger.Error("SSL received message, but decode failed...")
		return nil, err
	}
	common.Logger.DebugWithConn(conn.GetRemote(), "SSL accept message, which is ", string(data))

	var req message
	err = json.Unmarshal(data, &req)
	if err != nil {
		errMsg := "SSL received message failed, the client request is illegal, which is " + string(data)
		common.Logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if req.MType == setup {
		res, err := handleSetup(req, conn.key, conn.GetRemote())
		if err != nil {
			common.Logger.Error("SSL setup failed...")
			//return nil, err
		}
		ress, err := sslPacket(200, res, encrypted, resp)
		if err != nil {
			common.Logger.Error("SSL setup failed...")
			//return nil,err
		}
		conn.write(ress)
		conn.Close()
		return nil, err
	} else if req.MType == resp {
		return []byte(req.Data), nil
	} else {
		return nil, errors.New("SSL Server can't handle this message type")
	}

}

func (conn *TCPConnect) Write(b []byte) error {
	var mes message

	if conn.isEncrypt == false {
		mes.IsEncrypted = not_encrypted
		mes.Data = b
		mes.MType = normal
		mes.Code = 200
		data, _ := json.Marshal(mes)
		return conn.write(data)
	} else {
		mes.IsEncrypted = encrypted
		mes.MType = normal
		if conn.secretKey == "" {
			mes.Code = 500
			data, _ := json.Marshal(mes)
			//fmt.Println(string(data))
			conn.write(data)
			return errors.New("SSL depacket failed, the secretKey doesn't exist")
		}
		ans, err := keymanager.SM4Encrypt(conn.secretKey, b)
		if err != nil {
			mes.Code = 500
			data, _ := json.Marshal(mes)
			conn.write(data)
			return errors.New("SSL depacket failed, the secretKey can't use to decrypt the data, the err is " + err.Error())
		}
		mes.Data = ans
		if conn.role == 1 {
			mes.Nonce = conn.key.GetPubkey()
		}
		mes.Code = 200
		data, _ := json.Marshal(mes)
		return conn.write(data)
	}

}

func (conn *TCPConnect) write(b []byte) error {
	new_data, err := encoding.Encode(b)
	_, err = conn.connect.Write(new_data)
	if err != nil {
		return err
	}

	return nil
}

func (conn *TCPConnect) Close() {
	conn.connect.Close()
}

func (conn *TCPConnect) GetRemote() string {
	return conn.connect.RemoteAddr().String()
}

func (conn *TCPConnect) GetConn() *net.TCPConn {
	return conn.connect
}
