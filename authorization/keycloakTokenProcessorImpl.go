package authorization

import (
	"errors"
	"github.com/Masterminds/log-go"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type keycloakTokenProcessor struct {
	JwksUrl string
}

func NewKeycloakTokenProcessor(jwksUrl string) TokenProcessor {
	return &keycloakTokenProcessor{jwksUrl}
}

var _jwks *keyfunc.JWKS

func (k keycloakTokenProcessor) getKeycloakCert() (jwks *keyfunc.JWKS, err error) {

	if _jwks != nil {
		return _jwks, nil
	}

	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Info("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err = keyfunc.Get(k.JwksUrl, options)
	if err != nil {
		log.Error("Can't load certs from ", k.JwksUrl)
		return nil, err
	}
	_jwks = jwks
	log.Info("Loaded certs from ", k.JwksUrl)

	return _jwks, nil
}

func (t keycloakTokenProcessor) Process(bearerToken string, r *http.Request) (valid bool, roles []string, err error) {

	jwks, err := t.getKeycloakCert()

	if err != nil {
		return false, nil, err
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(bearerToken, claims, jwks.Keyfunc)
	if err != nil {
		log.Error("Token validation error. ip: ", r.RemoteAddr, ", msg: ", err.Error())
		return false, nil, err
	}

	if !token.Valid {
		log.Error("The token is not valid.")
		return false, nil, errors.New("token is not valid")
	}

	tokenRoles := claims["realm_access"].(map[string]interface{})["roles"].([]interface{})

	roleArray := make([]string, len(tokenRoles))
	for i, v := range tokenRoles {
		roleArray[i] = v.(string)
	}

	return true, roleArray, nil
}
