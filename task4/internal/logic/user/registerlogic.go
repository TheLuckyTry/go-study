// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"go-study/task4/internal/utils"
	"go-study/task4/model"

	"go-study/task4/internal/svc"
	"go-study/task4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// 检查用户名是否已存在
	var existingUser model.User
	result := l.svcCtx.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	result = l.svcCtx.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("邮箱已存在")
	}
	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}
	result = l.svcCtx.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &types.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
