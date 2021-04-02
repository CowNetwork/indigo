package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/dao"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/indigo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IndigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	Dao dao.DataAccessor
}

func (serv IndigoServiceServer) GetRole(_ context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
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

	return r, nil
}

func (serv IndigoServiceServer) InsertRole(_ context.Context, role *pb.Role) (*pb.InsertRoleResponse, error) {
	err := serv.Dao.InsertRole(model.FromProtoRole(role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	return &pb.InsertRoleResponse{
		Successful: true,
	}, nil
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
	roles, err := serv.Dao.GetUserRoles(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user roles: %v", err)
	}

	var protoRoles []*pb.Role
	for _, role := range roles {
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
