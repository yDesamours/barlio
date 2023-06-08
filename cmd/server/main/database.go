package main

import "database/sql"

func newDB(dsn string) (*sql.DB, error) {
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func (app *App) setDB(db *sql.DB) {
	app.DB = db
}
