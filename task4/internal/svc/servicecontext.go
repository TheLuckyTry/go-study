// Package svc Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
/**
服务上下文管理：ServiceContext 是 go-zero 框架中的核心服务上下文结构，用于管理应用程序的全局依赖和服务实例
依赖注入容器：作为依赖注入容器，存储和提供应用程序运行所需的各种服务组件，如数据库连接、配置信息等
*/
package svc

import (
	"go-study/task4/internal/config"
	"go-study/task4/model"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitDB()
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
