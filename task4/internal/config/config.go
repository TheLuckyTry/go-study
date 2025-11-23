// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// 自定义配置项使用不同的字段名
	AppLog struct {
		ServiceName string
		Mode        string
		Level       string
		Path        string
		Encoding    string
		TimeFormat  string
		KeepDays    int
		MaxBackups  int
		MaxSize     int
	} `json:",optional"`
	Database struct {
		Driver string
		Source string
	}
	JWT struct {
		AccessSecret string
		AccessExpire int64
	}
}
