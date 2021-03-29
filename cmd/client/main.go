package main

import (
	"context"
	pb "github.com/cownetwork/indigo/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":6969", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRolesServiceClient(conn)
	role, err := client.Get(context.Background(), &pb.RolesGetRequest{
		RoleId: "minecraft_player",
	})
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	log.Println(role)
}
