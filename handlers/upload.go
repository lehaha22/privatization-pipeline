package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(uploadDir string, maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		// 检查文件大小
		if file.Size > maxSize*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File is too large. Max size is %d MB", maxSize)})
			return
		}

		// 确保上传目录存在
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
				return
			}
		}

		// 保存文件
		filePath := filepath.Join(uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filePath": filePath})
	}
}
