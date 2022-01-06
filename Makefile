server: 
	cd example/server && go build && ./server
client: 
	CLIENT_ID=web CLIENT_SECRET=web ISSUER=http://localhost:9998/ SCOPES=openid PORT=5556 go run example/client/app/app.go