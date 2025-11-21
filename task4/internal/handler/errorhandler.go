package handler

import (
	"go-study/task4/internal/types"
	"go-study/task4/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

// CustomErrorHandler 自定义错误处理器 - 正确的函数签名
func CustomErrorHandler(err error) (int, interface{}) {
	logx.Errorf("请求错误: %v", err)

	var response types.BaseResponse

	// 根据错误类型返回不同的错误码
	switch e := err.(type) {
	case *utils.Error:
		// 处理自定义错误
		response = types.BaseResponse{
			Code:    e.Code,
			Message: e.Message,
		}
	default:
		// 处理其他错误
		response = types.BaseResponse{
			Code:    utils.CodeInternalError,
			Message: "服务器内部错误",
		}
	}

	// 返回 HTTP 状态码和响应体
	return getHttpStatus(response.Code), response
}

// 根据错误码获取 HTTP 状态码
func getHttpStatus(code int) int {
	switch code {
	case utils.CodeInvalidParams:
		return 400
	case utils.CodeUnauthorized:
		return 401
	case utils.CodeForbidden:
		return 403
	case utils.CodeNotFound:
		return 404
	default:
		return 500
	}
}
