package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataStore *PostgresDB

type PostgresDB struct {
	DB     *sql.DB
	GormDB *gorm.DB
}

func InitPostgressConnection() *PostgresDB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"))

	log.Println(connStr)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error while connecting to DB:", err)
		return nil
	}

	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(30)
	conn.SetConnMaxLifetime(time.Hour)
	err = conn.Ping()
	if err != nil {
		log.Println("Error while connecting to DB:", err)
		return nil
	}
	// gorm connection
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connStr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Println("Error while creating GORM connection to DB: ", err)
	}

	return &PostgresDB{
		DB:     conn,
		GormDB: gormDB,
	}
}

func Init() {
	log.Println("Connecting to Postgres DB")
	conn := InitPostgressConnection()
	if conn == nil {
		log.Fatal("Error While Connecting To Postgres DB")
		return
	}
	DataStore = conn

	log.Println("Successfully connected to Postgres DB")
}
