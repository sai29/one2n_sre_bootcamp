package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sai29/one2n_sre_bootcamp/internal/data"
	"github.com/sai29/one2n_sre_bootcamp/internal/jsonlog"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (dev|stage|prod)")

	pw := os.Getenv("DB_PW")
	flag.StringVar(&cfg.db.dsn, "db-dsn", fmt.Sprintf("postgres://student_api:%s@localhost/student_api?sslmode=disable",
		pw), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Postgres max open conns")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Postgres max idle conns")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Postgres max conn idle time")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
