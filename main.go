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
	// 1. 初始化数据库 (即使失败也继续启动 Web 服务，但在日志里报错)
	if err := db.InitDB(); err != nil {
		fmt.Printf("⚠️⚠️⚠️ DATABASE ERROR: %v\n", err)
		fmt.Println("Server will start, but database features will fail.")
	}

	r := gin.Default()

	// 2. 配置 CORS (允许所有跨域请求，这是最宽松的配置)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Sudoku API is running"})
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
		// 检查数据库是否就绪
		if db.DB == nil {
			fmt.Println("Error: Database not connected")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected. Check server logs."})
			return
		}

		var req ScoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.SaveScore(req.Username, req.Difficulty, req.TimeSpent); err != nil {
			fmt.Printf("Save Score Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Score saved!"})
	})

	r.GET("/api/leaderboard", func(c *gin.Context) {
		if db.DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
			return
		}
		difficulty := c.Query("difficulty")
		if difficulty == "" {
			difficulty = "easy"
		}

		scores, err := db.GetTopScores(difficulty)
		if err != nil {
			fmt.Printf("Leaderboard Error: %v\n", err)
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
