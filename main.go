package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化 Gin 引擎
	r := gin.Default()

	// 2. 定义一个简单的路由 (用于测试服务是否存活)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello! Sudoku Backend is running correctly on Render.",
			"status":  "success",
		})
	})

	// 3. 获取环境变量 PORT (Render 会自动注入这个变量)
	// 如果没有设置（比如在本地运行），默认为 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 4. 监听 0.0.0.0 (这是 Render 要求的，不能只监听 localhost)
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Server is starting on %s...\n", addr)

	// 启动服务
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
