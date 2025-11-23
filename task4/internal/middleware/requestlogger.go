package middleware

import (
	"go-study/task4/internal/utils"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// RequestLogger 请求日志中间件
type RequestLogger struct{}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// responseWriter 用于捕获响应状态码

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Handle 处理请求日志
func (m *RequestLogger) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建自定义的 ResponseWriter
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 记录请求开始
		logx.Infof("请求开始: %s %s - 客户端: %s",
			r.Method, r.URL.Path, getClientIP(r))

		// 处理请求
		next(rw, r)

		// 计算处理时间
		duration := time.Since(start)

		// 记录请求完成
		ctx := r.Context()
		utils.RequestInfo(ctx, r.Method, r.URL.Path, rw.statusCode, duration, getClientIP(r))
	}
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
