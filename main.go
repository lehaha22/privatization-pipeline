package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hainancicd/handlers"
	"hainancicd/middlewares"
	"hainancicd/utils"
	"log"
	"net/http"
	"sync"
)

var (
	// Using a mutex to ensure thread-safe access to the config
	configMutex sync.RWMutex
	config      *utils.Config
)

func main() {
	// 加载配置文件
	config, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	maxSize := int64(config.Upload.MaxSize) // 文件大小限制 (MB)

	// 文件上传限制
	router.MaxMultipartMemory = maxSize

	// 应用 AuthMiddleware 到所有需要的路由
	router.Use(middlewares.AuthMiddleware()) // 在这里应用中间件，保证所有请求都需要 Token 校验

	// 定义上传路由
	router.POST("/upload", handlers.UploadFile(config.Upload.Directory, maxSize))

	// 定义后端部署接口
	router.POST("/deploy_backend", func(c *gin.Context) {
		// 从请求中获取参数
		var request struct {
			ServiceName string `json:"service_name"`
		}

		// 绑定请求体
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		handlers.DeployService(config.Deploy.Backend.ScriptPath, request.ServiceName, config.Deploy.Backend.WorkingDir)(c)
	})

	// 定义前端部署接口
	router.POST("/deploy_frontend", func(c *gin.Context) {
		// 从请求中获取参数
		var request struct {
			ServiceName string `json:"service_name"`
		}

		// 绑定请求体
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		handlers.DeployService(config.Deploy.Frontend.ScriptPath, request.ServiceName, config.Deploy.Frontend.WorkingDir)(c)
	})

	log.Printf("Server is starting on port %d...", config.Server.Port)

	// 启动服务器
	port := config.Server.Port
	router.Run(fmt.Sprintf(":%d", port))

}
