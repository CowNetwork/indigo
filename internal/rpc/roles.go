package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	if role == nil {
		return nil, status.Errorf(codes.NotFound, "could not find role")
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
	role, err := serv.Dao.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role == nil {
		return nil, status.Error(codes.NotFound, "this role does not exists")
	}

	err = serv.Dao.DeleteRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete role: %v", err)
	}
	return &pb.DeleteRoleResponse{}, nil
}
