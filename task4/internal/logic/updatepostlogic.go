// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"go-study/task4/internal/utils"
	"go-study/task4/model"
	"strconv"
	"time"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePostLogic) UpdatePost(req *types.UpdatePostRequest) (resp *types.SinglePostResponse, err error) {
	// 从JWT token中获取用户ID
	userId, ok := utils.GetUserIDFromContext(l.ctx)
	if !ok {
		return nil, utils.ErrorUnauthorized
	}

	postId, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		l.Logger.Errorf("文章ID格式错误: %s, 错误: %v", req.Id, err)
		return nil, utils.ErrorInvalidParams.WithMessage("文章ID格式错误")
	}
	// 检查文章是否存在且属于当前用户
	var post model.Post
	result := l.svcCtx.DB.First(&post, postId)
	if result.Error != nil {
		l.Logger.Errorf("文章不存在, ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorPostNotFound
	}
	if post.UserID != userId {
		l.Logger.Errorf("用户 %d 尝试更新不属于自己的文章 %d", userId, postId)
		return nil, utils.ErrorNotAuthor
	}

	// 更新文章
	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}
	result = l.svcCtx.DB.Model(&post).Updates(updates)
	if result.Error != nil {
		l.Logger.Errorf("更新文章失败, ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorDatabaseError.WithMessage("更新文章失败")
	}

	// 重新查询文章信息（包含用户信息）
	l.svcCtx.DB.Preload("User").First(&post, uint(postId))

	postResponse := types.PostResponse{
		Id:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		Username:  post.User.Username,
		CreatedAt: post.CreatedAt.Format(time.DateTime),
		UpdatedAt: post.UpdatedAt.Format(time.DateTime),
	}

	l.Logger.Infof("用户 %d 成功更新文章 %d", userId, postId)

	return &types.SinglePostResponse{
		Code:    utils.CodeSuccess,
		Message: "更新文章成功",
		Data:    postResponse,
	}, nil
}
