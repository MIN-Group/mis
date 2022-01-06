/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/17 下午7:03
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sqlite

import (
	"fmt"
	"testing"
)

func TestSqliteOpen(t *testing.T) {
	//测试默认的数据库文件的连接
	OpenDefault()
	res, err := PathExists("/home/xzw/min/identity/identity.db")
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("no target db file!!")
	}

	//测试自定义的数据库文件的连接
	homePath, err := Home()

	if err != nil {
		fmt.Println(err)
		panic("The project can't get the system home path!!!")
	}

	filePath := homePath + "/test/"
	Open(filePath)

	res1, err := PathExists("/home/xzw/test/identity.db")
	if err != nil {
		t.Error(err)
	}
	if !res1 {
		t.Error("no target db file!!!")
	}
}

func TestSqliteSetPasswd(t *testing.T) {
	//测试数据库密码修改
	s := "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
	fmt.Println(s)
	fmt.Println(len(s))
	result := SetPasswd(s)
	if result != nil {
		t.Log(result)
	}
	//测试数据库连接
	_, err := GetConn()
	if err != nil {
		fmt.Println(err)
		panic("Connection Failure!!!")
	}
}
