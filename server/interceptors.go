package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

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

	fmt.Printf("[END] RPC: %s | Duration: %v | Error: %v \n", info.FullMethod, duration, err)

	return resp, err
}
