package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/dao"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IndigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	Dao dao.DataAccessor
}

func (serv IndigoServiceServer) ListRoles(_ context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	roles, err := serv.Dao.ListRoles()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list roles: %v", err)
	}

	var protoRoles []*pb.Role
	for _, role := range roles {
		protoRoles = append(protoRoles, role.ToProtoRole())
	}

	for _, role := range protoRoles {
		perms, err := serv.Dao.GetRolePermissions(role.Id)
		if err != nil {
			continue
		}

		for _, perm := range perms {
			role.Permissions = append(role.Permissions, perm.Permission)
		}
	}

	return &pb.ListRolesResponse{
		Roles: protoRoles,
	}, nil
}

func (serv IndigoServiceServer) GetRole(_ context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	role, err := serv.Dao.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}

	bindings, err := serv.Dao.GetRolePermissions(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role permissions: %v", err)
	}

	permissions := make([]string, len(bindings))
	for i, binding := range bindings {
		permissions[i] = binding.Permission
	}

	r := role.ToProtoRole()
	r.Permissions = permissions

	return &pb.GetRoleResponse{Role: r}, nil
}

func (serv IndigoServiceServer) InsertRole(_ context.Context, req *pb.InsertRoleRequest) (*pb.InsertRoleResponse, error) {
	role, err := serv.Dao.GetRole(req.Role.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role != nil {
		return nil, status.Error(codes.AlreadyExists, "this role already exists")
	}

	err = serv.Dao.InsertRole(model.FromProtoRole(req.Role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	return &pb.InsertRoleResponse{
		InsertedRole: req.Role,
	}, nil
}

func (serv IndigoServiceServer) UpdateRole(_ context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	err := serv.Dao.UpdateRole(req.RoleId, model.FromProtoRole(req.RoleData))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update role: %v", err)
	}

	role := req.RoleData
	role.Id = req.RoleId
	return &pb.UpdateRoleResponse{
		UpdatedRole: role,
	}, nil
}

func (serv IndigoServiceServer) DeleteRole(_ context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	err := serv.Dao.DeleteRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete role: %v", err)
	}
	return &pb.DeleteRoleResponse{}, nil
}

func (serv IndigoServiceServer) AddRolePermission(_ context.Context, req *pb.AddRolePermissionRequest) (*pb.AddRolePermissionResponse, error) {
	addedPerms, err := serv.Dao.AddRolePermissions(req.RoleId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add permissions: %v", err)
	}

	return &pb.AddRolePermissionResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv IndigoServiceServer) RemoveRolePermission(_ context.Context, req *pb.RemoveRolePermissionRequest) (*pb.RemoveRolePermissionResponse, error) {
	removedPerms, err := serv.Dao.RemoveRolePermissions(req.RoleId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove permissions: %v", err)
	}

	return &pb.RemoveRolePermissionResponse{
		RemovedPermissions: removedPerms,
	}, nil
}

func (serv IndigoServiceServer) GetUserRoles(_ context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roleBindings, err := serv.Dao.GetUserRoleBindings(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user roles bindings: %v", err)
	}

	var protoRoles []*pb.Role
	for _, binding := range roleBindings {
		role, err := serv.Dao.GetRole(binding.RoleId)
		if err != nil {
			continue
		}

		protoRoles = append(protoRoles, role.ToProtoRole())
	}

	return &pb.GetUserRolesResponse{
		Roles: protoRoles,
	}, nil
}

func (serv IndigoServiceServer) AddUserRole(_ context.Context, req *pb.AddUserRoleRequest) (*pb.AddUserRoleResponse, error) {
	addedRoles, err := serv.Dao.AddUserRoles(req.UserAccountId, req.RoleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.AddUserRoleResponse{
		AddedRoleIds: addedRoles,
	}, nil
}

func (serv IndigoServiceServer) RemoveUserRole(_ context.Context, req *pb.RemoveUserRoleRequest) (*pb.RemoveUserRoleResponse, error) {
	removedRoles, err := serv.Dao.RemoveUserRoles(req.UserAccountId, req.RoleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.RemoveUserRoleResponse{
		RemovedRoleIds: removedRoles,
	}, nil
}

func (serv IndigoServiceServer) AddUserPermission(_ context.Context, req *pb.AddUserPermissionRequest) (*pb.AddUserPermissionResponse, error) {
	addedPerms, err := serv.Dao.AddUserPermissions(req.UserAccountId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user permissions: %v", err)
	}

	return &pb.AddUserPermissionResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv IndigoServiceServer) RemoveUserPermission(_ context.Context, req *pb.RemoveUserPermissionRequest) (*pb.RemoveUserPermissionResponse, error) {
	removedPerms, err := serv.Dao.RemoveUserPermissions(req.UserAccountId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove user permissions: %v", err)
	}

	return &pb.RemoveUserPermissionResponse{
		RemovedPermissions: removedPerms,
	}, nil
}
