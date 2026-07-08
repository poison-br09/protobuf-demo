package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Logging interceptor
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	fmt.Printf("[START] RPC: %s\n", info.FullMethod)

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	fmt.Printf("[END] RPC: %s | Duration: %v | Error: %v\n",
		info.FullMethod, duration, err)

	return resp, err
}

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// Skip auth for health checks or public endpoints

	if info.FullMethod == "/UserService/GetUser" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

	agentID, err := validateToken(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	fmt.Printf("[AUTH] Request from agent: %s\n", agentID)

	return handler(ctx, req)
}
