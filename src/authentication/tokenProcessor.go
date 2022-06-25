package authentication

import "net/http"

type TokenProcessor interface {
	Process(bearerToken string, r *http.Request) (bool, []string, error)
}
