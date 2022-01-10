package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/minhthong176881/oidc-golang/example/internal/mock"
	"github.com/minhthong176881/oidc-golang/pkg/op"
)

func main() {
	ctx := context.Background()
	port := "9998"
	config := &op.Config{
		Issuer:    "http://localhost:9998/",
		CryptoKey: sha256.Sum256([]byte("test")),
	}
	storage := mock.NewAuthStorage()
	handler, err := op.NewOpenIDProvider(ctx, config, storage, op.WithCustomTokenEndpoint(op.NewEndpoint("test")))
	if err != nil {
		log.Fatal(err)
	}
	router := handler.HttpHandler().(*mux.Router)
	router.Methods("GET").Path("/login").HandlerFunc(HandleLogin)
	router.Methods("POST").Path("/login").HandlerFunc(HandleCallback)
	router.Methods("GET").Path("/interaction/{router}").HandlerFunc(HandleInteraction)
	router.Methods("POST").Path("/interaction/{router}/confirm").HandlerFunc(HandleInteractionConfirm)
	router.Methods("GET").Path("/interaction/{router}/abort").HandlerFunc(HandleInteractionAbort)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
}

func render(w http.ResponseWriter, path string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	tlp := template.Must(template.ParseFiles(path))
	tlp.Execute(w, data)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	render(w, "example/server/view/login.html", nil)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	client := r.FormValue("client")
	// http.Redirect(w, r, "/interaction?id=" + client, http.StatusFound)
	http.Redirect(w, r, "/authorize/callback?id="+client, http.StatusFound)
}

func HandleInteraction(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("/interaction/%s", mux.Vars(r)["router"]) + "/confirm"
	data := map[string]interface{}{
		"url": url,
	}
	render(w, "example/server/view/interaction.html", data)
}

func HandleInteractionConfirm(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "callback" {
			http.Redirect(w, r, cookie.Value, http.StatusFound)
		}
	}
}

func HandleInteractionAbort(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
