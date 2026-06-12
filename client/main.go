package main

import (
	"context"
	"fmt"
	"log"
	"protobuf-demo/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	defer conn.Close()

	client := generated.NewUserServiceClient(conn)

	response, err := client.GetUser(context.Background(), &generated.GetUserRequest{})

	if err != nil {
		log.Fatal("Error calling GetUser:", err)
	}

	fmt.Printf("Id: %d\nName: %s\nEmail: %s\nCity: %s\nPhones: %v\n",
		response.Id,
		response.Name,
		// response.Email,
		response.Address.City,
		response.Phones,
	)

	fmt.Println("\n--- Activity Log Stream ---")
	activityStream, err := client.GetActivityLog(context.Background(), &generated.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatal("Error calling GetActivityLog:", err)
	}

	for {
		event, err := activityStream.Recv()
		if err != nil {
			break // stream ended
		}
		fmt.Println("Event:", event.Name)
	}
}
