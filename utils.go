package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func getSession(r *http.Request, name string) *sessions.Session {
	s, err := store.Get(r, name)
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}
	return s

}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	templates := template.Must(template.ParseGlob("templates/*"))
	//	t, err := template.ParseFiles(tmpl + ".html")
	/*	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	//	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
