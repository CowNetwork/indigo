package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

var snakeCaseRegex = regexp.MustCompile("^[a-z]+(_[a-z]+)*$")

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

	bindings, err := serv.Dao.GetRolePermissions(role.Id)
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
	n := req.Role.Name
	t := req.Role.Type
	if len(n) == 0 || len(t) == 0 {
		return nil, status.Error(codes.InvalidArgument, "name or type can not be empty.")
	}
	if !snakeCaseRegex.MatchString(n) || !snakeCaseRegex.MatchString(t) {
		return nil, status.Error(codes.InvalidArgument, "name and type must be snake_case.")
	}
	if len(req.Role.Color) > 6 {
		return nil, status.Error(codes.InvalidArgument, "role color length must be 6 or less.")
	}

	role, err := serv.Dao.GetRole(&pb.RoleIdentifier{
		Id: &pb.RoleIdentifier_NameId{NameId: &pb.RoleNameIdentifier{
			Name: n,
			Type: t,
		}},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role != nil {
		return nil, status.Error(codes.AlreadyExists, "this role already exists")
	}

	roleUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not generate uuid: %v", err)
	}
	req.Role.Id = roleUuid.String()

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
	// if it is name+type, then the user may wanted to update
	// the name, so we only set the uuid for them.
	switch u := req.RoleId.Id.(type) {
	case *pb.RoleIdentifier_Uuid:
		role.Id = u.Uuid
	}

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
