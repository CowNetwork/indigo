package main

import (
	"context"
	"github.com/upper/db/v4/adapter/postgresql"
	"google.golang.org/grpc"
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
	s.RegisterService(&pb.RolesService_ServiceDesc, &rolesServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// TODO
	// dao:
	// - select roles
	// - insert roles
	// - modify role permissions
	// -> important: why does protoc-gen-go_rpc not work?
}

func getEnvOrDefault(env string, def string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return def
	}
	return value
}

type rolesServiceServer struct {
	pb.UnimplementedRolesServiceServer
}

func (rolesServiceServer) Get(context.Context, *pb.RolesGetRequest) (*pb.RolesGetResponse, error) {
	return &pb.RolesGetResponse{
		Role: &pb.Role{
			Id:        "minecraft_player",
			Priority:  1,
			Transient: false,
			Color:     "6699ff",
		},
	}, nil
}
