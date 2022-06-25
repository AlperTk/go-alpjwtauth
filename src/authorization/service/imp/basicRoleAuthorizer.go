package authorization

import (
	"JwtAuth/src/authorization/builder/roadBuillder"
	roleBuilder2 "JwtAuth/src/authorization/builder/roleBuilder"
	authorization3 "JwtAuth/src/authorization/model"
	"JwtAuth/src/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

type basicRoleAuthorizer struct {
	securityConfig roleBuilder2.SecurityConfig
	requestRoad    *roadBuillder.RoadBuilder[authorization3.RoleModel]
}

func NewBasicRoleAuthorizer(securityConfig roleBuilder2.SecurityConfig) *basicRoleAuthorizer {
	instance := basicRoleAuthorizer{securityConfig: securityConfig, requestRoad: roadBuillder.NewRoadBuilder[authorization3.RoleModel]()}
	instance.loadConfig()
	return &instance
}

func (b *basicRoleAuthorizer) Process(roles []string, w http.ResponseWriter, r *http.Request) (defined bool, err error) {
	defined, err = proccessRequestRoad(roles, r, b.requestRoad)
	if err != nil {
		responseNotAuthorized(w)
	}
	return defined, err
}

func proccessRequestRoad(roles []string, r *http.Request, requestRoad *roadBuillder.RoadBuilder[authorization3.RoleModel]) (defined bool, err error) {
	securityDef, _ := requestRoad.Get(r.RequestURI)
	if securityDef == nil {
		return false, fmt.Errorf("no securityDef find")
	}

	if securityDef.Authenticated {
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

func responseNotAuthorized(w http.ResponseWriter) {
	w.WriteHeader(403)
	_ = json.NewEncoder(w).Encode(errors.ForbiddenRequestError())
}
