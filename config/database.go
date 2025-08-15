package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDatabase() *sql.DB {
	connStr := "user=vistartr password=Password1! dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}
	fmt.Println("Berhasil terhubung ke database!")
	return db
}
