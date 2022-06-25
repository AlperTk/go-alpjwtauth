package authentication

import "github.com/gorilla/mux"

type AlpJwtAuth interface {
	SetupMux(router *mux.Router)
}
