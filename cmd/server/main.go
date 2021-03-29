package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/cownetwork/indigo/proto"
)

func main() {
	log.Println("Hello World!")

	lis, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	s.RegisterService(&pb.RolesService_ServiceDesc, &rolesServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
