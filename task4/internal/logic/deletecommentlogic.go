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

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentRequest) (resp *types.BaseResponse, err error) {
	// 从JWT token中获取用户ID
	userId, ok := utils.GetUserIDFromContext(l.ctx)
	if !ok {
		return nil, utils.ErrorUnauthorized
	}
	commentId, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		l.Logger.Errorf("文章ID格式错误: %s, 错误: %v", req.Id, err)
		return nil, utils.ErrorInvalidParams.WithMessage("评论ID格式错误")
	}

	// 检查评论是否存在且属于当前用户
	var comment model.Comment
	result := l.svcCtx.DB.First(&comment, uint(commentId))
	if result.Error != nil {
		l.Logger.Errorf("评论不存在, ID: %d, 错误: %v", commentId, result.Error)
		return nil, utils.ErrorCommentNotFound
	}

	if comment.UserID != userId {
		return nil, utils.ErrorForbidden.WithMessage("无权删除此评论")
	}

	// 删除评论
	result = l.svcCtx.DB.Delete(&comment)
	if result.Error != nil {
		l.Logger.Errorf("删除评论失败, ID: %d, 错误: %v", commentId, result.Error)
		return nil, utils.ErrorDatabaseError.WithMessage("删除评论失败")
	}

	l.Logger.Infof("用户 %d 成功删除评论 %d", userId, commentId)
	return &types.BaseResponse{
		Code:    utils.CodeSuccess,
		Message: "删除成功",
	}, nil
}
