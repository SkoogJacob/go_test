package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var pathToTemplates = "./templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (s *server) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)
	if s.Session.Exists(r.Context(), "test") {
		msg := s.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		s.Session.Put(r.Context(), "test", "hit this page at "+time.Now().UTC().String())
	}
	if s.Session.Exists(r.Context(), "login_error") {
		msg := s.Session.GetString(r.Context(), "login_error")
		td["error"] = msg
	}
	_ = s.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

func (s *server) Profile(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)
	if s.Session.Exists(r.Context(), "flash") {
		msg := s.Session.GetString(r.Context(), "flash")
		td["test"] = msg
	} else {
		_ = s.Session.RenewToken(r.Context())
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if s.Session.Exists(r.Context(), "error") {
		s.Session.Remove(r.Context(), "error")
	}
	_ = s.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

func (s *server) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	parsed, err := template.ParseFiles(
		path.Join(pathToTemplates, t),
		path.Join(pathToTemplates, "base.layout.gohtml"),
	)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}
	data.IP = s.ipFromContext(r.Context())
	err = parsed.Execute(w, data)
	return err
}

func (s *server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	form := NewForm(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		s.Session.Put(r.Context(), "login_error", "Login information not submitted properly")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	user, err := s.DB.GetUserByEmail(email)
	if err != nil {
		s.Session.Put(r.Context(), "login_error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Prevent fixation attacks
	_ = s.Session.RenewToken(r.Context())

	log.Println("From db: ", user)
	log.Println(user, password)

	s.Session.Put(r.Context(), "flash", "successfully logged in")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
