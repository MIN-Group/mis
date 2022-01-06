/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/20 上午2:04
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package network

import "encoding/json"

type NetRequest struct {
	innerNetRequest
	conn       Connect
	attributes map[string]interface{}
	model      interface{}
}

type innerNetRequest struct {
	Method     string
	Type       string
	Command    string
	Header     map[string]interface{}
	Parameters []byte
}

func NewNetRequest() *NetRequest {
	return &NetRequest{
		innerNetRequest: innerNetRequest{
			Header: make(map[string]interface{}),
		},
		attributes: make(map[string]interface{}),
	}
}

// SetJsonInfo
// @Title    		SetJsonInfo
// @Description   	设置json串
// @Param			buf []byte 字节缓存
// @Return			error 错误
func (req *NetRequest) SetJsonInfo(buf []byte) error {
	inn := innerNetRequest{}
	// buf:json--> innerNetRequest
	err := json.Unmarshal(buf, &inn)
	if err != nil {
		return err
	}
	// TODO 设置inn 到 request
	req.innerNetRequest = inn
	return nil
}

func (req *NetRequest) Read() ([]byte, error) {
	return req.Read()
}

func (req *NetRequest) GetRemote() string {
	return req.conn.GetRemote()
}

func (req *NetRequest) GetConnection() Connect {
	return req.conn
}

func (req *NetRequest) SetConnection(conn Connect) {
	req.conn = conn
}

func (req *NetRequest) SetAttribute(key string, value interface{}) {
	if req.attributes == nil {
		req.attributes = make(map[string]interface{})
	}

	req.attributes[key] = value
}

func (req *NetRequest) GetAttribute(key string) interface{} {
	return req.attributes[key]
}

func (req *NetRequest) GetModel() interface{} {
	return req.model
}

func (req *NetRequest) SetModel(model interface{}) {
	req.model = model
}
