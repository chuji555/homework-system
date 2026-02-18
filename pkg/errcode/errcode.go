package errcode

type ErrCode int

// 定义常用错误码
const (
	Success          ErrCode = 0
	ParamError       ErrCode = 10001
	AuthError        ErrCode = 10002
	PermissionDenied ErrCode = 10003
	DataNotFound     ErrCode = 10004
	DBError          ErrCode = 10005
	TokenExpired     ErrCode = 10006
)

// 获取错误信息
func (e ErrCode) Msg() string {
	switch e {
	case Success:
		return "success"
	case ParamError:
		return "参数错误"
	case AuthError:
		return "认证失败"
	case PermissionDenied:
		return "权限不足"
	case DataNotFound:
		return "数据不存在"
	case DBError:
		return "数据库操作失败"
	case TokenExpired:
		return "Token已过期"
	default:
		return "未知错误"
	}
}
