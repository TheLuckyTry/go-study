// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"go-study/task4/internal/types"
	"go-study/task4/internal/utils"
	"net/http"
	"time"

	"go-study/task4/internal/logic"
	"go-study/task4/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req types.GetPostRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析获取文章详情参数失败: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		utils.LogAPIStart(r.Context(), "获取文章", req)
		l := logic.NewGetPostLogic(r.Context(), svcCtx)
		resp, err := l.GetPost(&req) // 这里传递 &req 参数
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "获取文章", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, &req, &resp, err, duration, r)
	}
}
