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

type GetPostCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostCommentsLogic {
	return &GetPostCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostCommentsLogic) GetPostComments(req *types.CommentListRequest) (resp *types.CommentListResponse, err error) {
	//从请求中获取文章ID
	postId, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		l.Logger.Errorf("文章ID格式错误: %s, 错误: %v", req.Id, err)
		return nil, utils.ErrorInvalidParams.WithMessage("文章ID格式错误")
	}
	// 检查文章是否存在
	var post model.Post
	result := l.svcCtx.DB.First(&post, uint(postId))
	if result.Error != nil {
		l.Logger.Errorf("文章不存在, ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorPostNotFound
	}

	// 获取该文章的所有评论
	var comments []model.Comment
	result = l.svcCtx.DB.Preload("User").Where("post_id = ?", postId).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		l.Logger.Errorf("获取评论列表失败, 文章ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorDatabaseError.WithMessage("获取评论列表失败")
	}

	// 构建响应数据
	var commentResponses []types.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, types.CommentResponse{
			Id:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			Username:  comment.User.Username,
			PostID:    comment.PostID,
			CreatedAt: comment.CreatedAt.Format(time.DateTime),
		})
	}

	l.Logger.Infof("成功获取文章 %d 的评论列表，共 %d 条评论", postId, len(commentResponses))

	return &types.CommentListResponse{
		Code:    utils.CodeSuccess,
		Message: "获取评论列表成功",
		Data:    commentResponses,
	}, nil

}
