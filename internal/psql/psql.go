package psql

import (
	"fmt"
	"github.com/cownetwork/indigo/internal/model"
	"github.com/upper/db/v4"
)

type DataAccessor struct {
	Session db.Session
}

func (d *DataAccessor) ListRoles() ([]*model.Role, error) {
	coll := d.Session.Collection("role_definitions")
	res := coll.Find()

	var roles []*model.Role
	err := res.All(&roles)
	return roles, err
}

func (d *DataAccessor) InsertRole(role *model.Role) error {
	coll := d.Session.Collection("role_definitions")

	res, err := coll.Insert(role)
	if err != nil {
		return err
	}
	if res.ID() == nil {
		return fmt.Errorf("no psql id found after insertion of role %v", role.Id)
	}
	return nil
}

func (d *DataAccessor) UpdateRole(roleId string, role *model.Role) error {
	coll := d.Session.Collection("role_definitions")
	role.Id = roleId

	return coll.UpdateReturning(role)
}

func (d *DataAccessor) GetRole(roleId string) (*model.Role, error) {
	coll := d.Session.Collection("role_definitions")
	res := coll.Find("id", roleId)

	count, err := res.TotalEntries()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	var role model.Role
	err = res.One(&role)
	return &role, err
}

func (d *DataAccessor) DeleteRole(roleId string) error {
	coll := d.Session.Collection("role_definitions")
	err := coll.Find("id", roleId).Delete()
	if err != nil {
		return err
	}

	coll = d.Session.Collection("user_roles")
	err = coll.Find("role_id", roleId).Delete()
	if err != nil {
		return err
	}

	coll = d.Session.Collection("role_permissions")
	return coll.Find("role_id", roleId).Delete()
}

func (d *DataAccessor) GetRolePermissions(roleId string) ([]*model.RolePermissionBinding, error) {
	res := d.Session.Collection("role_permissions").Find("role_id", roleId)
	var bindings []*model.RolePermissionBinding

	err := res.All(&bindings)

	return bindings, err
}

func (d *DataAccessor) AddRolePermissions(roleId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("role_permissions")

	var addedPerms []string
	for _, perm := range permissions {
		binding := model.RolePermissionBinding{
			RoleId:     roleId,
			Permission: perm,
		}

		exists, _ := coll.Find("role_id", roleId).And("permission", perm).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(&binding)
		if err == nil && res.ID() != nil {
			addedPerms = append(addedPerms, perm)
		}
	}
	return addedPerms, nil
}

func (d *DataAccessor) RemoveRolePermissions(roleId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("role_permissions")

	var removedPerms []string
	for _, perm := range permissions {
		binding := model.RolePermissionBinding{
			RoleId:     roleId,
			Permission: perm,
		}

		exists, _ := coll.Find("role_id", roleId).And("permission", perm).Exists()
		if !exists {
			continue
		}

		err := coll.Find(&binding).Delete()
		if err == nil {
			removedPerms = append(removedPerms, perm)
		}
	}
	return removedPerms, nil
}

func (d *DataAccessor) GetUserRoleBindings(userAccountId string) ([]*model.UserRoleBinding, error) {
	coll := d.Session.Collection("user_roles")
	res := coll.Find("user_account_id", userAccountId)

	var roleBindings []*model.UserRoleBinding
	err := res.All(&roleBindings)

	return roleBindings, err
}

func (d *DataAccessor) AddUserRoles(userAccountId string, roleIds []string) ([]string, error) {
	coll := d.Session.Collection("user_roles")

	var addedRoles []string
	for _, id := range roleIds {
		binding := model.UserRoleBinding{
			UserAccountId: userAccountId,
			RoleId:        id,
		}

		exists, _ := coll.Find("user_account_id", userAccountId).And("role_id", id).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(binding)
		if err == nil && res.ID() != nil {
			addedRoles = append(addedRoles, id)
		}
	}
	return addedRoles, nil
}

func (d *DataAccessor) RemoveUserRoles(userAccountId string, roleIds []string) ([]string, error) {
	coll := d.Session.Collection("user_roles")

	var removedRoles []string
	for _, id := range roleIds {
		binding := model.UserRoleBinding{
			UserAccountId: userAccountId,
			RoleId:        id,
		}

		exists, _ := coll.Find("user_account_id", userAccountId).And("role_id", id).Exists()
		if exists {
			continue
		}

		err := coll.Find(&binding).Delete()
		if err == nil {
			removedRoles = append(removedRoles, id)
		}
	}
	return removedRoles, nil
}

func (d *DataAccessor) GetUserPermissions(userAccountId string) ([]*model.UserPermissionBinding, error) {
	coll := d.Session.Collection("user_permissions")

	res := coll.Find("user_account_id", userAccountId)

	var permBindings []*model.UserPermissionBinding
	err := res.All(&permBindings)

	return permBindings, err
}

func (d *DataAccessor) AddUserPermissions(userAccountId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("user_permissions")

	var addedPerms []string
	for _, permission := range permissions {
		binding := model.UserPermissionBinding{
			UserAccountId: userAccountId,
			Permission:    permission,
		}

		exists, _ := coll.Find("user_account_id", userAccountId).And("permission", permission).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(binding)
		if err == nil && res.ID() != nil {
			addedPerms = append(addedPerms, permission)
		}
	}
	return addedPerms, nil
}

func (d *DataAccessor) RemoveUserPermissions(userAccountId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("user_permissions")

	var removedPerms []string
	for _, permission := range permissions {
		binding := model.UserPermissionBinding{
			UserAccountId: userAccountId,
			Permission:    permission,
		}

		exists, _ := coll.Find("user_account_id", userAccountId).And("permission", permission).Exists()
		if exists {
			continue
		}

		err := coll.Find(&binding).Delete()
		if err == nil {
			removedPerms = append(removedPerms, permission)
		}
	}
	return removedPerms, nil
}
