package Node

import (
	"MIS-BC/MetaData"
	"MIS-BC/common"
	"MIS-BC/security/code"
	"MIS-BC/security/minsecurity"
	"MIS-BC/security/minsecurity/crypto/cert"
	"MIS-BC/security/minsecurity/crypto/sm2"
	"context"
	"fmt"
	"time"
)

type RpcServer struct {
	node *Node
}

type CommonResponse struct {
	Code    int
	Message string
	Data    interface{}
}

type Args struct {
	IdentityIdentifier string
	Pubkey             string
	Passwd             string
}

func (rpcServer *RpcServer) IdentityRegistryforTest(ctx context.Context, args Args, reply *CommonResponse) error {
	if args.IdentityIdentifier == "" || args.Pubkey == "" || args.Passwd == "" {
		reply = &CommonResponse{Code: code.LESS_PARAMETER, Message: "缺少字段", Data: nil}
		return fmt.Errorf("\"缺少字段\"")
	}

	if rpcServer.node.mongo.HasIdentityData("identityidentifier", args.IdentityIdentifier) {
		common.Logger.Error("数据库已经存在该身份标识，注册失败", args.IdentityIdentifier)
		reply.Code = code.BAD_REQUEST
		reply.Message = "数据库已经存在该身份标识"
		reply.Data = nil
		return fmt.Errorf("数据库已经存在该身份标识")
	} else if rpcServer.node.mongo.HasIdentityData("pubkey", args.Pubkey) {
		common.Logger.Error("用户公钥重复，注册失败")
		reply.Code = code.BAD_REQUEST
		reply.Message = "用户公钥重复，注册失败"
		reply.Data = nil
		return fmt.Errorf("用户公钥重复，注册失败")
	} else {
		var transaction MetaData.Identity
		transaction.Type = "identity-act"
		transaction.Command = "Registry"
		transaction.IdentityIdentifier = args.IdentityIdentifier
		transaction.KeyParam = MetaData.KeyParam{0, 0}
		transaction.Pubkey = args.Pubkey
		transaction.Passwd = args.Passwd
		transaction.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		transaction.IsValid = code.VALID
		transaction.IPIdentifier = ""
		transaction.ModifyRecords = append(transaction.ModifyRecords, MetaData.ModifyRecord{Type: "identity-act",
			Command: "Registry", Timestamp: time.Now().Format("2006-01-02 15:04:05")})

		// 填充证书内容
		pub := sm2.Sm2PublicKey{}
		pub.SetBytes([]byte(transaction.Pubkey))
		var pubkey minsecurity.PublicKey = &pub
		cert := cert.Certificate{}
		cert.Version = 0
		cert.SerialNumber = 1
		cert.PublicKey = pubkey
		cert.SignatureAlgorithm = 0
		cert.PublicKeyAlgorithm = 0
		cert.IssueTo = transaction.IdentityIdentifier
		cert.Issuer = "/root"
		cert.NotBefore = time.Now().Unix()
		cert.NotAfter = time.Now().AddDate(1, 0, 0).Unix()
		cert.KeyUsage = minsecurity.CertSign
		cert.IsCA = false
		cert.Timestamp = time.Now().Unix()

		pri := sm2.Sm2PrivateKey{}
		pri.SetBytes([]byte(rpcServer.node.keyManager.GetPriKey()))
		var prikey minsecurity.PrivateKey = &pri
		err := cert.SignCert(prikey)
		if err != nil {
			common.Logger.Error(err)
		}

		c, err := cert.ToPem([]byte(transaction.Passwd), 0)
		if err != nil {
			common.Logger.Error("Certificate issuance failed：", err)
			reply.Code = code.BAD_REQUEST
			reply.Message = "Certificate issuance failed"
			return fmt.Errorf("Certificate issuance failed")
		} else {
			var transactionHeader MetaData.TransactionHeader
			transactionHeader.TXType = MetaData.IdentityAction
			transaction.Cert = c
			reply.Code = code.SUCCESS
			reply.Message = "身份注册成功"
			reply.Data = transaction
			common.Logger.Info("身份申请注册中")
			go rpcServer.node.txPool.PushbackTransaction(transactionHeader, &transaction)
			rpcServer.node.registryList[transaction.Pubkey] = rpcServer.node.mongo.Height
			// common.Logger.Info("当前身份：", node.network.Keychain.GetAllIdentities())
		}
	}
	return nil
}

