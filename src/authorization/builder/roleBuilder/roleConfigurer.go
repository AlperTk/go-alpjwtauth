package roleBuilder

import (
	"github.com/AlperTk/go-jwt-role-based-auth/src/authorization/builder/roadBuillder"
	authorization "github.com/AlperTk/go-jwt-role-based-auth/src/authorization/model"
)

type RoleConfigurer struct {
	RequestRoad  *roadBuillder.RoadBuilder[authorization.RoleModel]
	levelCounter int
}

func (r RoleConfigurer) AntMatcher(endpoints ...string) *RoleBuilder {
	var roleBuilder = &RoleBuilder{
		Endpoints:   endpoints,
		Caller:      &r,
		requestRoad: r.RequestRoad,
	}
	r.levelCounter++
	return roleBuilder
}

func (r RoleConfigurer) AnyRequest() *RoleBuilder {
	var roleBuilder = &RoleBuilder{
		Endpoints:   []string{"/**"},
		Caller:      &r,
		requestRoad: r.RequestRoad,
	}
	r.levelCounter++
	return roleBuilder
}
