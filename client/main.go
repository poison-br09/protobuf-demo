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
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	defer conn.Close()

	client := generated.NewUserServiceClient(conn)

	response, err := client.GetUser(context.Background(), &generated.GetUserRequest{Id: 99})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			fmt.Println("Error: User not found!")
		case codes.InvalidArgument:
			fmt.Println("Error: Invalid argument!")
		default:
			fmt.Println("Unexpected error:", err)
		}
	} else {
		fmt.Printf("Name: %s, Status: %s\n", response.Name, response.Status)
		fmt.Printf("Id: %d\nName: %s\nCity: %s\nPhones: %v\n",
			response.Id,
			response.Name,
			response.Address.City,
			response.Phones,
		)
	}

	// fmt.Printf("Id: %d\nName: %s\nEmail: %s\nCity: %s\nPhones: %v\n",
	// 	response.Id,
	// 	response.Name,
	// 	// response.Email,
	// 	response.Address.City,
	// 	response.Phones,
	// )

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

	fmt.Println("\n--- Uploading Users ---")
	uploadStream, err := client.UploadUsers(context.Background())
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

	chatStream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatal("Error starting chat:", err)
	}

	names := []string{"Abhay", "Rahul", "Priya"}

	for _, name := range names {
		// Send message
		chatStream.Send(&generated.User{Name: name})
		fmt.Printf("Client sent: %s\n", name)

		// Immediately wait for reply
		reply, err := chatStream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Server replied: %s\n", reply.Message)

		time.Sleep(300 * time.Millisecond)
	}

	chatStream.CloseSend()

	// fmt.Printf("Name: %s, Status: %s\n",
	// 	response.Name,
	// 	response.Status,
	// )

	// switch c := response.Contact.(type) {
	// case *generated.User_EmailContact:
	// 	fmt.Println("Contact via Email:", c.EmailContact)
	// case *generated.User_PhoneContact:
	// 	fmt.Println("Contact via Phone:", c.PhoneContact)
	// default:
	// 	fmt.Println("No contact info")
	// }
}
