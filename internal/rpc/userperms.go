package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/eventhandler"
	"github.com/cownetwork/indigo/internal/model"
	"github.com/cownetwork/indigo/internal/perm"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (serv IndigoServiceServer) GetUserPermissions(_ context.Context, req *pb.GetUserPermissionsRequest) (*pb.GetUserPermissionsResponse, error) {
	permBindings, err := serv.Dao.GetUserPermissions(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user permissions: %v", err)
	}

	perms := make([]string, len(permBindings))
	for i, binding := range permBindings {
		perms[i] = binding.Permission
	}

	return &pb.GetUserPermissionsResponse{Permissions: perms}, nil
}

func (serv IndigoServiceServer) AddUserPermissions(_ context.Context, req *pb.AddUserPermissionsRequest) (*pb.AddUserPermissionsResponse, error) {
	user := model.NewUser(req.UserAccountId)
	permBindings, err := serv.Dao.GetUserPermissions(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user permissions: %v", err)
	}
	user.SetPermissions(permBindings)

	// only take those permissions that match the regex.
	var perms []string
	for _, permission := range req.Permissions {
		if perm.ValidatePermission(permission) {
			perms = append(perms, permission)
		}
	}

	addedPerms, err := serv.Dao.AddUserPermissions(req.UserAccountId, perms)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user permissions: %v", err)
	}
	user.AddPermissions(addedPerms)

	eventhandler.SendUserPermUpdateEvent(user.ToProtoUser(), pb.UserPermissionUpdateEvent_ACTION_PERM_ADDED)

	return &pb.AddUserPermissionsResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv IndigoServiceServer) RemoveUserPermissions(_ context.Context, req *pb.RemoveUserPermissionsRequest) (*pb.RemoveUserPermissionsResponse, error) {
	user := model.NewUser(req.UserAccountId)
	permBindings, err := serv.Dao.GetUserPermissions(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user permissions: %v", err)
	}
	user.SetPermissions(permBindings)

	removedPerms, err := serv.Dao.RemoveUserPermissions(req.UserAccountId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove user permissions: %v", err)
	}
	user.RemovePermissions(removedPerms)

	eventhandler.SendUserPermUpdateEvent(user.ToProtoUser(), pb.UserPermissionUpdateEvent_ACTION_PERM_REMOVED)

	return &pb.RemoveUserPermissionsResponse{
		RemovedPermissions: removedPerms,
	}, nil
}
