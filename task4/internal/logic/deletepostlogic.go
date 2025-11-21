// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"go-study/task4/internal/utils"
	"go-study/task4/model"
	"strconv"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePostLogic) DeletePost(req *types.DeletePostRequest) (resp *types.BaseResponse, err error) {
	//从请求中获取文章ID
	postId, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		l.Logger.Errorf("文章ID格式错误: %s, 错误: %v", req.Id, err)
		return nil, utils.ErrorInvalidParams.WithMessage("文章ID格式错误")
	}
	// 从JWT token中获取用户ID
	userId, ok := utils.GetUserIDFromContext(l.ctx)
	if !ok {
		return nil, utils.ErrorUnauthorized
	}
	// 检查文章是否存在且属于当前用户
	var post model.Post
	result := l.svcCtx.DB.First(&post, uint(postId))
	if result.Error != nil {
		l.Logger.Errorf("文章不存在, ID: %d, 错误: %v", postId, result.Error)
		return nil, utils.ErrorPostNotFound
	}
	if post.UserID != userId {
		l.Logger.Errorf("用户 %d 尝试删除不属于自己的文章 %d", userId, postId)
		return nil, utils.ErrorNotAuthor
	}

	// 删除文章及相关评论
	tx := l.svcCtx.DB.Begin()

	// 先删除相关评论
	if err := tx.Where("post_id = ?", postId).Delete(&model.Comment{}).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除文章评论失败, 文章ID: %d, 错误: %v", postId, err)
		return nil, utils.ErrorDatabaseError.WithMessage("删除文章失败")
	}

	// 再删除文章
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("删除文章失败, ID: %d, 错误: %v", postId, err)
		return nil, utils.ErrorDatabaseError.WithMessage("删除文章失败")
	}

	tx.Commit()

	l.Logger.Infof("用户 %d 成功删除文章 %d", userId, postId)
	return &types.BaseResponse{
		Code:    utils.CodeSuccess,
		Message: "删除成功",
	}, nil
}
