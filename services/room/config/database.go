package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "postgres://postgres.iurddbrzlrhqcdwjauvw:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Iml1cmRkYnJ6bHJocWNkd2phdXZ3Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDUyODgzMzIsImV4cCI6MjA2MDg2NDMzMn0.lvzzluq_W9Nj4YLJMKCwq_8Ei0YAhCv9Q4gjhjlBReI@db.iurddbrzlrhqcdwjauvw.supabase.co:5432/postgres"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Successfully connected to database")
}
