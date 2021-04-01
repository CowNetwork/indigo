package main

import (
	"context"
	"github.com/upper/db/v4"
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
	s.RegisterService(&pb.IndigoService_ServiceDesc, &indigoServiceServer{sess: sess})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// TODO
	// - bind role to user
	// - unbind roles from user
	// - bind permission to user
	// - unbind permission from user
}

func getEnvOrDefault(env string, def string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return def
	}
	return value
}

type Role struct {
	Id        string `db:"id"`
	Priority  int32  `db:"priority"`
	Transient bool   `db:"transient"`
	Color     string `db:"color"`
}

func FromProtoRole(r *pb.Role) *Role {
	return &Role{
		Id:        r.Id,
		Priority:  r.Priority,
		Transient: r.Transient,
		Color:     r.Color,
	}
}

func (r Role) ToProtoRole() *pb.Role {
	return &pb.Role{
		Id:        r.Id,
		Priority:  r.Priority,
		Transient: r.Transient,
		Color:     r.Color,
	}
}

type RolePermissionBinding struct {
	RoleId     string `db:"role_id"`
	Permission string `db:"permission"`
}

type UserRoleBinding struct {
	UserAccountId string `db:"user_account_id"`
	RoleId        string `db:"role_id"`
}

type UserPermissionBinding struct {
	UserAccountId string `db:"user_account_id"`
	Permission    string `db:"permission"`
}

type indigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	sess db.Session
}

func (serv indigoServiceServer) GetRole(_ context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	coll := serv.sess.Collection("role_definitions")
	res := coll.Find("id", req.RoleId)

	var role Role
	err := res.One(&role)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "could not fine role with id %s", req.RoleId)
	}

	res = serv.sess.Collection("role_permissions").Find("role_id", req.RoleId)

	var bindings []RolePermissionBinding

	// ignore error, the role just does not have any permissions
	_ = res.All(&bindings)

	permissions := make([]string, len(bindings))
	for i, binding := range bindings {
		permissions[i] = binding.Permission
	}

	r := role.ToProtoRole()
	r.Permissions = permissions

	return r, nil
}

func (serv indigoServiceServer) InsertRole(_ context.Context, role *pb.Role) (*pb.InsertRoleResponse, error) {
	coll := serv.sess.Collection("role_definitions")

	res, err := coll.Insert(FromProtoRole(role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	return &pb.InsertRoleResponse{
		Successful: res.ID() != nil,
	}, nil
}

func (serv indigoServiceServer) AddRolePermission(_ context.Context, req *pb.AddRolePermissionRequest) (*pb.AddRolePermissionResponse, error) {
	coll := serv.sess.Collection("role_permissions")

	addedPerms := make([]string, len(req.Permissions))
	for _, perm := range req.Permissions {
		binding := RolePermissionBinding{
			RoleId:     req.RoleId,
			Permission: perm,
		}

		exists, _ := coll.Find("role_id", req.RoleId).And("permission", perm).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(&binding)
		if err == nil && res.ID() != nil {
			addedPerms = append(addedPerms, perm)
		}
	}

	return &pb.AddRolePermissionResponse{
		AddedPermissions: addedPerms,
	}, nil
}

func (serv indigoServiceServer) RemoveRolePermission(_ context.Context, req *pb.RemoveRolePermissionRequest) (*pb.RemoveRolePermissionResponse, error) {
	coll := serv.sess.Collection("role_permissions")

	filter := make([]db.LogicalExpr, len(req.Permissions))
	for i, permission := range req.Permissions {
		filter[i] = db.Cond{"permission": permission}
	}

	err := coll.Find("role_id", req.RoleId).And(db.Or(filter...)).Delete()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete permission binding: %v", err)
	}

	return &pb.RemoveRolePermissionResponse{
		RemovedPermissions: req.Permissions,
	}, nil
}

func (serv indigoServiceServer) GetUserRoles(_ context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	coll := serv.sess.Collection("user_roles")
	res := coll.Find("user_account_id", req.UserAccountId)

	var roles []Role
	err := res.All(roles)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get roles of user account: %v", err)
	}

	var protoRoles []*pb.Role
	for _, role := range roles {
		res = serv.sess.Collection("role_permissions").Find("role_id", role.Id)

		var bindings []RolePermissionBinding

		// ignore error, the role just does not have any permissions
		_ = res.All(&bindings)

		permissions := make([]string, len(bindings))
		for i, binding := range bindings {
			permissions[i] = binding.Permission
		}

		r := role.ToProtoRole()
		r.Permissions = permissions

		protoRoles = append(protoRoles, r)
	}

	return &pb.GetUserRolesResponse{
		Roles: protoRoles,
	}, nil
}

func (serv indigoServiceServer) AddUserRole(_ context.Context, req *pb.AddUserRoleRequest) (*pb.AddUserRoleResponse, error) {
	coll := serv.sess.Collection("user_roles")

	addedRoles := make([]string, len(req.RoleIds))
	for _, id := range req.RoleIds {
		binding := UserRoleBinding{
			UserAccountId: req.UserAccountId,
			RoleId:        id,
		}

		exists, _ := coll.Find("user_account_id", req.UserAccountId).And("role_id", id).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(binding)
		if err == nil && res.ID() != nil {
			addedRoles = append(addedRoles, id)
		}
	}

	return &pb.AddUserRoleResponse{
		AddedRoleIds: nil,
	}, nil
}

func (serv indigoServiceServer) RemoveUserRole(_ context.Context, req *pb.RemoveUserRoleRequest) (*pb.RemoveUserRoleResponse, error) {
	coll := serv.sess.Collection("user_roles")

	filter := make([]db.LogicalExpr, len(req.RoleIds))
	for i, permission := range req.RoleIds {
		filter[i] = db.Cond{"role_id": permission}
	}

	err := coll.Find("user_account_id", req.UserAccountId).And(db.Or(filter...)).Delete()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete role binding: %v", err)
	}

	return &pb.RemoveUserRoleResponse{
		RemovedRoleIds: req.RoleIds,
	}, nil
}

func (serv indigoServiceServer) AddUserPermission(_ context.Context, req *pb.AddUserPermissionRequest) (*pb.AddUserPermissionResponse, error) {
	coll := serv.sess.Collection("user_permissions")

	addedPermissions := make([]string, len(req.Permissions))
	for _, permission := range req.Permissions {
		binding := UserPermissionBinding{
			UserAccountId: req.UserAccountId,
			Permission:    permission,
		}

		exists, _ := coll.Find("user_account_id", req.UserAccountId).And("permission", permission).Exists()
		if exists {
			continue
		}

		res, err := coll.Insert(binding)
		if err == nil && res.ID() != nil {
			addedPermissions = append(addedPermissions, permission)
		}
	}

	return &pb.AddUserPermissionResponse{
		AddedPermissions: addedPermissions,
	}, nil
}

func (serv indigoServiceServer) RemoveUserPermission(_ context.Context, req *pb.RemoveUserPermissionRequest) (*pb.RemoveUserPermissionResponse, error) {
	coll := serv.sess.Collection("user_permissions")

	filter := make([]db.LogicalExpr, len(req.Permissions))
	for i, permission := range req.Permissions {
		filter[i] = db.Cond{"permission": permission}
	}

	err := coll.Find("user_account_id", req.UserAccountId).And(db.Or(filter...)).Delete()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete permission binding: %v", err)
	}

	return &pb.RemoveUserPermissionResponse{
		RemovedPermissions: req.Permissions,
	}, nil
}
