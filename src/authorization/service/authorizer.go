package authorization

import "net/http"

type Authorizer interface {
	ProcessUnauthorized(w http.ResponseWriter, r *http.Request) (defined bool, err error)
	ProcessAuthorized(roles []string, w http.ResponseWriter, r *http.Request) (defined bool, err error)
}
