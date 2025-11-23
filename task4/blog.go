package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"go-study/task4/internal/config"
	"go-study/task4/internal/handler"
	"go-study/task4/internal/svc"
	"go-study/task4/internal/utils"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "D:\\go-study\\go-study\\task4\\etc\\blog.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化日志系统
	setupLogger(&c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 记录系统启动信息
	utils.LogSystem(nil, "startup", "博客系统启动成功",
		"port", c.Port,
		"host", c.Host,
	)

	fmt.Printf("🚀 服务器启动成功！\n")
	fmt.Printf("📍 访问地址: http://localhost:%d\n", c.Port)
	fmt.Printf("⏰ 启动时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	server.Start()
}

// setupLogger 配置日志系统
func setupLogger(c *config.Config) {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
	}
	// 测试日志输出
	logx.Infof("=== 日志系统初始化 ===")
	logx.Infof("日志目录: %s", getLogPath(c))
	logx.Infof("日志级别: %s", c.Log.Level)
	logx.Infof("服务名称: %s", c.Log.ServiceName)
	logx.Infof("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	logx.Info("日志系统初始化完成")
}

// getLogPath 获取日志路径
func getLogPath(c *config.Config) string {
	if c.Log.Path != "" {
		return c.Log.Path
	}
	return "logs"
}
