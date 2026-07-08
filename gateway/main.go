package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"protobuf-demo/generated"
)

func main() {
	// Connect to your existing gRPC server
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", err)
	}
	defer conn.Close()

	// Create gateway mux
	mux := runtime.NewServeMux()

	// Register your gRPC service with the gateway
	err = generated.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatal("Failed to register gateway:", err)
	}

	fmt.Println("REST Gateway running on port 8080...")
	fmt.Println("Try: curl http://localhost:8080/v1/users/1")

	// Start HTTP server
	http.ListenAndServe(":8080", mux)
}
