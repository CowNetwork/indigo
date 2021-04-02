package psql

import (
	"fmt"
	"github.com/cownetwork/indigo/internal/model"
	"github.com/upper/db/v4"
)

type DataAccessor struct {
	Session db.Session
}

func (d *DataAccessor) InsertRole(role *model.Role) error {
	coll := d.Session.Collection("role_definitions")

	res, err := coll.Insert(role)
	if err != nil {
		return err
	}
	if res.ID() != nil {
		return fmt.Errorf("no psql id found after insertion of role %v", role.Id)
	}
	return nil
}

func (d *DataAccessor) GetRole(roleId string) (*model.Role, error) {
	coll := d.Session.Collection("role_definitions")
	res := coll.Find("id", roleId)

	var role model.Role
	err := res.One(&role)
	return &role, err
}

func (d *DataAccessor) GetRolePermissions(roleId string) ([]*model.RolePermissionBinding, error) {
	res := d.Session.Collection("role_permissions").Find("role_id", roleId)
	var bindings []*model.RolePermissionBinding

	err := res.All(&bindings)

	return bindings, err
}

func (d *DataAccessor) AddRolePermissions(roleId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("role_permissions")

	addedPerms := make([]string, len(permissions))
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

	filter := make([]db.LogicalExpr, len(permissions))
	for i, permission := range permissions {
		filter[i] = db.Cond{"permission": permission}
	}

	err := coll.Find("role_id", roleId).And(db.Or(filter...)).Delete()
	return permissions, err
}

func (d *DataAccessor) GetUserRoles(userAccountId string) ([]*model.Role, error) {
	coll := d.Session.Collection("user_roles")
	res := coll.Find("user_account_id", userAccountId)

	var roles []*model.Role
	err := res.All(roles)
	return roles, err
}

func (d *DataAccessor) AddUserRoles(userAccountId string, roleIds []string) ([]string, error) {
	coll := d.Session.Collection("user_roles")

	addedRoles := make([]string, len(roleIds))
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

	filter := make([]db.LogicalExpr, len(roleIds))
	for i, permission := range roleIds {
		filter[i] = db.Cond{"role_id": permission}
	}

	err := coll.Find("user_account_id", userAccountId).And(db.Or(filter...)).Delete()
	return roleIds, err
}

func (d *DataAccessor) AddUserPermissions(userAccountId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("user_permissions")

	addedPermissions := make([]string, len(permissions))
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
			addedPermissions = append(addedPermissions, permission)
		}
	}
	return addedPermissions, nil
}

func (d *DataAccessor) RemoveUserPermissions(userAccountId string, permissions []string) ([]string, error) {
	coll := d.Session.Collection("user_permissions")

	filter := make([]db.LogicalExpr, len(permissions))
	for i, permission := range permissions {
		filter[i] = db.Cond{"permission": permission}
	}

	err := coll.Find("user_account_id", userAccountId).And(db.Or(filter...)).Delete()
	return permissions, err
}
