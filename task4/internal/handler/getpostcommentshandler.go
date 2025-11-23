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

func getPostCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListRequest
		start := time.Now()
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析获取文章列表请求参数失败: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		utils.LogAPIStart(r.Context(), "获取文章列表", req)
		l := logic.NewGetPostCommentsLogic(r.Context(), svcCtx)
		resp, err := l.GetPostComments(&req)
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "获取文章列表", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, &req, &resp, err, duration, r)
	}
}
