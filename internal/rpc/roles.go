package rpc

import (
	"context"
	"github.com/cownetwork/indigo/internal/eventhandler"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (serv IndigoServiceServer) ListRoles(_ context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	roles, err := serv.Dao.ListRoles()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list roles: %v", err)
	}

	for _, role := range roles {
		perms, err := serv.Dao.GetRolePermissions(role.Id)
		if err != nil {
			continue
		}

		role.SetPermissions(perms)
	}

	var protoRoles []*pb.Role
	for _, role := range roles {
		protoRoles = append(protoRoles, role.ToProtoRole())
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
	role.SetPermissions(bindings)

	return &pb.GetRoleResponse{
		Role: role.ToProtoRole(),
	}, nil
}

func (serv IndigoServiceServer) InsertRole(_ context.Context, req *pb.InsertRoleRequest) (*pb.InsertRoleResponse, error) {
	role := model.FromProtoRole(req.Role)

	err := ValidateRole(role)
	if err != nil {
		return nil, err
	}

	r, err := serv.Dao.GetRole(model.ToRoleNameIdentifier(role.Name, role.Type))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if r != nil {
		return nil, status.Error(codes.AlreadyExists, "this role already exists")
	}

	roleUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not generate uuid: %v", err)
	}
	role.Id = roleUuid.String()

	err = serv.Dao.InsertRole(role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	if len(req.Role.Permissions) > 0 {
		_, err = serv.Dao.AddRolePermissions(role.Id, role.Permissions)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not initialize role permissions: %v", err)
		}
	}

	pr := role.ToProtoRole()
	eventhandler.SendRoleUpdateEvent(pr, pb.RoleUpdateEvent_ACTION_ADDED)

	return &pb.InsertRoleResponse{
		InsertedRole: pr,
	}, nil
}

func (serv IndigoServiceServer) UpdateRole(_ context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	role, err := serv.Dao.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}
	if role == nil {
		return nil, status.Errorf(codes.NotFound, "could not find role")
	}

	prevPerms := append([]string(nil), role.Permissions...)
	for _, mask := range req.FieldMasks {
		role.Merge(req.RoleData, mask)
	}
	removed, added := funk.DifferenceString(prevPerms, role.Permissions)

	err = ValidateRole(role)
	if err != nil {
		return nil, err
	}

	err = serv.Dao.UpdateRole(req.RoleId, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update role: %v", err)
	}

	if funk.Contains(req.FieldMasks, pb.UpdateRoleRequest_FIELD_MASK_ALL) ||
		funk.Contains(req.FieldMasks, pb.UpdateRoleRequest_FIELD_MASK_PERMISSIONS) {
		_, err = serv.Dao.AddRolePermissions(role.Id, added)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not update role permissions: %v", err)
		}

		_, err = serv.Dao.RemoveRolePermissions(role.Id, removed)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not update role permissions: %v", err)
		}
	}

	r := role.ToProtoRole()
	eventhandler.SendRoleUpdateEvent(r, pb.RoleUpdateEvent_ACTION_UPDATED)

	return &pb.UpdateRoleResponse{
		UpdatedRole: r,
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

	eventhandler.SendRoleUpdateEvent(role.ToProtoRole(), pb.RoleUpdateEvent_ACTION_DELETED)

	return &pb.DeleteRoleResponse{}, nil
}
