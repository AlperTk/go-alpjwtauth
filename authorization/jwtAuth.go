package authorization

import (
	"crypto/tls"
	"encoding/json"
	errors2 "errors"
	"github.com/AlperTk/go-alpjwtauth/accesscontrol"
	"github.com/AlperTk/go-alpjwtauth/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type jwtAuth struct {
	TokenProcessor TokenProcessor
	RoleAuthor     accesscontrol.AlpAuthorizer
}

func NewJwtAuth(processor TokenProcessor) AlpJwtAuth {
	return jwtAuth{TokenProcessor: processor}
}

func NewJwtAuthWithAccessControl(processor TokenProcessor, authorizer accesscontrol.AlpAuthorizer) AlpJwtAuth {
	return jwtAuth{TokenProcessor: processor, RoleAuthor: authorizer}
}

func init() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}

func (j jwtAuth) SetupMux(router *mux.Router) {
	router.Use(j.protect)
	router.NotFoundHandler = j.responseHandler(notFound)
	router.MethodNotAllowedHandler = j.responseHandler(methodNotAllowed)
}

func (j jwtAuth) protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > 0 {
			roles, err := j.tokenValidate(r)
			if err != nil {
				responseUnauthorized(w)
				return
			}

			if j.RoleAuthor != nil {
				_, e := j.RoleAuthor.ProcessAuthorized(roles, w, r)
				if e != nil {
					return
				}
			}

			next.ServeHTTP(w, r)
		} else {
			_, e := j.RoleAuthor.ProcessUnauthorized(w, r)
			if e == nil {
				next.ServeHTTP(w, r)
				return
			}
		}

	})
}

func (j jwtAuth) responseHandler(response func(w http.ResponseWriter)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > 0 {
			_, err := j.tokenValidate(r)

			if err == nil {
				response(w)
				return
			}

		} else {
			if j.RoleAuthor != nil {
				defined, e := j.RoleAuthor.ProcessUnauthorized(w, r)
				if defined && e == nil {
					response(w)
				}
				return
			}
		}

		responseUnauthorized(w)
	})
}

func (j jwtAuth) tokenValidate(r *http.Request) ([]string, error) {
	authHeader := r.Header.Get("Authorization")

	if len(authHeader) < 1 {
		return nil, errors2.New("no accesscontrol token find")
	}

	bearerToken := strings.Split(authHeader, " ")[1]

	valid, roles, err := j.TokenProcessor.Process(bearerToken, r)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors2.New("token not valid")
	}
	return roles, nil
}

func responseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	_ = json.NewEncoder(w).Encode(errors.UnauthorizedError())
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	_ = json.NewEncoder(w).Encode(errors.NotFound())
}

func methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(405)
	_ = json.NewEncoder(w).Encode(errors.MethodNotAllowed())
}
