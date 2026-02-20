package main

import (
	"net/http"
	"os"

	// 引入我们刚才写的 sudoku 包
	// 注意：这里的路径 "sudoku-backend/sudoku" 必须和你 go.mod 里的 module 名字一致
	// 如果你 go.mod 第一行是 "module sudoku-backend"，这里就是 "sudoku-backend/sudoku"
	"github.com/IndigoIceZLB/sudoku-backend/sudoku"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 允许跨域请求 (CORS) - 这对后续前端开发很重要
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Sudoku API is Ready!"})
	})

	// 新增接口：生成数独
	// 访问示例：/api/new-game?level=easy
	r.GET("/api/new-game", func(c *gin.Context) {
		// 获取难度参数，默认为 easy
		levelStr := c.Query("level")
		holes := 30 // 默认挖 30 个洞 (简单)

		switch levelStr {
		case "medium":
			holes = 40
		case "hard":
			holes = 50
		case "expert":
			holes = 55
		}

		// 调用核心算法
		puzzle, solution := sudoku.Generate(holes)

		// 返回 JSON
		c.JSON(http.StatusOK, gin.H{
			"difficulty": levelStr,
			"holes":      holes,
			"puzzle":     puzzle,   // 题目 (0 代表空格)
			"solution":   solution, // 答案 (前端可以用来验证，或者暂时不发给前端防作弊)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run("0.0.0.0:" + port)
}
