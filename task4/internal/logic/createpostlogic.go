// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostRequest) (resp *types.PostResponse, err error) {
	// 从JWT token中获取用户ID
	// 从JWT token中获取用户ID - 修复版本
	userId, ok := utils.GetUserIDFromContext(l.ctx)
	if !ok {
		l.Logger.Error("无法从上下文中获取用户ID")
		return nil, nil
	}

	// 获取用户名
	username, _ := utils.GetUsernameFromContext(l.ctx)

	if username == "" {
		// 如果上下文中没有用户名，从数据库查询
		var user model.User
		result := l.svcCtx.DB.First(&user, userId)
		if result.Error != nil {
			l.Logger.Errorf("查询用户信息失败, 用户ID: %d, 错误: %v", userId, result.Error)
			return nil, nil
		}
		username = user.Username
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  uint(userId),
	}

	result := l.svcCtx.DB.Create(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询用户信息
	var user model.User
	l.svcCtx.DB.First(&user, userId)

	return &types.PostResponse{
		Id:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    userId,
		Username:  username,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}
