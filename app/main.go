package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	s := server{}
	flag.StringVar(&s.DSN,
		"dsn",
		"host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connect info")
	flag.Parse()
	d, err := s.connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	s.DB = d
	s.Session = getSession()
	mux := s.routes()
	log.Print("Starting server...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
