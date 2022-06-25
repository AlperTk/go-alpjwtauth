package authorization

import "net/http"

type Authorizer interface {
	Process(roles []string, w http.ResponseWriter, r *http.Request) (defined bool, err error)
}
