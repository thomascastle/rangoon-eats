package main

import (
	"database/sql"
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

type application struct {
	config configuration
	logger *structuredlog.Logger
}

func main() {
	logger := structuredlog.New(os.Stdout, structuredlog.LevelInfo)

	config, e := configure()
	if e != nil {
		logger.Fatal(e, nil)
	}

	_, e = openDB(config)
	if e != nil {
		logger.Fatal(e, nil)
	}

	logger.Info("database connection pool established", nil)

	app := application{
		config: config,
		logger: logger,
	}

	if e := app.serve(); e != nil {
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
