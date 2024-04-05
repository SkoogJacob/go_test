package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to db, running a ping...")
	err = db.Ping()
	log.Println("Ping successful!")
	return db, err
}

func (s *server) connectToDB() (*sql.DB, error) {
	connection, err := openDB(s.DSN)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
