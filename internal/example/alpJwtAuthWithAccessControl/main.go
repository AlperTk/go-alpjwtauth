package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlperTk/go-alpjwtauth/internal/example/alpJwtAuthWithAccessControl/config"
	"github.com/AlperTk/go-alpjwtauth/src/accesscontrol"
	"github.com/AlperTk/go-alpjwtauth/src/authorization"
	"github.com/gorilla/mux"
	"net/http"
)

type ApplicationStarter struct {
	AlpJwtAuth authorization.AlpJwtAuth
}

func main() {
	fmt.Println("Server starting...")
	load().run()
}

func load() ApplicationStarter {
	tokenProcessor := authorization.NewKeycloakTokenProcessor("https://localhost:8443/auth/realms/marsrealm/protocol/openid-connect/certs")

	webSecurity := securityConfig.SecurityConfig{}
	alpAuthorizer := accesscontrol.NewBasicRoleAuthorizer(webSecurity)
	alpJwtAuth := authorization.NewJwtAuthWithAccessControl(tokenProcessor, alpAuthorizer)

	p := ApplicationStarter{
		AlpJwtAuth: alpJwtAuth,
	}
	return p
}

func (p ApplicationStarter) run() {
	router := mux.NewRouter().StrictSlash(true)
	p.AlpJwtAuth.SetupMux(router)

	router.Handle("/api/v1/test", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_ = json.NewEncoder(writer).Encode("AlperTk")
	})).Methods("POST")

	_ = http.ListenAndServe(":9702", router)
}
