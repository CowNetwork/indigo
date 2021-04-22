package model

import (
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/thoas/go-funk"
)

type Role struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Type        string `db:"type"`
	Priority    int32  `db:"priority"`
	Transient   bool   `db:"transient"`
	Color       string `db:"color"`
	Permissions []string
}

func FromProtoRole(r *pb.Role) *Role {
	return &Role{
		Id:          r.Id,
		Name:        r.Name,
		Type:        r.Type,
		Priority:    r.Priority,
		Transient:   r.Transient,
		Color:       r.Color,
		Permissions: r.Permissions,
	}
}

func (r *Role) ToProtoRole() *pb.Role {
	return &pb.Role{
		Id:          r.Id,
		Name:        r.Name,
		Type:        r.Type,
		Priority:    r.Priority,
		Transient:   r.Transient,
		Color:       r.Color,
		Permissions: r.Permissions,
	}
}

func (r *Role) SetPermissions(perms []*RolePermissionBinding) {
	r.Permissions = []string{}
	for _, perm := range perms {
		r.Permissions = append(r.Permissions, perm.Permission)
	}
}

func (r *Role) AddPermissions(perms []string) {
	for _, perm := range perms {
		r.Permissions = append(r.Permissions, perm)
	}
}

func (r *Role) RemovePermissions(perms []string) {
	r.Permissions = funk.SubtractString(r.Permissions, perms)
}

func (r *Role) Merge(r2 *pb.Role, fm pb.UpdateRoleRequest_FieldMask) {
	switch fm {
	case pb.UpdateRoleRequest_FIELD_MASK_ALL:
	case pb.UpdateRoleRequest_FIELD_MASK_ALL_PROPERTIES:

		break
	case pb.UpdateRoleRequest_FIELD_MASK_NAME:
		r.Name = r2.Name
		break
	case pb.UpdateRoleRequest_FIELD_MASK_TYPE:
		r.Type = r2.Type
		break
	case pb.UpdateRoleRequest_FIELD_MASK_PRIORITY:
		r.Priority = r2.Priority
		break
	case pb.UpdateRoleRequest_FIELD_MASK_TRANSIENCE:
		r.Transient = r2.Transient
		break
	case pb.UpdateRoleRequest_FIELD_MASK_COLOR:
		r.Color = r2.Color
		break
	case pb.UpdateRoleRequest_FIELD_MASK_PERMISSIONS:
		r.Permissions = r2.Permissions
		break
	}
}

type RolePermissionBinding struct {
	RoleId     string `db:"role_id"`
	Permission string `db:"permission"`
}

type UserRoleBinding struct {
	UserAccountId string `db:"user_account_id"`
	RoleId        string `db:"role_id"`
}

type UserPermissionBinding struct {
	UserAccountId string `db:"user_account_id"`
	Permission    string `db:"permission"`
}

func ToRoleUuidIdentifier(roleId string) *pb.RoleIdentifier {
	return &pb.RoleIdentifier{
		Id: &pb.RoleIdentifier_Uuid{
			Uuid: roleId,
		},
	}
}

func ToRoleNameIdentifier(name string, t string) *pb.RoleIdentifier {
	return &pb.RoleIdentifier{
		Id: &pb.RoleIdentifier_NameId{NameId: &pb.RoleNameIdentifier{
			Name: name,
			Type: t,
		}},
	}
}
