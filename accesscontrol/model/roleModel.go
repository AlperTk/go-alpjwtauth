package authorization

type RoleModel struct {
	RoleMaps      []string
	Authenticated bool
	Denied        bool
	Permitted     bool
}
