/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/20 上午2:05
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package code

const (
	SUCCESS               = 200
	UNAUTHORIZED          = 401
	BAD_REQUEST           = 400
	FORBIDDEN             = 403
	NOT_FOUND             = 404
	INTERNAL_SERVER_ERROR = 500
)

// 登录状态码
const (
	NORMAL_SUCCESS  = 200
	ROOT_SUCCESS    = 201
	PENDING_SUCCESS = 203

	UNARTHORIZED = 401

	LESS_PARAMETER = 410
	SERVER_WRONG   = 500

	NORMAL_WARNING  = 300
	ROOT_WARNING    = 299
	PENDING_WARNING = 298

	EXPIRED                = 301
	MORE_THAN_FIVE_FAILURE = 302
)

// 身份状态码
const (
	INVALID        = 0
	PENDING_REVIEW = 1
	VALID          = 2
	WITHOUT_CERT   = 3
)

// 日志级别码
const (
	NORMAL  = 1
	WARNING = 2
)

// 内外网日志码
const (
	outer = 0
	Inner = 1
)
