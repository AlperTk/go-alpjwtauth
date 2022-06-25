package roleBuilder

import (
	"JwtAuth/src/authorization/builder/roadBuillder"
	authorization3 "JwtAuth/src/authorization/model"
)

type RoleConfigurer struct {
	RequestRoad  *roadBuillder.RoadBuilder[authorization3.RoleModel]
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
