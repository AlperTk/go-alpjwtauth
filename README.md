# Go JWT Role Based Authorization

## Usage

```golang
import "github.com/AlperTk/go-jwt-role-based-auth/src/authorization/builder/roleBuilder"

type SecurityConfig struct {
}

func (s SecurityConfig) Config(security *roleBuilder.RoleConfigurer) {
	security.
		AntMatcher("/api/v1/test", "/api/v1/test2").HasAnyRoles("admin", "test-role").
		AntMatcher("/api/v2/**").PermitAll().
		AntMatcher("/api/v3/**").DenyAll().
		AntMatcher("/api/v4/**").Authenticated().
		AnyRequest().DenyAll()
}
```

```golang
package main

import (
	"fmt"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication/impl/keycloak"
	authorization "github.com/AlperTk/go-jwt-role-based-auth/src/authorization/service/imp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"test/src/config"
)

type ApplicationStarter struct {
	JwtAuth authentication.JwtAuth
}

func main() {
	fmt.Println("Server starting...")
	load().run()
}

func load() ApplicationStarter {
	tokenProcessor := keycloak.NewKeycloakTokenProcessor("https://localhost:8443/auth/realms/marsrealm/protocol/openid-connect/certs")

	securityConfig := config.SecurityConfig{}

	jwtAuth := authentication.JwtAuth{
		TokenProcessor: tokenProcessor,
		RoleAuthor:     authorization.NewBasicRoleAuthorizer(securityConfig),
	}

	p := ApplicationStarter{
		JwtAuth: jwtAuth,
	}
	return p
}

func (p ApplicationStarter) run() {
	router := mux.NewRouter().StrictSlash(true)
	p.JwtAuth.SetupMux(router)

	//router.Handle("/api/v1/test", http.HandlerFunc(postTest)).Methods("POST")

	log.Fatal(http.ListenAndServe(":9701", router))
}

```
