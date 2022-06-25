package securityConfig

import (
	"github.com/AlperTk/go-jwt-role-based-auth/src/authorization/builder/roleBuilder"
)

type WebSecurityConfig struct {
}

func (s WebSecurityConfig) Config(security *roleBuilder.RoleConfigurer) {
	security.
		//AntMatcher("/api/v1/**").PermitAll().
		//AntMatcher("/api2/v1/test").HasAnyRoles("Admin").
		//AntMatcher("/api3/v1/**").DenyAll().
		AnyRequest().DenyAll()
}
