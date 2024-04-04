package main

import (
	"github.com/alexedwards/scs/v2"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	Session *scs.SessionManager
}

func (s *server) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(s.addIpToContext)
	mux.Use(s.Session.LoadAndSave)

	mux.Get("/", s.Home)
	mux.Post("/login", s.Login)

	filSrv := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", filSrv))
	return mux
}
