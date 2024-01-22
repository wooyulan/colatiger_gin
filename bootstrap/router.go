package bootstrap

import (
	"colatiger/app/middleware"
	"colatiger/global"
	"colatiger/routers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// 初始化路由
	router := gin.Default()
	// 设置控制台日志级别
	gin.SetMode(global.App.Config.App.RunLogType)

	// 为 multipart forms 设置较低的内存限制
	router.MaxMultipartMemory = 8 << 20 // 8MB

	// 错误日志
	router.Use(gin.Logger(), middleware.CustomRecovery())

	// 跨域处理
	router.Use(middleware.Cors())
	// 限流
	// 验证接口是否合法
	// 验证token
	// 找不到路由

	// 绑定路由 TODO

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routers.SetApiGroupRoutes(apiGroup)

	return router

}

// RunServer 启动服务器
func RunServer() {
	// 加载路由
	r := setupRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	log.Println("server start success!启动端口：" + global.App.Config.App.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			str := fmt.Sprintf("listen: %s\n", err) //拼接字符串
			global.App.Log.Error(str)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
