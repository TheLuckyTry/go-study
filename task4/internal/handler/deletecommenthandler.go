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

func deleteCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req types.DeleteCommentRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析删除评论参数失败: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		utils.LogAPIStart(r.Context(), "删除评论", req)
		l := logic.NewDeleteCommentLogic(r.Context(), svcCtx)
		resp, err := l.DeleteComment(&req)
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "删除评论", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, &req, &resp, err, duration, r)
	}
}
