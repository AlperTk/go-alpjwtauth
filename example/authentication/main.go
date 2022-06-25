package main

import (
	"fmt"
	"github.com/AlperTk/go-jwt-role-based-auth/example/authentication/config"
	controllers2 "github.com/AlperTk/go-jwt-role-based-auth/example/authentication/controllers"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication"
	"github.com/AlperTk/go-jwt-role-based-auth/src/authentication/impl"
	authorization "github.com/AlperTk/go-jwt-role-based-auth/src/authorization/service/imp"
	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/logrus"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/mux"
	logrusImp "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

type ApplicationStarter struct {
	AlpJwtAuth authentication.AlpJwtAuth
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

	registerRoutes(router)
	log.Fatal(http.ListenAndServe(":9702", router))
}

func registerRoutes(router *mux.Router) {
	registerControllerRoutes(controllers2.EventController{}, router)
}

func registerControllerRoutes(controller controllers2.Controller, router *mux.Router) {
	controller.RegisterRoutes(router)
}
