package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go-study/task4/model"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RequestLogData struct {
	Method     string
	Path       string
	StatusCode int
	UserID     uint
	Username   string
	IPAddress  string
	UserAgent  string
	Request    interface{}
	Response   interface{}
	Error      string
	Duration   time.Duration
}

// LogInfo 记录业务信息日志
func LogInfo(ctx context.Context, module, action, message string, data ...interface{}) {
	logx.WithContext(ctx).Info(formatLog(module, action, message, data...))
}

// LogError 记录错误日志
func LogError(ctx context.Context, err error, module, operation string, data ...interface{}) {
	logx.WithContext(ctx).Error(formatLog(module, operation, "操作失败: "+err.Error(), data...))
}

// LogWarn 记录警告日志
func LogWarn(ctx context.Context, module, action, message string, data ...interface{}) {
	logx.WithContext(ctx).Error(formatLog(module, action, message, data...))
}

// LogDatabase 记录数据库操作日志
func LogDatabase(ctx context.Context, operation, table string, duration time.Duration, rowsAffected int64) {
	message := fmt.Sprintf("数据库操作: %s %s 耗时: %v", operation, table, duration)
	if rowsAffected > 0 {
		message += fmt.Sprintf(", 影响行数: %d", rowsAffected)
	}

	// 慢查询警告
	if duration > 100*time.Millisecond {
		logx.WithContext(ctx).Error(formatLog("database", "slow", message))
	} else {
		logx.WithContext(ctx).Debug(formatLog("database", operation, message))
	}
}

// LogRequest 记录请求日志，支持入参和出参
func LogRequest(ctx context.Context, db *gorm.DB, req interface{}, resp interface{},
	err error, duration time.Duration, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	// 构建基础消息
	message := fmt.Sprintf("%s %s, 耗时: %v", method, path, duration)
	// 准备日志字段
	var fields []interface{}
	// 添加入参（如果存在）
	if req != nil {
		fields = append(fields, fmt.Sprintf("req: %v", req))
	}
	// 添加出参（如果存在）
	if resp != nil {
		fields = append(fields, fmt.Sprintf("resp: %v", resp))
	}
	// 添加错误信息（如果存在）
	if err != nil {
		fields = append(fields, fmt.Sprintf("error: %v", err))
	}
	// 格式化完整消息
	if len(fields) > 0 {
		message += fmt.Sprintf(" | %v", fields)
	}
	// 根据是否有错误决定日志级别
	if err != nil {
		logx.WithContext(ctx).Error(formatLog("request", "error", message))
	} else {
		logx.WithContext(ctx).Info(formatLog("request", "success", message))
	}
	// 记录到数据库
	logData := RequestLogData{
		Method:    method,
		Path:      path,
		Duration:  duration,
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
	}
	// 设置用户信息
	if userID, ok := GetUserIDFromContext(ctx); ok {
		logData.UserID = userID
	}
	if username, ok := GetUsernameFromContext(ctx); ok {
		logData.Username = username
	}

	// 设置请求和响应数据
	logData.Request = req
	logData.Response = resp
	if err != nil {
		logData.Error = err.Error()
	}
	//异步保存到日志
	go func() {
		// 保存到数据库
		LogRequestToDB(ctx, db, logData)
	}()
}

// LogSystem 记录系统日志
func LogSystem(ctx context.Context, operation, message string, data ...interface{}) {
	logx.WithContext(ctx).Info(formatLog("system", operation, message, data...))
}

// formatLog 格式化日志消息
func formatLog(module, action, message string, data ...interface{}) string {
	base := fmt.Sprintf("[%s.%s] %s", module, action, message)
	if len(data) > 0 {
		base += fmt.Sprintf(" | 数据: %v", data)
	}
	return base
}

func LogAPIStart(ctx context.Context, apiName string, req interface{}) {
	var reqJSON []byte
	var err error

	// 只有当 req 不为 nil 时才进行序列化
	if req != nil {
		reqJSON, err = json.Marshal(req)
		if err != nil {
			logx.WithContext(ctx).Errorf("序列化请求参数失败: %v", err)
			return
		}
	} else {
		// 当 req 为 nil 时，使用 null 字符串表示
		reqJSON = []byte("null")
	}

	logx.WithContext(ctx).Infow(fmt.Sprintf("%s 请求开始", apiName),
		logx.Field("request_json", string(reqJSON)),
	)
}

func LogAPIEnd(ctx context.Context, apiName string, resp interface{}, err error, duration time.Duration) {
	if err != nil {
		logx.WithContext(ctx).Errorw(fmt.Sprintf("%s 请求失败", apiName),
			logx.Field("error", err.Error()),
			logx.Field("duration", duration.String()),
		)
	} else {
		respJSON, jsonErr := json.Marshal(resp)
		if jsonErr != nil {
			logx.WithContext(ctx).Errorf("序列化响应结果失败: %v", jsonErr)
			return
		}

		logx.WithContext(ctx).Infow(fmt.Sprintf("%s 请求成功", apiName),
			logx.Field("response_json", string(respJSON)),
			logx.Field("duration", duration.String()),
		)
	}
}

func LogRequestToDB(ctx context.Context, db *gorm.DB, data RequestLogData) {
	// 限制请求体和响应体的大小，避免存储过大
	requestStr := ConvertToJSONString(data.Request, 1000)
	responseStr := ConvertToJSONString(data.Response, 2000)
	errorStr := limitStringSize(data.Error, 500)

	requestLog := model.RequestLog{
		Method:     data.Method,
		Path:       data.Path,
		StatusCode: data.StatusCode,
		UserID:     data.UserID,
		Username:   data.Username,
		IPAddress:  data.IPAddress,
		UserAgent:  data.UserAgent,
		Request:    requestStr,
		Response:   responseStr,
		Error:      errorStr,
		Duration:   int64(data.Duration / time.Millisecond), // 转换为毫秒
	}

	// 保存到数据库
	result := db.Create(&requestLog)
	if result.Error != nil {
		logx.Errorf("保存请求日志失败: %v", result.Error)
	} else {
		logx.Infof("请求日志保存成功 - ID: %d, 路径: %s", requestLog.ID, data.Path)
	}
}

// limitStringSize 限制字符串大小
func limitStringSize(str string, maxSize int) string {
	if len(str) > maxSize {
		return str[:maxSize] + "..."
	}
	return str
}

func ConvertToJSONString(data interface{}, maxSize int) string {
	if data == nil {
		return "null"
	}

	var dataStr string
	// 尝试JSON序列化
	if jsonData, err := json.Marshal(data); err == nil {
		dataStr = string(jsonData)
	} else {
		dataStr = fmt.Sprintf("%v", data)
	}
	if len(dataStr) > maxSize {
		return dataStr[:maxSize] + "..."
	}
	return dataStr
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
