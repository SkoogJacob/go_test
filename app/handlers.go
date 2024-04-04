package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

var pathToTemplates = "./templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (s *server) Home(w http.ResponseWriter, r *http.Request) {
	_ = s.render(w, r, "home.page.gohtml", &TemplateData{})
}

func (s *server) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	parsed, err := template.ParseFiles(path.Join(pathToTemplates, t))
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Errors in submitted form:")
		for k, v := range form.Errors {
			fmt.Fprintf(w, "%s: %v\n", k, v)
		}
	} else {
		w.WriteHeader(200)
		fmt.Fprintln(w, "Login form was formed correctly, bu no action is implemented")
	}
}
