package main

import (
	"context"
	"fmt"
	"log"
	"protobuf-demo/generated"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	client := generated.NewUserServiceClient(conn)

	// Generate token
	token, err := generateToken("abhay-agent")
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}

	// Attach token to context metadata
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx, cancel := context.WithTimeout(
		metadata.NewOutgoingContext(context.Background(), md),
		10*time.Second,
	)
	defer cancel()

	// response, err := client.GetUser(ctx, &generated.GetUserRequest{Id: 1})
	response, err := getUserWithRetry(client, 1)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			fmt.Println("Error: User not found!")
		case codes.InvalidArgument:
			fmt.Println("Error: Invalid argument!")
		case codes.DeadlineExceeded:
			fmt.Println("Error: Request timed out!")
		case codes.Unauthenticated:
			fmt.Println("Error: Unauthenticated!")
		default:
			fmt.Println("Unexpected error:", err)
		}
	} else {
		fmt.Printf("Name: %s, Status: %s\n", response.Name, response.Status)
		fmt.Printf("Id: %d, City: %s, Phones: %v\n",
			response.Id,
			response.Address.City,
			response.Phones,
		)
	}

	fmt.Println("\n--- Activity Log Stream ---")
	activityStream, err := client.GetActivityLog(ctx, &generated.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatal("Error calling GetActivityLog:", err)
	}
	for {
		event, err := activityStream.Recv()
		if err != nil {
			break
		}
		fmt.Println("Event:", event.Name)
	}

	fmt.Println("\n--- Uploading Users ---")
	uploadStream, err := client.UploadUsers(ctx)
	if err != nil {
		log.Fatal("Error starting upload:", err)
	}
	users := []string{"Abhay", "Rahul", "Priya"}
	for _, name := range users {
		uploadStream.Send(&generated.User{Name: name})
		time.Sleep(300 * time.Millisecond)
	}
	summary, err := uploadStream.CloseAndRecv()
	if err != nil {
		log.Fatal("Error closing stream:", err)
	}
	fmt.Printf("Upload complete: %s (count: %d)\n", summary.Message, summary.Count)

	fmt.Println("\n--- Bidirectional Chat ---")
	chatStream, err := client.Chat(ctx)
	if err != nil {
		log.Fatal("Error starting chat:", err)
	}
	for _, name := range []string{"Abhay", "Rahul", "Priya"} {
		chatStream.Send(&generated.User{Name: name})
		fmt.Printf("Client sent: %s\n", name)
		reply, err := chatStream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Server replied: %s\n", reply.Message)
		time.Sleep(300 * time.Millisecond)
	}
	chatStream.CloseSend()
}
