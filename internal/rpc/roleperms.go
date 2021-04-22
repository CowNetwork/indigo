package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/eventhandler"
	"github.com/cownetwork/indigo/internal/perm"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/thoas/go-funk"
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

	bindings, err := serv.Dao.GetRolePermissions(role.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role permissions: %v", err)
	}
	role.SetPermissions(bindings)

	// only take those permissions that match the regex.
	perms := funk.FilterString(req.Permissions, func(s string) bool {
		return perm.ValidatePermission(s)
	})

	addedPerms, err := serv.Dao.AddRolePermissions(role.Id, perms)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add permissions: %v", err)
	}
	role.AddPermissions(addedPerms)

	eventhandler.SendRoleUpdateEvent(role.ToProtoRole(), pb.RoleUpdateEvent_ACTION_UPDATED)

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

	bindings, err := serv.Dao.GetRolePermissions(role.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role permissions: %v", err)
	}
	role.SetPermissions(bindings)

	removedPerms, err := serv.Dao.RemoveRolePermissions(role.Id, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove permissions: %v", err)
	}
	role.RemovePermissions(removedPerms)

	eventhandler.SendRoleUpdateEvent(role.ToProtoRole(), pb.RoleUpdateEvent_ACTION_UPDATED)

	return &pb.RemoveRolePermissionsResponse{
		RemovedPermissions: removedPerms,
	}, nil
}
