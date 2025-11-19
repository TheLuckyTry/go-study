// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"go-study/task4/model"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	// 从JWT token中获取用户ID
	userId := l.ctx.Value("userId").(int64)

	var user model.User
	result := l.svcCtx.DB.First(&user, uint(userId))
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	return &types.UserInfoResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
