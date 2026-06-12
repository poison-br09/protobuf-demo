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
