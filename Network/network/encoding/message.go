/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/17 下午8:34
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package encoding

import (
	"MIS-BC/common"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

// Encode 将消息编码
func Encode(message []byte) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// DecodeTcp Decode 解码消息
func DecodeTcp(conn net.Conn) ([]byte, error) {
	// 读取消息的长度
	buf := make([]byte, 4)           //4字节长度缓冲
	_, err := io.ReadFull(conn, buf) // 读取前4个字节的数据
	if err != nil {
		common.Logger.Error("Error IO reading failed, the reason is ", err.Error())
		return make([]byte, 0), err
	}
	length := binary.LittleEndian.Uint32(buf) //
	// 读取真正的消息数据
	pack := make([]byte, length)
	_, err = io.ReadFull(conn, pack)
	if err != nil {
		common.Logger.Error("Error IO reading failed, the reason is ", err.Error())
		return make([]byte, 0), err
	}

	return pack, nil
}
