package app

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (app *Application) openDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", app.Cfg.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) ConnectToDB() (*sql.DB, error) {
	conn, err := app.openDB()
	if err != nil {
		return nil, err
	}

	return conn, err
}
