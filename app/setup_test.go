package main

import (
	"encoding/gob"
	"log"
	"os"
	"testing"
	"web_test/pkg/data"
	"web_test/pkg/db"
)

var s server

func TestMain(m *testing.M) {
	pathToTemplates = "../templates/"
	// Set up DB connection

	s.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	d, err := s.connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	s.DB = db.PostgresConn{DB: d}
	defer s.closeDB()

	s.Session = getSession()
	gob.Register(data.User{})

	os.Exit(m.Run())
}
