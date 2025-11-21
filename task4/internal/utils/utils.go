package utils

import (
	"errors"
	"go-study/task4/internal/types"

	"golang.org/x/crypto/bcrypt"
)

// 错误码定义
const (
	CodeSuccess       = 200
	CodeInvalidParams = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
	CodeDatabaseError = 501
)

// 业务错误码
const (
	CodeUserExists      = 1001
	CodeUserNotFound    = 1002
	CodeWrongPassword   = 1003
	CodePostNotFound    = 2001
	CodeNotAuthor       = 2002
	CodeCommentNotFound = 3001
)

var (
	ErrorSuccess       = NewError(CodeSuccess, "成功")
	ErrorInvalidParams = NewError(CodeInvalidParams, "请求参数错误")
	ErrorUnauthorized  = NewError(CodeUnauthorized, "未授权")
	ErrorForbidden     = NewError(CodeForbidden, "禁止访问")
	ErrorNotFound      = NewError(CodeNotFound, "资源不存在")
	ErrorInternalError = NewError(CodeInternalError, "内部服务器错误")
	ErrorDatabaseError = NewError(CodeDatabaseError, "数据库错误")

	// 业务错误
	ErrorUserExists      = NewError(CodeUserExists, "用户已存在")
	ErrorUserNotFound    = NewError(CodeUserNotFound, "用户不存在")
	ErrorWrongPassword   = NewError(CodeWrongPassword, "用户名或密码错误")
	ErrorPostNotFound    = NewError(CodePostNotFound, "文章不存在")
	ErrorNotAuthor       = NewError(CodeNotAuthor, "无权操作此文章")
	ErrorCommentNotFound = NewError(CodeCommentNotFound, "评论不存在")
)

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithMessage(message string) *Error {
	return &Error{
		Code:    e.Code,
		Message: message,
	}
}

func SuccessBaseResponse(msg string) types.BaseResponse {
	if msg == "" {
		msg = "成功"
	}
	return types.BaseResponse{
		Code:    CodeSuccess,
		Message: msg,
	}
}
func ToBaseResponse(err error) types.BaseResponse {
	if err == nil {
		return types.BaseResponse{
			Code:    CodeSuccess,
			Message: "成功",
		}
	}

	var e *Error
	if errors.As(err, &e) {
		return types.BaseResponse{
			Code:    e.Code,
			Message: e.Message,
		}
	}

	return types.BaseResponse{
		Code:    CodeInternalError,
		Message: err.Error(),
	}
}

type Error struct {
	Code    int
	Message string
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
