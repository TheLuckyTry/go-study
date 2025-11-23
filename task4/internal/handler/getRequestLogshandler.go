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

func getRequestLogshandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req types.RequestLogQueryRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析获取日志请求参数失败: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//设置默认值
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}
		l := logic.NewGetRequestLogslogic(r.Context(), svcCtx)
		utils.LogAPIStart(r.Context(), "获取日志列表", nil)
		resp, err := l.GetRequestLogs(&req)
		duration := time.Since(start)
		utils.LogAPIEnd(r.Context(), "获取文件列表", resp, err, duration)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		utils.LogRequest(r.Context(), svcCtx.DB, &req, &resp, err, duration, r)
	}
}
