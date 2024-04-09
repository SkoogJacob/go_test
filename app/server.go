package main

import (
	"log"
	"net/http"
	"web_test/pkg/data"
	"web_test/pkg/repository"

	"github.com/alexedwards/scs/v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	Session *scs.SessionManager
	DB      repository.DatabaseRepo
	DSN     string
}

func (s *server) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(s.addIpToContext)
	mux.Use(s.Session.LoadAndSave)

	mux.Get("/", s.Home)
	mux.Post("/login", s.Login)

	mux.Route("/user", func(r chi.Router) {
		r.Use(s.auth)
		r.Get("/profile", s.Profile)
	})

	filSrv := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", filSrv))
	return mux
}

func (s *server) closeDB() {
	err := s.DB.Connection().Close()
	if err != nil {
		log.Fatalf("Unable to close DB connection: %v\n", err)
	}
}

func (s *server) authenticateUser(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}
	s.Session.Put(r.Context(), "user", user)
	return true
}
