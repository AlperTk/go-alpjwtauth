package securityConfig

import (
	roleAuth "JwtAuth/src/authorization/builder/roleBuilder"
)

type WebSecurityConfig struct {
}

func (s WebSecurityConfig) Config(security *roleAuth.RoleConfigurer) {
	security.
		//AnyRequest().PermitAll().
		AntMatcher("/api/v1/test").HasAnyRoles("Admin").
		AntMatcher("/api/v1/**").DenyAll()
}
