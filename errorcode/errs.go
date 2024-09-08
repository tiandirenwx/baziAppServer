package errorcode

// Package errorcode 错误码包

import (
	"fmt"
	"regexp"
	"strconv"

	"trpc.group/trpc-go/trpc-go/errs"
)

var defaultErrorCode = -10000

// baseError 通用基础Error, 包含错误码,Message,上次的Error
type baseError struct {
	Code    int
	Message string
	Err     error // 原始错误
}

// Error baseError实现Error接口
func (e baseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code:%d, message:%s, cause: \n\t%v", e.Code, e.Message, e.Err)
	} else {
		return fmt.Sprintf("code:%d, message:%s, cause: nil", e.Code, e.Message)
	}
}

var errRgx = regexp.MustCompile(`code:(-\d+), message:(.*?), cause: .*?`)

// FromText 根据异常文本生成错误, 比如事件消费过程中的错误存储在redis中
func FromText(desc string) error {
	if len(desc) < 5 { // code: 前缀都没有
		return nil
	}

	subs := errRgx.FindStringSubmatch(desc)
	if len(subs) != 3 {
		return nil
	}

	sCode, message := subs[1], subs[2]
	code, _ := strconv.Atoi(sCode)
	return baseError{Code: code, Message: message}
}

// WithMessage 修改error的Message,使之更具体明确
func WithMessage(e error, msg string) error {
	ec := getErrorCode(e)
	// return baseError{Code: ec, Message: msg}
	return errs.New(ec, msg)
}

// IsBaseError 是否是 baseError
func IsBaseError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(baseError)
	if ok {
		return true
	}

	_, ok = err.(*errs.Error)
	return ok
}

// IsBizError 是否是业务层抛出的异常，包含ClientError,ServerError,SystermError
func IsBizError(err error) bool {
	if err == nil {
		return false
	}

	code := getErrorCode(err)
	return code < 0 && code != defaultErrorCode
}

// IsClientError 是否是 clientError
func IsClientError(err error) bool {
	if err != nil {
		_, ok := err.(ClientError)
		return ok
	}

	return false
}

// IsSameBizErr 判断是否是相同类型的错误
func IsSameBizErr(err error, code int) bool {
	return getErrorCode(err) == code
}

func getErrorCode(err error) int {
	if err == nil {
		return 0
	}

	var errCode = defaultErrorCode
	switch e := err.(type) {
	case ClientError:
		errCode = e.Code
	case ServerError:
		errCode = e.Code
	case SystemError:
		errCode = e.Code
	case baseError:
		errCode = e.Code
	case *errs.Error:
		errCode = int(e.Code)
	}
	return errCode
}

// WrapError 一般用于异常转换/包装
func WrapError(base error, err error) error {
	switch e := base.(type) {
	case ClientError:
		return WrapClientError(e, err)
	case ServerError:
		return WrapServerError(e, err)
	case SystemError:
		return WrapSystemError(e, err)
	default:
		return ServerError{baseError{Code: defaultErrorCode, Message: base.Error(), Err: err}}
	}
}

// ClientError 客户端错误: 由于客户端输入或者不满足业务条件导致的错误
// 1.参数错误类, 邮箱格式错误, 账号不存在
// 2.安全保护: 请求过于频繁,请求参数过长
// 无需告警, 正常业务拦截
type ClientError struct {
	baseError
}

// NewClientError 一般用于异常转换/包装
func NewClientError(code int, msg string, err error) ClientError {
	return ClientError{baseError{Code: code, Message: msg, Err: err}}
}

// WrapClientError 一般用于异常转换/包装
func WrapClientError(ce ClientError, err error) ClientError {
	return ClientError{baseError{Code: ce.Code, Message: ce.Message, Err: err}}
}

// ServerError 服务端错误: 客户端输入正确,条件满足,但业务还处理失败
// 1.是本业务代码问题,如未考虑的边界情况: 请求获取全局锁失败; 由于并发导致乐观锁失败,CAS; 由于服务过载而丢弃请求
// 2.违反产品规则: 国外用户登录使用国内版本
// 3.流程控制,通过返回特定的异常,调用函数进行相应的处理
// 需要告警,一般较少发生
type ServerError struct {
	baseError
}

// NewServerError 一般用于异常转换/包装
func NewServerError(code int, msg string, err error) ServerError {
	return ServerError{baseError{Code: code, Message: msg, Err: err}}
}

// WrapServerError 一般用于异常转换/包装
func WrapServerError(se ServerError, err error) ServerError {
	return ServerError{baseError{Code: se.Code, Message: se.Message, Err: err}}
}

// SystemError 系统错误: 业务组件失败,跟本业务代码无关,需要告警,紧急处理
// 如:MySQL超时, Redis超时, 下游RPC超时, 外部服务/依赖不可用,返回预期之外的错误(证书过期,密钥失效,IP白名单)
type SystemError struct {
	baseError
}

// NewSystemError 一般用于异常转换/包装
func NewSystemError(code int, msg string, err error) SystemError {
	return SystemError{baseError{Code: code, Message: msg, Err: err}}
}

// WrapSystemError 一般用于异常转换/包装
func WrapSystemError(se SystemError, err error) SystemError {
	return SystemError{baseError{Code: se.Code, Message: se.Message, Err: err}}
}
