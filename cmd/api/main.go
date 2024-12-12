package main

import (
	"air-monolith/internal/app"
	"air-monolith/internal/jsonlog"
	"air-monolith/internal/repository/dbrepo"
	"flag"
	"fmt"
	"os"
	"time"
)

const TIMEOUT = time.Duration(120 * time.Second)

func main() {
	var cfg app.Config

	cfg.Timeout = TIMEOUT

	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Domain, "domain", "localhost", "domain")
	flag.StringVar(&cfg.Env, "env", "dev", "dev|main")
	flag.StringVar(&cfg.DSN, "dsn", "host=postgres port=5432 user=postgres password=password dbname=segments sslmode=disable timezone=UTC connect_timeout=5", "Database Source Name")

	var app app.Application
	app.Cfg = cfg
	app.Logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	conn, err := app.ConnectToDB()
	if err != nil {
		fmt.Println(err)
		app.Logger.PrintFatal(err, nil)
		return
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	err = app.Server()
	if err != nil {
		app.Logger.PrintFatal(err, nil)
	}
}
