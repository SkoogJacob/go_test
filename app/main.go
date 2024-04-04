package main

import (
	"log"
	"net/http"
)

func main() {
	srv := server{}
	mux := srv.routes()
	log.Print("Starting server...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
