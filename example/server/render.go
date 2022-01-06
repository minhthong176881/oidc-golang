package main

import (
	"html/template"
	"net/http"
)

var (
	loginTemplate = template.Must(template.ParseFiles("./view/login.html"))
)

func renderLogin(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	loginTemplate.Execute(w, nil)
}
