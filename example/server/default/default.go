package main

import (
	"context"
	"crypto/sha256"
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
	// tpl := `
	// <!DOCTYPE html>
	// <html>
	// 	<head>
	// 		<meta charset="UTF-8">
	// 		<title>Login</title>
	// 	</head>
	// 	<body>
	// 		<form method="POST" action="/login">
	// 			<input name="client"/>
	// 			<button type="submit">Login</button>
	// 		</form>
	// 	</body>
	// </html>`
	tpl := `<!DOCTYPE html>
	<html >
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<meta http-equiv="x-ua-compatible" content="ie=edge">
		<title>Sign-in</title>
		<style>
		  @import url(https://fonts.googleapis.com/css?family=Roboto:400,100);
	
		  body {
			font-family: 'Roboto', sans-serif;
			margin-top: 25px;
			margin-bottom: 25px;
		  }
	
		  .login-card {
			padding: 40px;
			padding-top: 0px;
			padding-bottom: 10px;
			width: 274px;
			background-color: #F7F7F7;
			margin: 0 auto 10px;
			border-radius: 2px;
			box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
			overflow: hidden;
		  }
	
		  .login-card + .login-card {
			padding-top: 10px;
		  }
	
		  .login-card h1 {
			font-weight: 100;
			text-align: center;
			font-size: 2.3em;
		  }
	
		  .login-card h1 + p {
			font-weight: 100;
			text-align: center;
		  }
	
		  .login-card [type=submit] {
			width: 100%;
			display: block;
			margin-bottom: 10px;
			position: relative;
		  }
	
		  .login-card input[type=text], input[type=email], input[type=password] {
			height: 44px;
			font-size: 16px;
			width: 100%;
			margin-bottom: 10px;
			-webkit-appearance: none;
			background: #fff;
			border: 1px solid #d9d9d9;
			border-top: 1px solid #c0c0c0;
			padding: 0 8px;
			box-sizing: border-box;
			-moz-box-sizing: border-box;
		  }
	
		  .login {
			text-align: center;
			font-size: 14px;
			font-family: 'Arial', sans-serif;
			font-weight: 700;
			height: 36px;
			padding: 0 8px;
		  }
	
		  .login-submit {
			border: 0px;
			color: #fff;
			text-shadow: 0 1px rgba(0,0,0,0.1);
			background-color: #4d90fe;
		  }
	
		  .login-card a {
			text-decoration: none;
			color: #666;
			font-weight: 400;
			text-align: center;
			display: inline-block;
			opacity: 0.6;
		  }
	
		  .login-help {
			width: 100%;
			text-align: center;
			font-size: 12px;
		  }
	
		  .login-client-image img {
			margin-bottom: 20px;
			display: block;
			margin-left: auto;
			margin-right: auto;
			width: 20%;
		  }
	
		  .login-card input[type=checkbox] {
			margin-bottom: 10px;
		  }
	
		  .login-card label {
			color: #999;
		  }
	
		  ul {
			font-weight: 100;
			padding-left: 1em;
			list-style-type: circle;
		  }
	
		  li + ul, ul + li, li + li {
			padding-top: 0.3em;
		  }
	
		  button {
			cursor: pointer;
		  }
		</style>
	  </head>
	  <body>
		<div class="login-card">
		  <h1>Login</h1>
		  <form autocomplete="off" action="/login" method="post">
		  	<input required type="text" name="client" placeholder="Enter id to login">
			<button type="submit" class="login login-submit">Login</button>
		  </form>
		</div>
	  </body>
	</html>
	`
	t, err := template.New("login").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	client := r.FormValue("client")
	http.Redirect(w, r, "/authorize/callback?id="+client, http.StatusFound)
}
