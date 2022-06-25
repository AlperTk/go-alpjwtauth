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
	"github.com/AlperTk/go-jwt-role-based-auth/example/authentication/config"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication/impl"
	authorization "github.com/AlperTk/go-jwt-role-based-auth/src/authorization/service/imp"
	"github.com/Masterminds/log-go"
	"github.com/gorilla/mux"
	"net/http"
)

type ApplicationStarter struct {
	AlpJwtAuth authentication.AlpJwtAuth
}

func main() {
	fmt.Println("Server starting...")
	load().run()
}

func load() ApplicationStarter {
	tokenProcessor := impl.NewKeycloakTokenProcessor("https://localhost:8443/auth/realms/marsrealm/protocol/openid-connect/certs")

	webSecurity := securityConfig.WebSecurityConfig{}
	alpAuthorizer := authorization.NewBasicRoleAuthorizer(webSecurity)
	alpJwtAuth := impl.NewJwtAuthWithAccessControl(tokenProcessor, alpAuthorizer)

	p := ApplicationStarter{
		AlpJwtAuth: alpJwtAuth,
	}
	return p
}

func (p ApplicationStarter) run() {
	router := mux.NewRouter().StrictSlash(true)
	p.AlpJwtAuth.SetupMux(router)

	//router.Handle("/api/v1/test", http.HandlerFunc(testFunc)).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":9702", router))
}

```
