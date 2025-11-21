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

type GetPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostLogic) GetPost(req *types.GetPostRequest) (resp *types.SinglePostResponse, err error) {
	postId, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		l.Logger.Errorf("文章ID格式错误: %s, 错误: %v", req.Id, err)
		return nil, utils.ErrorInvalidParams.WithMessage("文章ID格式错误")
	}

	var post model.Post
	result := l.svcCtx.DB.Preload("User").First(&post, uint(postId))
	if result.Error != nil {
		l.Logger.Errorf("获取文章失败, ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorPostNotFound
	}

	postResponse := types.PostResponse{
		Id:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		Username:  post.User.Username,
		CreatedAt: post.CreatedAt.Format(time.DateTime),
		UpdatedAt: post.UpdatedAt.Format(time.DateTime),
	}
	l.Logger.Infof("成功获取文章: %s (ID: %d)", post.Title, post.ID)
	return &types.SinglePostResponse{
		Code:    utils.CodeSuccess,
		Message: "成功获取文章",
		Data:    postResponse,
	}, nil
}