func (rpcServer *RpcServer) GetOneIdentityInfByIdentityIdentifierforTest(ctx context.Context, args Args, reply *CommonResponse) error {
	if args.IdentityIdentifier == "" {
		reply = &CommonResponse{Code: code.LESS_PARAMETER, Message: "缺少字段", Data: nil}
		return fmt.Errorf("\"缺少字段\"")
	}
	if !rpcServer.node.mongo.HasIdentityData("identityidentifier", args.IdentityIdentifier) {
		reply.Code = code.NOT_FOUND
		reply.Message = "不存在该身份"
		reply.Data = MetaData.Identity{}
		return fmt.Errorf("不存在该身份")
	} else {
		identity := rpcServer.node.mongo.GetOneIdentityFromDatabase("identityidentifier", args.IdentityIdentifier)

		reply.Code = code.SUCCESS
		reply.Message = "成功获得该身份"
		reply.Data = identity
	}
	common.Logger.Info("身份"+args.IdentityIdentifier+"查询成功")
	return nil
}

func (rpcServer *RpcServer) DestroyByIdentityIdentifierforTest(ctx context.Context, args Args, reply *CommonResponse) error {
	if args.IdentityIdentifier == "" {
		reply = &CommonResponse{Code: code.LESS_PARAMETER, Message: "缺少字段", Data: nil}
		return fmt.Errorf("\"缺少字段\"")
	}

	if rpcServer.node.mongo.HasIdentityData("identityidentifier", args.IdentityIdentifier) {
		var transaction MetaData.Identity
		transaction.Type = "identity-act"
		transaction.Command = "DestroyByIdentityIdentifier"
		transaction.IdentityIdentifier = args.IdentityIdentifier

		var transactionHeader MetaData.TransactionHeader
		transactionHeader.TXType = MetaData.IdentityAction
		rpcServer.node.txPool.PushbackTransaction(transactionHeader, &transaction)

		i := rpcServer.node.mongo.GetOneIdentityFromDatabase("identityidentifier", args.IdentityIdentifier)
		flag, err := rpcServer.node.network.Keychain.DeleteIdentityByName(transaction.IdentityIdentifier, i.Passwd)
		if err != nil {
			common.Logger.Error(err)
		} else if flag == true {
			common.Logger.Info("sqlite删除身份成功")
		} else {
			common.Logger.Info("sqlite删除身份失败")

		}

		reply.Code = code.SUCCESS
		reply.Message = "注销成功"
		reply.Data = nil

	} else {
		reply.Code = code.NOT_FOUND
		reply.Message = "数据库不存在该用户"
		reply.Data = nil
		return fmt.Errorf("数据库不存在该用户")
	}

	common.Logger.Info("当前身份：", rpcServer.node.network.Keychain.GetAllIdentities())
	return nil
}

type restPasswdArgs struct {
	IdentityIdentifier string
	Previous		   string
	Passwd             string
}


func (rpcServer *RpcServer) ResetPasswordforTest(ctx context.Context, args Args, reply *CommonResponse) error {
	if args.IdentityIdentifier == "" || args.Passwd == "" {
		reply = &CommonResponse{Code: code.LESS_PARAMETER, Message: "缺少字段", Data: nil}
		return fmt.Errorf("\"缺少字段\"")
	}

	if !rpcServer.node.mongo.HasIdentityData("identityidentifier", args.IdentityIdentifier) {
		reply.Code = code.NOT_FOUND
		reply.Message = "数据库不存在该身份"
		reply.Data = nil
		return fmt.Errorf("数据库不存在该身份")
	} else {
		// identity := node.mongo.GetOneIdentityFromDatabase("identityidentifier", res["IdentityIdentifier"].(string))
		var transaction MetaData.Identity
		transaction.Type = "identity-act"
		transaction.Command = "ResetPassword"
		transaction.IdentityIdentifier = args.IdentityIdentifier
		transaction.Passwd = args.Passwd

		var transactionHeader MetaData.TransactionHeader
		transactionHeader.TXType = MetaData.IdentityAction
		rpcServer.node.txPool.PushbackTransaction(transactionHeader, &transaction)

		reply.Code = code.SUCCESS
		reply.Message = "修改成功"
		reply.Data = nil
	}

	return nil
}