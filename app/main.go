package main

import (
	"log"
	"net/http"
)

func main() {
	s := server{}
	s.Session = getSession()
	mux := s.routes()
	log.Print("Starting server...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
