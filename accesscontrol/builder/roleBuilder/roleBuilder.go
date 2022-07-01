package roleBuilder

import (
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/builder/roadBuillder"
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/model"
	"github.com/Masterminds/log-go"
)

type RoleBuilder struct {
	Endpoints   []string
	Caller      *RoleConfigurer
	requestRoad *roadBuillder.RoadBuilder[authorization.RoleModel]
}

func createIfNotExist(roleRoad *roadBuillder.RoadBuilder[authorization.RoleModel], endpoint string) *authorization.RoleModel {
	data := &authorization.RoleModel{
		Authenticated: false,
		Denied:        false,
		Permitted:     false,
		RoleMaps:      []string{},
	}

	err := roleRoad.Put(endpoint, data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (r *RoleBuilder) HasAnyRoles(roles ...string) *RoleConfigurer {
	for _, endpoint := range r.Endpoints {
		roleModel := createIfNotExist(r.requestRoad, endpoint)
		for _, role := range roles {
			roleModel.RoleMaps = append(roleModel.RoleMaps, role)
		}
	}
	return r.Caller
}

func (r *RoleBuilder) Authenticated() *RoleConfigurer {
	for _, endpoint := range r.Endpoints {
		roleModel := createIfNotExist(r.requestRoad, endpoint)
		roleModel.Authenticated = true
	}
	return r.Caller
}

func (r *RoleBuilder) PermitAll() *RoleConfigurer {
	for _, endpoint := range r.Endpoints {
		roleModel := createIfNotExist(r.requestRoad, endpoint)
		roleModel.Permitted = true
	}
	return r.Caller
}

func (r *RoleBuilder) DenyAll() *RoleConfigurer {
	for _, endpoint := range r.Endpoints {
		roleModel := createIfNotExist(r.requestRoad, endpoint)
		roleModel.Denied = true
	}
	return r.Caller
}
