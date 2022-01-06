/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/20 上午2:04
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

import "encoding/json"

type NetResponse struct {
	innerNetResponse
	conn Connect
}

type innerNetResponse struct {
	Code     int
	Header   map[string]interface{}
	Data     []byte
	ErrorMsg string
	Msg      string
}

func NewNetResponse() *NetResponse {
	var response NetResponse
	response.Header = make(map[string]interface{})
	return &response
}

func (net *NetResponse) Write(b []byte) error {
	return net.conn.Write(b)
}

func (net *NetResponse) SetConnect(connect Connect) {
	net.conn = connect
}

func (net *NetResponse) GetJsonInfo() []byte {
	var inn innerNetResponse
	inn.Code = net.Code
	inn.Header = net.Header
	inn.ErrorMsg = net.ErrorMsg
	inn.Data = net.Data
	inn.Msg = net.Msg
	buf, _ := json.Marshal(inn)
	return buf
}
