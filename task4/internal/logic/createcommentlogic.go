// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"go-study/task4/internal/utils"
	"go-study/task4/model"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentRequest) (resp *types.BaseResponse, err error) {
	// 从JWT token中获取用户ID
	userId, ok := utils.GetUserIDFromContext(l.ctx)
	if !ok {
		l.Logger.Error("无法从JWT token中获取用户ID")
		return nil, utils.ErrorUnauthorized
	}
	// 获取用户名
	username, _ := utils.GetUsernameFromContext(l.ctx)
	if username == "" {
		// 从数据库查询用户名
		var user model.User
		if result := l.svcCtx.DB.First(&user, userId); result.Error != nil {
			l.Logger.Errorf("查询用户信息失败: %v", result.Error)
			return nil, utils.ErrorUserNotFound
		}
		username = user.Username
	}

	// 检查文章是否存在
	var post model.Post
	result := l.svcCtx.DB.First(&post, req.PostID)
	if result.Error != nil {
		l.Logger.Errorf("文章不存在, ID: %d, 错误: %v", req.PostID, result.Error)
		return nil, utils.ErrorPostNotFound
	}

	// 创建评论
	comment := model.Comment{
		Content: req.Content,
		UserID:  userId,
		PostID:  req.PostID,
	}

	if result := l.svcCtx.DB.Create(&comment); result.Error != nil {
		l.Logger.Errorf("创建评论失败, 用户ID: %d, 文章ID: %d, 错误: %v", userId, req.PostID, result.Error)
		return nil, utils.ErrorDatabaseError.WithMessage("创建评论失败")
	}

	l.Logger.Infof("用户 %s (ID: %d) 成功为文章 %d 创建评论 (ID: %d)", username, userId, req.PostID, comment.ID)

	return &types.BaseResponse{
		Code:    utils.CodeSuccess,
		Message: "评论创建成功",
	}, nil
}
