// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"go-study/task4/internal/types"
	"net/http"

	"go-study/task4/internal/logic"
	"go-study/task4/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPostRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetPostLogic(r.Context(), svcCtx)
		resp, err := l.GetPost(&req) // 这里传递 &req 参数
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
