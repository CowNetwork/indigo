package dao

import "github.com/cownetwork/indigo/internal/model"

type DataAccessor interface {
	ListRoles() ([]*model.Role, error)
	InsertRole(role *model.Role) error
	UpdateRole(roleId string, role *model.Role) error
	GetRole(roleId string) (*model.Role, error)
	DeleteRole(roleId string) error
	GetRolePermissions(roleId string) ([]*model.RolePermissionBinding, error)
	AddRolePermissions(roleId string, permissions []string) ([]string, error)
	RemoveRolePermissions(roleId string, permissions []string) ([]string, error)
	GetUserRoleBindings(userAccountId string) ([]*model.UserRoleBinding, error)
	AddUserRoles(userAccountId string, roleIds []string) ([]string, error)
	RemoveUserRoles(userAccountId string, roleIds []string) ([]string, error)
	AddUserPermissions(userAccountId string, permissions []string) ([]string, error)
	RemoveUserPermissions(userAccountId string, permissions []string) ([]string, error)
}
