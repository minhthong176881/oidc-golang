server: 
	CAOS_OIDC_DEV=1 go run example/server/default/default.go
client: 
	CLIENT_ID=web CLIENT_SECRET=web ISSUER=http://localhost:9998/ SCOPES=openid PORT=5556 go run example/client/app/app.go