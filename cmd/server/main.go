package main

import (
	"context"
	"github.com/cownetwork/indigo/internal/dao"
	"github.com/cownetwork/indigo/internal/model"
	"github.com/cownetwork/indigo/internal/psql"
	"github.com/upper/db/v4/adapter/postgresql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"

	pb "github.com/cownetwork/indigo/proto"
)

func main() {
	log.Println("Hello World!")

	connUrl := &postgresql.ConnectionURL{
		Host:     getEnvOrDefault("POSTGRES_URL", "localhost:5432"),
		User:     getEnvOrDefault("POSTGRES_USER", "test"),
		Password: getEnvOrDefault("POSTGRES_PASSWORD", "password"),
		Database: getEnvOrDefault("POSTGRES_DB", "test_database"),
	}

	log.Printf("Connecting to PostgresSQL at %s ...", connUrl.Host)

	sess, err := postgresql.Open(connUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer sess.Close()

	log.Println("Connected to PostgresSQL.")

	// setup grpc server
	lis, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	s.RegisterService(&pb.IndigoService_ServiceDesc, &indigoServiceServer{
		da: &psql.DataAccessor{Session: sess},
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnvOrDefault(env string, def string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return def
	}
	return value
}

type indigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	da dao.DataAccessor
}

func (serv indigoServiceServer) GetRole(_ context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	role, err := serv.da.GetRole(req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get role: %v", err)
	}

	bindings, err := serv.da.GetRolePermissions(req.RoleId)
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

func (serv indigoServiceServer) InsertRole(_ context.Context, role *pb.Role) (*pb.InsertRoleResponse, error) {
	err := serv.da.InsertRole(model.FromProtoRole(role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	return &pb.InsertRoleResponse{
		Successful: true,
	}, nil
}

func (serv indigoServiceServer) AddRolePermission(_ context.Context, req *pb.AddRolePermissionRequest) (*pb.AddRolePermissionResponse, error) {
	addedPerms, err := serv.da.AddRolePermissions(req.RoleId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add permissions: %v", err)
	}

	return &pb.AddRolePermissionResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv indigoServiceServer) RemoveRolePermission(_ context.Context, req *pb.RemoveRolePermissionRequest) (*pb.RemoveRolePermissionResponse, error) {
	removedPerms, err := serv.da.RemoveRolePermissions(req.RoleId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove permissions: %v", err)
	}

	return &pb.RemoveRolePermissionResponse{
		RemovedPermissions: removedPerms,
	}, nil
}

func (serv indigoServiceServer) GetUserRoles(_ context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roles, err := serv.da.GetUserRoles(req.UserAccountId)
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

func (serv indigoServiceServer) AddUserRole(_ context.Context, req *pb.AddUserRoleRequest) (*pb.AddUserRoleResponse, error) {
	addedRoles, err := serv.da.AddUserRoles(req.UserAccountId, req.RoleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.AddUserRoleResponse{
		AddedRoleIds: addedRoles,
	}, nil
}

func (serv indigoServiceServer) RemoveUserRole(_ context.Context, req *pb.RemoveUserRoleRequest) (*pb.RemoveUserRoleResponse, error) {
	removedRoles, err := serv.da.RemoveUserRoles(req.UserAccountId, req.RoleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user roles: %v", err)
	}

	return &pb.RemoveUserRoleResponse{
		RemovedRoleIds: removedRoles,
	}, nil
}

func (serv indigoServiceServer) AddUserPermission(_ context.Context, req *pb.AddUserPermissionRequest) (*pb.AddUserPermissionResponse, error) {
	addedPerms, err := serv.da.AddUserPermissions(req.UserAccountId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add user permissions: %v", err)
	}

	return &pb.AddUserPermissionResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv indigoServiceServer) RemoveUserPermission(_ context.Context, req *pb.RemoveUserPermissionRequest) (*pb.RemoveUserPermissionResponse, error) {
	removedPerms, err := serv.da.RemoveUserPermissions(req.UserAccountId, req.Permissions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not remove user permissions: %v", err)
	}

	return &pb.RemoveUserPermissionResponse{
		RemovedPermissions: removedPerms,
	}, nil
}
