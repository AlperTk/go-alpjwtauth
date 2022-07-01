package accesscontrol

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/builder/roadBuillder"
	roleBuilder2 "github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/builder/roleBuilder"
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/model"
	"github.com/AlperTk/go-alpjwtauth/v2/errors"
	"net/http"
)

type basicRoleAuthorizer struct {
	securityConfig roleBuilder2.SecurityConfig
	requestRoad    *roadBuillder.RoadBuilder[authorization.RoleModel]
}

func NewBasicRoleAuthorizer(securityConfig roleBuilder2.SecurityConfig) AlpAuthorizer {
	instance := basicRoleAuthorizer{securityConfig: securityConfig, requestRoad: roadBuillder.NewRoadBuilder[authorization.RoleModel]()}
	instance.loadConfig()
	return &instance
}

func (b *basicRoleAuthorizer) ProcessUnauthorized(w http.ResponseWriter, r *http.Request) (defined bool, err error) {
	defined, err = processRequestRoadBeforeAuth(r, b.requestRoad)
	if defined && err != nil {
		responseUnauthorized(w)
	}
	return defined, err
}

func (b *basicRoleAuthorizer) ProcessAuthorized(roles []string, w http.ResponseWriter, r *http.Request) (defined bool, err error) {
	defined, err = processAuthorizedRequestRoad(roles, r, b.requestRoad)
	if err != nil {
		responseForbidden(w)
	}
	return defined, err
}

func processRequestRoadBeforeAuth(r *http.Request, requestRoad *roadBuillder.RoadBuilder[authorization.RoleModel]) (defined bool, err error) {
	securityDef, _ := requestRoad.Get(r.RequestURI)
	if securityDef == nil {
		return false, fmt.Errorf("no securityDef find")
	}

	if securityDef.Denied {
		return true, errors2.New("endpoint denied")
	}

	if securityDef.Permitted {
		return true, nil
	}

	return true, fmt.Errorf("not authorized")
}

func processAuthorizedRequestRoad(roles []string, r *http.Request, requestRoad *roadBuillder.RoadBuilder[authorization.RoleModel]) (defined bool, err error) {
	securityDef, _ := requestRoad.Get(r.RequestURI)
	if securityDef == nil {
		return false, errors2.New("no securityDef find")
	}

	if securityDef.Denied {
		return true, errors2.New("request denied")
	}

	if securityDef.Authenticated {
		return true, nil
	}

	if securityDef.Permitted {
		return true, nil
	}

	tokenRoles := make(map[string]bool)
	for _, role := range roles {
		tokenRoles[role] = true
	}

	for _, role := range securityDef.RoleMaps {
		res := tokenRoles[role]
		if res {
			return true, nil
		}
	}

	return true, fmt.Errorf("not authorized")
}

func (b *basicRoleAuthorizer) loadConfig() {
	roleConfigurer := roleBuilder2.RoleConfigurer{RequestRoad: b.requestRoad}
	b.securityConfig.Config(&roleConfigurer)
}

func responseForbidden(w http.ResponseWriter) {
	w.WriteHeader(403)
	_ = json.NewEncoder(w).Encode(errors.ForbiddenRequestError())
}

func responseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	_ = json.NewEncoder(w).Encode(errors.UnauthorizedError())
}
