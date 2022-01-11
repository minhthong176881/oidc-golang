package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/minhthong176881/oidc-golang/pkg/client/rp"
	httphelper "github.com/minhthong176881/oidc-golang/pkg/http"
	"github.com/minhthong176881/oidc-golang/pkg/oidc"
)

var (
	callbackPath = "/auth/callback"
	key          = []byte("test1234test1234")
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	keyPath := os.Getenv("KEY_PATH")
	issuer := os.Getenv("ISSUER")
	port := os.Getenv("PORT")
	scopes := strings.Split(os.Getenv("SCOPES"), " ")

	templatePath := "example/client/view/"

	redirectURI := fmt.Sprintf("http://localhost:%v%v", port, callbackPath)
	cookieHandler := httphelper.NewCookieHandler(key, key, httphelper.WithUnsecure())

	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
	}
	if clientSecret == "" {
		options = append(options, rp.WithPKCE(cookieHandler))
	}
	if keyPath != "" {
		options = append(options, rp.WithClientKey(keyPath))
	}

	provider, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectURI, scopes, options...)
	if err != nil {
		logrus.Fatalf("error creating provider %s", err.Error())
	}

	//generate some state (representing the state of the user in your application,
	//e.g. the page where he was before sending him to login
	state := func() string {
		return uuid.New().String()
	}

	render := func(w http.ResponseWriter, path string, data interface{}) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.WriteHeader(http.StatusOK)
		tlp := template.Must(template.ParseFiles(path))
		tlp.Execute(w, data)
	}

	http.Handle("/", rp.Dashboard(render, templatePath + "index.html", provider))

	//register the AuthURLHandler at your preferred path
	//the AuthURLHandler creates the auth request and redirects the user to the auth server
	//including state handling with secure cookie and the possibility to use PKCE
	http.Handle("/login", rp.AuthURLHandler(state, provider))

	//for demonstration purposes the returned userinfo response is written as JSON object onto response
	// marshalUserinfo := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty, info oidc.UserInfo) {
	// 	data, err := json.Marshal(info)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Write(data)
	// }
	marshalUserinfo := func(w http.ResponseWriter, r *http.Request, info oidc.UserInfo) {
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}

	//you could also just take the access_token and id_token without calling the userinfo endpoint:
	marshalToken := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
		// data, err := json.Marshal(tokens)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokens.AccessToken, Path: "/"})
		http.SetCookie(w, &http.Cookie{Name: "token_type", Value: tokens.TokenType, Path: "/"})
		http.SetCookie(w, &http.Cookie{Name: "id_token", Value: tokens.IDToken, Path: "/"})
		render(w, templatePath + "callback.html", tokens)
	}

	//register the CodeExchangeHandler at the callbackPath
	//the CodeExchangeHandler handles the auth response, creates the token request and calls the callback function
	//with the returned tokens from the token endpoint
	//in this example the callback function itself is wrapped by the UserinfoCallback which
	//will call the Userinfo endpoint, check the sub and pass the info into the callback function
	// http.Handle(callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(marshalUserinfo), provider))

	//if you would use the callback without calling the userinfo endpoint, simply switch the callback handler for:
	//
	http.Handle(callbackPath, rp.CodeExchangeHandler(marshalToken, provider))

	http.Handle("/userinfo", rp.UserInfoExchangeHandler(marshalUserinfo, provider))

	http.Handle("/logout", rp.LogoutHandler(provider))

	lis := fmt.Sprintf("127.0.0.1:%s", port)
	logrus.Infof("listening on http://%s/", lis)
	logrus.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}
