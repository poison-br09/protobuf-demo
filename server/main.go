package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"protobuf-demo/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	generated.UnimplementedUserServiceServer
}

func (s *userServiceServer) GetUser(ctx context.Context, req *generated.GetUserRequest) (*generated.User, error) {
	fmt.Println("Server received request for ID:", req.Id)

	// Validate input
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	// Simulate user not found
	if req.Id == 99 {
		return nil, status.Error(codes.NotFound, "user with given id not found")
	}

	// Simulate slow processing
	fmt.Println("Processing...")

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Processing complete!")
	case <-ctx.Done():
		fmt.Println("Client cancelled! Stopping work.")
		return nil, status.Error(codes.Canceled, "request cancelled by client")
	}

	return &generated.User{
		Id:     req.Id,
		Name:   "Abhay Kumar",
		Age:    25,
		Status: generated.UserStatus_USER_STATUS_ACTIVE,
		Address: &generated.Address{
			Street: "123 MG Road",
			City:   "Hyderabad",
		},
		Phones: []string{"9999999999", "8888888888"},
		Contact: &generated.User_EmailContact{
			EmailContact: "abhay@zenwork.com",
		},
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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			authInterceptor,
		),
	)
	generated.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	fmt.Println("gRPC server running on port 50051")
	grpcServer.Serve(lis)
}
