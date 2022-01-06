/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/18 上午1:59
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

type enc int

const (
	encrypted     enc = 0
	not_encrypted enc = 1
)

type messagetype int

const (
	normal messagetype = 0 + iota
	setup
	resp
)

type message struct {
	IsEncrypted enc
	MType       messagetype
	Data        []byte
	Nonce       string
	Code        int
}

type setupMessage struct {
	Nonce     string
	Secretkey string
}

type respcode int

const (
	success        respcode = 200
	bad_request    respcode = 400
	not_authorized respcode = 403
	server_error   respcode = 500
)

type respMessage struct {
	Code respcode
}
