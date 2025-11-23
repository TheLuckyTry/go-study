// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"go-study/task4/internal/utils"
	"net/http"
	"time"

	"go-study/task4/internal/logic"
	"go-study/task4/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPostsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := logic.NewGetPostsLogic(r.Context(), svcCtx)
		utils.LogAPIStart(r.Context(), "获取文章列表", nil)
		resp, err := l.GetPosts()
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "获取文章列表", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, nil, &resp, err, duration, r)
	}
}
