package model

type Role struct {
	Id      int
	Service string
	Role    string
}

type RoleOperations interface {
	GetRoleByMethod(string) (string, error)
	SaveRole(Role) (int, error)
}
