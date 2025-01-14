package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thomascastle/rangoon-eats/internal/structuredlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configuration struct {
	db struct {
		dsn string
	}
	env  string
	port int
}

func configure() (configuration, error) {
	if e := godotenv.Load(); e != nil {
		return configuration{}, e
	}

	dsn := os.Getenv("DSN")

	env := os.Getenv("ENV")

	port, e := strconv.Atoi(os.Getenv("PORT"))
	if e != nil {
		return configuration{}, e
	}

	var config configuration
	config.db.dsn = dsn
	config.env = env
	config.port = port

	return config, nil
}

func main() {
	logger := structuredlog.New(os.Stdout, structuredlog.LevelInfo)

	config, e := configure()
	if e != nil {
		logger.Fatal(e, nil)
	}

	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"env":    config.env,
			"status": "available",
		}

		data_JSON, e := json.Marshal(data)
		if e != nil {
			logger.Fatal(e, nil)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data_JSON)
	})

	_, e = openDB(config)
	if e != nil {
		logger.Fatal(e, nil)
	}

	logger.Info("database connection pool established", nil)

	logger.Info("server started", map[string]string{"addr": fmt.Sprintf(":%d", config.port), "env": config.env})
	if e := http.ListenAndServe(fmt.Sprintf(":%d", config.port), nil); e != nil {
		logger.Fatal(e, nil)
	}
}

func openDB(config configuration) (*gorm.DB, error) {
	db, e := openUnderlyingConnection(config.db.dsn)
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
