package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/dao"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (serv IndigoServiceServer) GetUserRoles(_ context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roleBindings, err := serv.Dao.GetUserRoleBindings(req.UserAccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user roles bindings: %v", err)
	}
	protoRoles := UserRoleBindingsToProtoRoles(serv.Dao, roleBindings)

	return &pb.GetUserRolesResponse{
		Roles: protoRoles,
	}, nil
}

// UserRoleBindingsToProtoRoles fetches a role for
// every binding and that way fills in the permissions as well.
func UserRoleBindingsToProtoRoles(da dao.DataAccessor, roleBindings []*model.UserRoleBinding) []*pb.Role {
	var protoRoles []*pb.Role
	for _, binding := range roleBindings {
		role, err := da.GetRole(model.IdToRoleIdentifier(binding.RoleId))
		if err != nil || role == nil {
			continue
		}

		protoRoles = append(protoRoles, role.ToProtoRole())
	}
	return protoRoles
}

func (serv IndigoServiceServer) AddUserRoles(_ context.Context, req *pb.AddUserRolesRequest) (*pb.AddUserRolesResponse, error) {
	var roleIds []string
	for _, id := range req.RoleIds {
		r, err := serv.Dao.GetRole(id)
		if err != nil {
			continue
		}
		roleIds = append(roleIds, r.Id)
	}
	if len(roleIds) == 0 {
		return nil, status.Error(codes.NotFound, "could not find any roles")
	}

	addedRoles, err := serv.Dao.AddUserRoles(req.UserAccountId, roleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.AddUserRolesResponse{
		AddedRoleIds: addedRoles,
	}, nil
}

func (serv IndigoServiceServer) RemoveUserRoles(_ context.Context, req *pb.RemoveUserRolesRequest) (*pb.RemoveUserRolesResponse, error) {
	var roleIds []string
	for _, id := range req.RoleIds {
		r, err := serv.Dao.GetRole(id)
		if err != nil {
			continue
		}
		roleIds = append(roleIds, r.Id)
	}
	if len(roleIds) == 0 {
		return nil, status.Error(codes.NotFound, "could not find any roles")
	}

	removedRoles, err := serv.Dao.RemoveUserRoles(req.UserAccountId, roleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.RemoveUserRolesResponse{
		RemovedRoleIds: removedRoles,
	}, nil
}
