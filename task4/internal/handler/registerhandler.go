// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"go-study/task4/internal/utils"
	"net/http"
	"time"

	"go-study/task4/internal/logic"
	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/logx"
)

func registerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析注册请求参数失败: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 记录请求开始
		utils.LogAPIStart(r.Context(), "注册", req)
		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "注册", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, &req, &resp, err, duration, r)
	}
}
