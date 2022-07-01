package roleBuilder

type SecurityConfig interface {
	Config(security *RoleConfigurer)
}
