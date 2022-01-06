package main

import (
	"context"
	"crypto/sha256"
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	renderLogin(w)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	client := r.FormValue("client")
	http.Redirect(w, r, "/authorize/callback?id="+client, http.StatusFound)
}
