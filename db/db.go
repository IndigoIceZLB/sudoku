package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

type Score struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Difficulty string `json:"difficulty"`
	TimeSpent  int    `json:"time_spent"`
	CreatedAt  string `json:"created_at"`
}

// InitDB 改为返回 error，而不是直接退出程序
func InitDB() error {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	if url == "" || token == "" {
		return fmt.Errorf("environment variables TURSO_DATABASE_URL or TURSO_AUTH_TOKEN are missing")
	}

	connStr := fmt.Sprintf("%s?authToken=%s", url, token)

	var err error
	DB, err = sql.Open("libsql", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}

	// 尝试 Ping 一下
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to turso: %v", err)
	}

	fmt.Println("Connected to Turso database successfully!")

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
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func SaveScore(username, difficulty string, timeSpent int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	query := `INSERT INTO scores (username, difficulty, time_spent) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, username, difficulty, timeSpent)
	return err
}

func GetTopScores(difficulty string) ([]Score, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
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
