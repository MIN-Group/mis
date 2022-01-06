/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/1/24 下午6:16
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package sha512

import (
	"crypto/sha512"
)
import sec "MIS-BC/security/minsecurity"

func init(){
	sec.RegisterHash(sec.SHA512, sha512.New)
}