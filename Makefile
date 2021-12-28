server: 
	CAOS_OIDC_DEV=1 go run github.com/caos/oidc/example/server/default
client: 
	CLIENT_ID=web CLIENT_SECRET=web ISSUER=http://localhost:9998/ SCOPES=openid PORT=5556 go run github.com/caos/oidc/example/client/app