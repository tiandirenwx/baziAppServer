package errorcode

// 客户端相关错误码
var (
	ParamsErr = NewClientError(-40001, "参数错误", nil)
)
