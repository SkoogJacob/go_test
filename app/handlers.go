package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
	"web_test/pkg/data"
)

var pathToTemplates = "./templates/"

type TemplateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (s *server) Home(w http.ResponseWriter, r *http.Request) {
	var td = TemplateData{
		Data: make(map[string]any),
	}
	if s.Session.Exists(r.Context(), "test") {
		msg := s.Session.GetString(r.Context(), "test")
		td.Data["test"] = msg
	} else {
		s.Session.Put(r.Context(), "test", "hit this page at "+time.Now().UTC().String())
	}
	_ = s.render(w, r, "home.page.gohtml", &td)
}

func (s *server) Profile(w http.ResponseWriter, r *http.Request) {
	var td = TemplateData{Data: make(map[string]any)}
	if !s.Session.Exists(r.Context(), "user") {
		_ = s.Session.RenewToken(r.Context())
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if s.Session.Exists(r.Context(), "error") {
		s.Session.Remove(r.Context(), "error")
	}
	_ = s.render(w, r, "profile.page.gohtml", &td)
}

func (s *server) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	parsed, err := template.ParseFiles(
		path.Join(pathToTemplates, t),
		path.Join(pathToTemplates, "base.layout.gohtml"),
	)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}
	td.IP = s.ipFromContext(r.Context())
	td.Error = s.Session.PopString(r.Context(), "error")
	td.Flash = s.Session.PopString(r.Context(), "flash")
	user := s.Session.Get(r.Context(), "user")
	if user != nil {
		td.User = user.(data.User)
	}
	err = parsed.Execute(w, td)
	return err
}

// Login Not a page, handles a post request for login and then redirect to
// the proper page based on whether the login was successful or not
func (s *server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	form := NewForm(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		log.Printf("Login form did not have all required information, got %v\n", form.Data)
		s.Session.Put(r.Context(), "error", "Login information not submitted properly")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	user, err := s.DB.GetUserByEmail(email)
	if err != nil {
		log.Printf("Failed to authenticate user with error [%v\n]", err)
		s.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Prevent fixation attacks
	_ = s.Session.RenewToken(r.Context())

	if !s.authenticateUser(r, user, password) {
		s.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	s.Session.Put(r.Context(), "flash", "successfully logged in")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
