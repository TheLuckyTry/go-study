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

type GetPostsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostsLogic {
	return &GetPostsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostsLogic) GetPosts() (resp *types.PostListResponse, err error) {
	var posts []model.Post
	result := l.svcCtx.DB.Preload("User").Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		l.Logger.Errorf("获取文章列表失败: %v", result.Error)
		return nil, utils.ErrorDatabaseError.WithMessage("获取文章列表失败")
	}
	var postResponses []types.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, types.PostResponse{
			Id:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			Username:  post.User.Username,
			CreatedAt: post.CreatedAt.Format(time.DateTime),
			UpdatedAt: post.UpdatedAt.Format(time.DateTime),
		})
	}
	l.Logger.Infof("成功获取文章列表，共 %d 篇文章", len(postResponses))

	return &types.PostListResponse{
		Code:    utils.CodeSuccess,
		Message: "成功获取文章",
		Data:    postResponses,
	}, nil
}
