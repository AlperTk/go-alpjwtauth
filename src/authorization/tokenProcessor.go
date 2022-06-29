package authorization

import "net/http"

type TokenProcessor interface {
	Process(bearerToken string, r *http.Request) (valid bool, roles []string, err error)
}
