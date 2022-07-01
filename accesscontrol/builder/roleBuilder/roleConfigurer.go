package roleBuilder

import (
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/builder/roadBuillder"
	"github.com/AlperTk/go-alpjwtauth/v2/accesscontrol/model"
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
