package model

import pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"

type Role struct {
	Id        string `db:"id"`
	Name      string `db:"name"`
	Type      string `db:"type"`
	Priority  int32  `db:"priority"`
	Transient bool   `db:"transient"`
	Color     string `db:"color"`
}

func FromProtoRole(r *pb.Role) *Role {
	return &Role{
		Id:        r.Id,
		Name:      r.Name,
		Type:      r.Type,
		Priority:  r.Priority,
		Transient: r.Transient,
		Color:     r.Color,
	}
}

func (r *Role) ToProtoRole() *pb.Role {
	return &pb.Role{
		Id:        r.Id,
		Name:      r.Name,
		Type:      r.Type,
		Priority:  r.Priority,
		Transient: r.Transient,
		Color:     r.Color,
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

func IdToRoleIdentifier(roleId string) *pb.RoleIdentifier {
	return &pb.RoleIdentifier{
		Id: &pb.RoleIdentifier_Uuid{
			Uuid: roleId,
		},
	}
}
