package main

import (
	"context"
	pb "github.com/cownetwork/indigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	conn, err := grpc.Dial(":6969", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRolesServiceClient(conn)
	role, err := client.Get(context.Background(), &pb.StringValue{
		Value: "minecraft_player",
	})

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalf("failed to get role: %v", err)
	}

	if st.Code() != codes.OK {
		// insert
		_, err := client.Insert(context.Background(), &pb.Role{
			Id:        "minecraft_player",
			Priority:  1,
			Transient: false,
			Color:     "6699ff",
		})

		if err != nil {
			log.Fatalf("failed to insert role: %v", err)
		}
		log.Println("Inserted role.")
	} else {
		log.Println(role)
	}

	_, err = client.AddPermission(context.Background(), &pb.RolesAddPermissionRequest{
		RoleId: "minecraft_player",
		Permissions: []string{
			"bukkit.command.help",
			"bukkit.command.gamemode",
		},
	})
	if err != nil {
		log.Fatalf("failed to add permission: %v", err)
	}
	log.Println("Added permission.")

}
