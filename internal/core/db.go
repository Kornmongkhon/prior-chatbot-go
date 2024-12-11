package core

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDb(hostname string, port string, username string, password string, database string) *sql.DB {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, hostname, port, database)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
		return nil
	}
	// ทดสอบการเชื่อมต่อ
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
		return nil
	}
	log.Println("Database connected")
	return db
}
