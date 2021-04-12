package dao

import "github.com/cownetwork/indigo/internal/model"

type DataAccessor interface {
	InsertRole(role *model.Role) error
	GetRole(roleId string) (*model.Role, error)
	GetRolePermissions(roleId string) ([]*model.RolePermissionBinding, error)
	AddRolePermissions(roleId string, permissions []string) ([]string, error)
	RemoveRolePermissions(roleId string, permissions []string) ([]string, error)
	GetUserRoleBindings(userAccountId string) ([]*model.UserRoleBinding, error)
	AddUserRoles(userAccountId string, roleIds []string) ([]string, error)
	RemoveUserRoles(userAccountId string, roleIds []string) ([]string, error)
	AddUserPermissions(userAccountId string, permissions []string) ([]string, error)
	RemoveUserPermissions(userAccountId string, permissions []string) ([]string, error)
}
