package logic

import (
	"context"
	"go-study/task4/internal/utils"
	"go-study/task4/model"
	"time"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequestLogslogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRequestLogslogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequestLogslogic {
	return &GetRequestLogslogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRequestLogslogic) GetRequestLogs(req *types.RequestLogQueryRequest) (resp *types.RequestLogListResponse, err error) {
	// 创建数据库查询
	db := l.svcCtx.DB.Model(&model.RequestLog{})

	// 添加过滤条件（非必填）
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Username != "" {
		db = db.Where("username = ?", req.Username)
	}
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	}

	// 计算偏移量
	offset := (req.Page - 1) * req.PageSize

	// 查询总数
	var total int64
	if err = db.Count(&total).Error; err != nil {
		utils.LogError(l.ctx, err, "request_log", "count")
		return nil, err
	}

	// 查询数据
	var logs []model.RequestLog
	if err = db.Order("id DESC").Limit(req.PageSize).Offset(offset).Find(&logs).Error; err != nil {
		utils.LogError(l.ctx, err, "request_log", "query")
		return nil, err
	}

	// 转换为响应格式
	data := make([]types.RequestLogResponse, len(logs))
	for i, log := range logs {
		data[i] = types.RequestLogResponse{
			ID:         log.ID,
			Method:     log.Method,
			Path:       log.Path,
			StatusCode: log.StatusCode,
			UserID:     log.UserID,
			Username:   log.Username,
			IPAddress:  log.IPAddress,
			UserAgent:  log.UserAgent,
			Request:    log.Request,
			Response:   log.Response,
			Error:      log.Error,
			Duration:   log.Duration,
			CreatedAt:  log.CreatedAt.Format(time.DateTime),
		}
	}

	return &types.RequestLogListResponse{
		Code:     0,
		Message:  "success",
		Page:     req.Page,
		PageSize: req.PageSize,
		Data:     data,
	}, nil
}
