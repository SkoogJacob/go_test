package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"web_test/pkg/data"
	"web_test/pkg/repository/dbrepo"
)

func main() {
	s := server{}

	// Set up DB connection
	flag.StringVar(&s.DSN,
		"dsn",
		"host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connect info")
	flag.Parse()
	d, err := s.connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	s.DB = &dbrepo.PostgresDBRepo{DB: d}
	defer s.closeDB()

	s.Session = getSession()
	gob.Register(data.User{})

	mux := s.routes()
	log.Print("Starting server...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
