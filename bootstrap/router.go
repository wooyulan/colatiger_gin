package bootstrap

import (
	"colatiger/app/middleware"
	"colatiger/global"
	"colatiger/routers"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupRouter() *gin.Engine {

	if global.App.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), middleware.CustomRecovery())

	// 跨域处理
	router.Use(middleware.Cors())

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routers.SetApiGroupRoutes(apiGroup)

	return router

}

// RunServer 启动服务器
func RunServer() {
	r := setupRouter()
	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
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
