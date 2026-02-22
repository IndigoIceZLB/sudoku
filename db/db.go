package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

// Score 结构体用于映射数据库表
type Score struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Difficulty string `json:"difficulty"`
	TimeSpent  int    `json:"time_spent"` // 单位：秒
	CreatedAt  string `json:"created_at"`
}

// InitDB 初始化数据库连接并创建表
func InitDB() {
	// 从环境变量读取连接信息
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	if url == "" || token == "" {
		log.Fatal("Database URL or Token not found in environment variables")
	}

	// 拼接完整的连接字符串
	// 格式: libsql://dbname.turso.io?authToken=xxxx
	connStr := fmt.Sprintf("%s?authToken=%s", url, token)

	var err error
	DB, err = sql.Open("libsql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to Turso database successfully!")

	// 创建排行榜表 (如果不存在)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		time_spent INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table 'scores' is ready.")
}

// SaveScore 保存分数
func SaveScore(username, difficulty string, timeSpent int) error {
	query := `INSERT INTO scores (username, difficulty, time_spent) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, username, difficulty, timeSpent)
	return err
}

// GetTopScores 获取指定难度的前 10 名
func GetTopScores(difficulty string) ([]Score, error) {
	query := `
	SELECT id, username, difficulty, time_spent, created_at 
	FROM scores 
	WHERE difficulty = ? 
	ORDER BY time_spent ASC 
	LIMIT 10
	`
	rows, err := DB.Query(query, difficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var s Score
		if err := rows.Scan(&s.ID, &s.Username, &s.Difficulty, &s.TimeSpent, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}
