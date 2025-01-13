package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if e := godotenv.Load(); e != nil {
		log.Fatal("error: failed to load .env file")
	}

	env := os.Getenv("APP_ENV")

	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"env":    env,
			"status": "available",
		}

		data_JSON, e := json.Marshal(data)
		if e != nil {
			log.Fatal(e)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data_JSON)
	})

	_, e := openDB()
	if e != nil {
		log.Fatal(e)
	}

	log.Println("database connection pool established")

	log.Println("server started")
	if e := http.ListenAndServe(":4000", nil); e != nil {
		log.Fatal(e)
	}
}

func openDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	port, e := strconv.Atoi(os.Getenv("DB_PORT"))
	if e != nil {
		return nil, e
	}
	username := os.Getenv("DB_USERNAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, name)

	db, e := openUnderlyingConnection(dsn)
	if e != nil {
		return nil, e
	}

	db_GORM, e := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if e != nil {
		return nil, e
	}

	return db_GORM, nil
}

func openUnderlyingConnection(dsn string) (*sql.DB, error) {
	db, e := sql.Open("postgres", dsn)
	if e != nil {
		return nil, e
	}

	e = db.Ping()
	if e != nil {
		return nil, e
	}

	return db, nil
}
