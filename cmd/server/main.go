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
	s.RegisterService(&pb.RolesService_ServiceDesc, &rolesServiceServer{sess: sess})

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

type rolesServiceServer struct {
	pb.UnimplementedRolesServiceServer
	sess db.Session
}

func (serv rolesServiceServer) Get(_ context.Context, roleId *pb.StringValue) (*pb.Role, error) {
	coll := serv.sess.Collection("role_definitions")
	res := coll.Find("id", roleId.Value)

	var role Role
	err := res.One(&role)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "could not fine role with id %s", roleId.Value)
	}

	res = serv.sess.Collection("role_permissions").Find("role_id", roleId.Value)

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

func (serv rolesServiceServer) Insert(_ context.Context, role *pb.Role) (*pb.BoolValue, error) {
	coll := serv.sess.Collection("role_definitions")

	res, err := coll.Insert(FromProtoRole(role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert role: %v", err)
	}

	return &pb.BoolValue{
		Value: res.ID() != nil,
	}, nil
}

func (serv rolesServiceServer) AddPermission(_ context.Context, perm *pb.RolePermissionBinding) (*pb.BoolValue, error) {
	coll := serv.sess.Collection("role_permissions")

	binding := RolePermissionBinding{
		RoleId:     perm.RoleId,
		Permission: perm.Permission,
	}
	res, err := coll.Insert(&binding)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert permission binding: %v", err)
	}

	return &pb.BoolValue{
		Value: res.ID() != nil,
	}, nil
}

func (serv rolesServiceServer) RemovePermission(_ context.Context, perm *pb.RolePermissionBinding) (*pb.BoolValue, error) {
	coll := serv.sess.Collection("role_permissions")

	err := coll.Find("role_id", perm.RoleId).Delete()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete permission binding: %v", err)
	}

	return &pb.BoolValue{
		Value: true,
	}, nil
}
