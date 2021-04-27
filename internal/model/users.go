package model

import (
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/thoas/go-funk"
)

type User struct {
	AccountId         string
	Roles             []string
	CustomPermissions []string
}

type UserRoleBinding struct {
	UserAccountId string `db:"user_account_id"`
	RoleId        string `db:"role_id"`
}

type UserPermissionBinding struct {
	UserAccountId string `db:"user_account_id"`
	Permission    string `db:"permission"`
}

func NewUser(accountId string) *User {
	return &User{
		AccountId: accountId,
	}
}

func (u *User) SetPermissions(perms []*UserPermissionBinding) {
	u.CustomPermissions = []string{}
	for _, perm := range perms {
		u.CustomPermissions = append(u.CustomPermissions, perm.Permission)
	}
}

func (u *User) AddPermissions(perms []string) {
	u.CustomPermissions = append(u.CustomPermissions, perms...)
}

func (u *User) RemovePermissions(perms []string) {
	u.CustomPermissions = funk.SubtractString(u.CustomPermissions, perms)
}

func (u *User) SetRoles(bindings []*UserRoleBinding) {
	u.Roles = []string{}
	for _, perm := range bindings {
		u.Roles = append(u.Roles, perm.RoleId)
	}
}

func (u *User) AddRoles(roles []string) {
	u.Roles = append(u.Roles, roles...)
}

func (u *User) RemoveRoles(roles []string) {
	u.Roles = funk.SubtractString(u.Roles, roles)
}

func (u *User) ToProtoUser() *pb.User {
	var protoRoles []*pb.Role
	for _, roleId := range u.Roles {
		protoRoles = append(protoRoles, &pb.Role{Id: roleId})
	}

	return &pb.User{
		AccountId:         u.AccountId,
		Roles:             protoRoles,
		CustomPermissions: u.CustomPermissions,
	}
}
