package main

import (
	"fmt"
	securityConfig "github.com/AlperTk/go-alpjwtauth/example/alpJwtAuthWithAccessControl/config"
	"github.com/AlperTk/go-alpjwtauth/src/authentication"
	"github.com/AlperTk/go-alpjwtauth/src/authentication/impl"
	authorization "github.com/AlperTk/go-alpjwtauth/src/authorization/service/imp"
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

	//router.Handle("/api/v1/test", http.HandlerFunc(testFunc)).Methods("POST")

	log.Fatal(http.ListenAndServe(":9702", router))
}
