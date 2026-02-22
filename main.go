package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// 引入我们刚才写的 sudoku 包
	// 注意：这里的路径 "sudoku-backend/sudoku" 必须和你 go.mod 里的 module 名字一致
	// 如果你 go.mod 第一行是 "module sudoku-backend"，这里就是 "sudoku-backend/sudoku"
	"github.com/IndigoIceZLB/sudoku-backend/sudoku"

	"github.com/gin-gonic/gin"

	"github.com/IndigoIceZLB/sudoku-backend/db"

	"github.com/gin-contrib/cors" // 引入官方 CORS 包
)

// 定义接收前端提交数据的结构
type ScoreRequest struct {
	Username   string `json:"username" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required"`
	TimeSpent  int    `json:"time_spent" binding:"required"`
}

func main() {
	// 初始化数据库
	db.InitDB()

	r := gin.Default()

	// 🛑 核心修复：使用官方 CORS 中间件配置
	// 这能解决 99% 的 "提交失败" 问题
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源（生产环境可以改成你的前端域名）
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Sudoku API is Ready!"})
	})

	r.GET("/api/new-game", func(c *gin.Context) {
		levelStr := c.Query("level")
		holes := 30
		switch levelStr {
		case "medium":
			holes = 40
		case "hard":
			holes = 50
		case "expert":
			holes = 55
		}
		puzzle, solution := sudoku.Generate(holes)
		c.JSON(http.StatusOK, gin.H{
			"difficulty": levelStr,
			"holes":      holes,
			"puzzle":     puzzle,
			"solution":   solution,
		})
	})

	r.POST("/api/submit-score", func(c *gin.Context) {
		var req ScoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println("Bind Error:", err) // 打印日志到 Render 控制台
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
			return
		}

		fmt.Printf("Receiving score: %+v\n", req) // 打印接收到的数据

		if err := db.SaveScore(req.Username, req.Difficulty, req.TimeSpent); err != nil {
			fmt.Println("DB Error:", err) // 打印数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Score saved!"})
	})

	r.GET("/api/leaderboard", func(c *gin.Context) {
		difficulty := c.Query("difficulty")
		if difficulty == "" {
			difficulty = "easy"
		}
		scores, err := db.GetTopScores(difficulty)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"leaderboard": scores})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}
