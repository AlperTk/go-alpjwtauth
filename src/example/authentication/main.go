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
	"github.com/joho/godotenv"
	logrusImp "github.com/sirupsen/logrus"
	"net/http"
	"os"
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
		return " <-- (test)"
	}})
	log.Current = logrus.New(logger)
}

func load() ApplicationStarter {
	err := godotenv.Load("properties.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var tokenProcessor authentication.TokenProcessor
	tokenProcessor = keycloak.NewKeycloakTokenProcessor(os.Getenv("KEYCOLAK_ADDRESS"))

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
	router.Use(defaultResponseTypeSetter)
	p.JwtAuth.SetupMux(router)

	registerRoutes(router)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDR"), router))
}

func registerRoutes(router *mux.Router) {
	registerControllerRoutes(controllers.EventController{}, router)
}

func registerControllerRoutes(controller controllers.Controller, router *mux.Router) {
	controller.RegisterRoutes(router)
}

func defaultResponseTypeSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
