package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/perm"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (serv IndigoServiceServer) AddRolePermissions(_ context.Context, req *pb.AddRolePermissionsRequest) (*pb.AddRolePermissionsResponse, error) {
	role, err := serv.Dao.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role == nil {
		return nil, status.Error(codes.NotFound, "this role does not exists")
	}

	// only take those permissions that match the regex.
	var perms []string
	for _, permission := range req.Permissions {
		if perm.ValidatePermission(permission) {
			perms = append(perms, permission)
		}
	}

	addedPerms, err := serv.Dao.AddRolePermissions(req.RoleId, perms)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add permissions: %v", err)
	}

	return &pb.AddRolePermissionsResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv IndigoServiceServer) RemoveRolePermissions(_ context.Context, req *pb.RemoveRolePermissionsRequest) (*pb.RemoveRolePermissionsResponse, error) {
	role, err := serv.Dao.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role == nil {
		return nil, status.Error(codes.NotFound, "this role does not exists")
	}

	removedPerms, err := serv.Dao.RemoveRolePermissions(req.RoleId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove permissions: %v", err)
	}

	return &pb.RemoveRolePermissionsResponse{
		RemovedPermissions: removedPerms,
	}, nil
}
