package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal koneksi ke MySQL: %v", err)
	}

	fmt.Printf("Berhasil terhubung ke database MySQL %s!\n", dbName)
	return db
}