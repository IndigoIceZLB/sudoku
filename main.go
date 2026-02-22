package main

import (
	"net/http"
	"os"

	// 引入我们刚才写的 sudoku 包
	// 注意：这里的路径 "sudoku-backend/sudoku" 必须和你 go.mod 里的 module 名字一致
	// 如果你 go.mod 第一行是 "module sudoku-backend"，这里就是 "sudoku-backend/sudoku"
	"github.com/IndigoIceZLB/sudoku-backend/sudoku"

	"github.com/gin-gonic/gin"

	"github.com/IndigoIceZLB/sudoku-backend/db"
)

// 定义接收前端提交数据的结构
type ScoreRequest struct {
	Username   string `json:"username" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required"`
	TimeSpent  int    `json:"time_spent" binding:"required"`
}

func main() {
	// 1. 初始化数据库
	// 注意：在本地运行时，如果没设置环境变量，这里会报错退出
	// 建议先在本地设置环境变量，或者部署后再测
	db.InitDB()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 允许 POST
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // 允许 JSON Header
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Sudoku API with Database is Ready!"})
	})

	// 原有的生成游戏接口
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

	// --- 新增接口 ---

	// 1. 提交分数
	r.POST("/api/submit-score", func(c *gin.Context) {
		var req ScoreRequest
		// 绑定 JSON 数据
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 保存到数据库
		if err := db.SaveScore(req.Username, req.Difficulty, req.TimeSpent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Score saved!"})
	})

	// 2. 获取排行榜
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
