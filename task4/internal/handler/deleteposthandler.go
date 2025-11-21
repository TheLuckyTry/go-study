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

func deletePostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.DeletePostRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeletePostLogic(r.Context(), svcCtx)
		resp, err := l.DeletePost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
