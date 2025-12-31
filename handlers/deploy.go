package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

// 执行部署脚本的函数
func DeployService(scriptPath, serviceName, workingDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查脚本文件是否存在
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Script file %s not found", scriptPath)})
			return
		}

		// 检查工作目录是否存在
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Working directory %s not found", workingDir)})
			return
		}

		// 执行部署脚本
		cmd := exec.Command("/bin/bash", scriptPath, serviceName)
		cmd.Dir = workingDir // 设置脚本的工作目录

		// 捕获执行结果
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   fmt.Sprintf("Deployment failed: %v", err),
				"details": string(cmdOutput),
			})
			return
		}

		// 返回执行结果
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Deployment of %s completed successfully", serviceName),
			"output":  string(cmdOutput),
		})
	}
}
