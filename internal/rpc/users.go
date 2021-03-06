package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/model"
	"github.com/cownetwork/indigo/internal/perm"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (serv IndigoServiceServer) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	roleBindings, err := serv.Dao.GetUserRoleBindings(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user roles bindings: %v", err)
	}
	protoRoles := UserRoleBindingsToProtoRoles(serv.Dao, roleBindings)
	for _, role := range protoRoles {
		permBindings, err := serv.Dao.GetRolePermissions(role.Id)
		if err != nil {
			continue
		}
		rolePerms := make([]string, len(permBindings))
		for i, binding := range permBindings {
			rolePerms[i] = binding.Permission
		}
		role.Permissions = rolePerms
	}

	permBindings, err := serv.Dao.GetUserPermissions(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user permission bindings: %v", err)
	}
	perms := make([]string, len(permBindings))
	for i, binding := range permBindings {
		perms[i] = binding.Permission
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			AccountId:         req.UserAccountId,
			Roles:             protoRoles,
			CustomPermissions: perms,
		},
	}, nil
}

func (serv IndigoServiceServer) HasPermission(_ context.Context, req *pb.HasPermissionRequest) (*pb.HasPermissionResponse, error) {
	roleBindings, err := serv.Dao.GetUserRoleBindings(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user roles bindings: %v", err)
	}

	permBindings, err := serv.Dao.GetUserPermissions(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user permission bindings: %v", err)
	}

	var prevPrio int32
	v := perm.NewValidator([]string{})
	for _, binding := range roleBindings {
		r, err := serv.Dao.GetRole(model.ToRoleUuidIdentifier(binding.RoleId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
		}

		b, err := serv.Dao.GetRolePermissions(binding.RoleId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get role permissions: %v", err)
		}

		var rolePerms []string
		for _, permissionBinding := range b {
			rolePerms = append(rolePerms, permissionBinding.Permission)
		}

		v.Append(rolePerms, prevPrio < r.Priority)
		prevPrio = r.Priority
	}

	var perms []string
	for _, binding := range permBindings {
		perms = append(perms, binding.Permission)
	}
	validator := perm.NewValidator(perms)

	res := false
	for _, permission := range req.Permissions {
		if !validator.Validate(permission) {
			res = false
			break
		}
		res = true
	}

	return &pb.HasPermissionResponse{
		Result: res,
	}, nil
}
