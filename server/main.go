package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"protobuf-demo/generated"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	generated.UnimplementedUserServiceServer
}

func (s *userServiceServer) GetUser(ctx context.Context, req *generated.GetUserRequest) (*generated.User, error) {
	fmt.Println("Server received request for ID:", req.Id)

	return &generated.User{
		Id:   req.Id,
		Name: "Abhay Kumar",
		// Email: "abhay@zenwork.com",
		Address: &generated.Address{
			Street: "123 MG Road",
			City:   "Hyderabad",
		},
		Phones: []string{"9999999999", "8888888888"},
	}, nil
}

func (s *userServiceServer) GetActivityLog(req *generated.GetUserRequest, stream grpc.ServerStreamingServer[generated.User]) error {
	fmt.Println("Server streaming activity log for ID:", req.Id)

	//simulates streaming 5 activity events

	activities := []string{
		"User logged in",
		"User viewed dashboard",
		"User updated profile",
		"User uploaded file",
		"User logged out",
	}

	for _, activity := range activities {
		err := stream.Send(&generated.User{
			Id:   req.Id,
			Name: activity,
		})
		if err != nil {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (s *userServiceServer) UploadUsers(stream grpc.ClientStreamingServer[generated.User, generated.UploadSummary]) error {
	count := 0

	for {
		user, err := stream.Recv()
		if err != nil {
			break
		}

		fmt.Println("Server recieved user: %s\n", user.Name)
		count++
	}

	return stream.SendAndClose(&generated.UploadSummary{
		Count:   int32(count),
		Message: fmt.Sprintf("Succefully received %d users", count),
	})
}

func (s *userServiceServer) Chat(stream grpc.BidiStreamingServer[generated.User, generated.ServerReply]) error {
	for {
		user, err := stream.Recv()
		if err != nil {
			break
		}

		fmt.Printf("Server received message from: %s\n", user.Name)

		stream.Send(&generated.ServerReply{
			Message: fmt.Sprintf("Hello %s! Welcome to the server.", user.Name),
		})
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	generated.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	fmt.Println("gRPC server running on port 50051")
	grpcServer.Serve(lis)
}
