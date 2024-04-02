package utility

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Postgres *sql.DB

func InitPostgres() {
	err := godotenv.Load()
	if err != nil {
		Slogger.Error("加载环境变量失败")
		log.Fatal(err.Error())
	}
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		database,
	)
	Postgres, err = sql.Open("postgres", dsn)
	if err != nil {
		Slogger.Error(err.Error())
		log.Fatal(err.Error())
	}
	Postgres.SetConnMaxLifetime(time.Second * 30)
	Postgres.SetMaxIdleConns(runtime.NumCPU()*2 + 1)
	if err = Postgres.Ping(); err != nil {
		Slogger.Error("连接数据库失败")
		log.Fatal(err.Error())
	}
}
