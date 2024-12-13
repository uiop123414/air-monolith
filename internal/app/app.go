package app

import (
	"air-monolith/internal/jsonlog"
	"air-monolith/internal/repository"
	"time"
)

type Config struct {
	Port    int
	Domain  string
	Env     string
	DSN     string
	Timeout time.Duration
}

type Application struct {
	Cfg    Config
	Logger *jsonlog.Logger
	DB     repository.DatabaseRepo
}
