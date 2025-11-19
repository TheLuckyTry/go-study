// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
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
	userId := l.ctx.Value("userId").(int64)

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
		UserID:    post.UserID,
		Username:  user.Username,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil

}
