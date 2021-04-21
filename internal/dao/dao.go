package dao

import (
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
)

type DataAccessor interface {
	ListRoles() ([]*model.Role, error)
	InsertRole(role *model.Role) error
	UpdateRole(roleId *pb.RoleIdentifier, role *model.Role) error
	GetRole(roleId *pb.RoleIdentifier) (*model.Role, error)
	DeleteRole(roleId *pb.RoleIdentifier) error
	GetRolePermissions(roleId string) ([]*model.RolePermissionBinding, error)
	AddRolePermissions(roleId string, permissions []string) ([]string, error)
	RemoveRolePermissions(roleId string, permissions []string) ([]string, error)
	GetUserRoleBindings(userAccountId string) ([]*model.UserRoleBinding, error)
	AddUserRoles(userAccountId string, roleIds []string) ([]string, error)
	RemoveUserRoles(userAccountId string, roleIds []string) ([]string, error)
	GetUserPermissions(userAccountId string) ([]*model.UserPermissionBinding, error)
	AddUserPermissions(userAccountId string, permissions []string) ([]string, error)
	RemoveUserPermissions(userAccountId string, permissions []string) ([]string, error)
}
