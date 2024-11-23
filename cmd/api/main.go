package main

import (
	"air-monolith/internal/jsonlog"
	"air-monolith/internal/repository"
	"air-monolith/internal/repository/dbrepo"
	"flag"
	"fmt"
	"os"
	"time"
)

type config struct {
	port   int
	Domain string
	env    string
	DSN    string
	db     struct {
		maxOpensConns int
		maxIdleConns  int
		maxIdleTime   string
	}

	timeout time.Duration
}

type application struct {
	cfg         config
	logger      *jsonlog.Logger
	DB          repository.DatabaseRepo
	statusWrote chan bool
}

func main() {
	var cfg config

	cfg.timeout = time.Duration(120 * time.Second)

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Domain, "domain", "localhost", "domain")
	flag.StringVar(&cfg.env, "env", "dev", "docker|dev|main|")
	flag.StringVar(&cfg.DSN, "dsn", "host=localhost port=5432 user=postgres password=password dbname=segments sslmode=disable timezone=UTC connect_timeout=5", "Database Source Name")

	flag.IntVar(&cfg.db.maxOpensConns, "db-max-opens-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max edle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connections idle time")

	var app application
	app.cfg = cfg
	app.logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app.statusWrote = make(chan bool)

	conn, err := app.connectToDB()
	if err != nil {
		fmt.Println(err)
		app.logger.PrintFatal(err, nil)
		return
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	err = app.server()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
