package main

import (
	"JwtAuth/src/authentication"
	"JwtAuth/src/authentication/impl/keycloak"
	roleAuth "JwtAuth/src/authorization/service/imp"
	securityConfig "JwtAuth/src/example/authentication/config"
	"JwtAuth/src/example/authentication/controllers"
	"fmt"
	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/logrus"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/mux"
	logrusImp "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

type ApplicationStarter struct {
	JwtAuth authentication.JwtAuth
}

func main() {
	fmt.Println("Server starting...")
	load().run()
}

func init() {
	// init logger

	logger := logrusImp.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&nested.Formatter{CustomCallerFormatter: func(frame *runtime.Frame) string {
		pc, _, line, ok := runtime.Caller(10)
		details := runtime.FuncForPC(pc)
		var funcName string
		if ok && details != nil {
			funcName = details.Name()
			return fmt.Sprintf(" <-- (%s:%d)", funcName, line)
		}
		return " <-- (Unknown)"
	}})
	log.Current = logrus.New(logger)
}

func load() ApplicationStarter {
	tokenProcessor := keycloak.NewKeycloakTokenProcessor("https://localhost:8443/auth/realms/marsrealm/protocol/openid-connect/certs")

	webSecurity := securityConfig.WebSecurityConfig{}

	jwtAuth := authentication.JwtAuth{
		TokenProcessor: tokenProcessor,
		RoleAuthor:     roleAuth.NewBasicRoleAuthorizer(webSecurity),
	}

	p := ApplicationStarter{
		JwtAuth: jwtAuth,
	}
	return p
}

func (p ApplicationStarter) run() {
	router := mux.NewRouter().StrictSlash(true)
	p.JwtAuth.SetupMux(router)

	registerRoutes(router)
	log.Fatal(http.ListenAndServe(":9702", router))
}

func registerRoutes(router *mux.Router) {
	registerControllerRoutes(controllers.EventController{}, router)
}

func registerControllerRoutes(controller controllers.Controller, router *mux.Router) {
	controller.RegisterRoutes(router)
}
