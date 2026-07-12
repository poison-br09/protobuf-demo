package main

import (
	"context"
	"fmt"
	"time"

	"protobuf-demo/generated"

	retry "github.com/avast/retry-go/v4"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Circuit breaker — shared across all requests
var cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
	Name:        "UserService",
	MaxRequests: 1,                // allow 1 request in half-open state
	Interval:    10 * time.Second, // reset counts every 10 seconds
	Timeout:     5 * time.Second,  // wait 5 seconds before half-open
	ReadyToTrip: func(counts gobreaker.Counts) bool {
		// Trip after 3 consecutive failures
		return counts.ConsecutiveFailures >= 3
	},
	OnStateChange: func(name string, from, to gobreaker.State) {
		fmt.Printf("Circuit breaker [%s]: %s → %s\n", name, from, to)
	},
})

func getUserWithRetry(client generated.UserServiceClient, id int32) (*generated.User, error) {
	var user *generated.User

	err := retry.Do(
		func() error {
			// Wrap call in circuit breaker
			result, err := cb.Execute(func() (interface{}, error) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				resp, err := client.GetUser(ctx, &generated.GetUserRequest{Id: id})
				if err != nil {
					return nil, err
				}
				return resp, nil
			})

			if err != nil {
				// Circuit breaker is open — fail fast
				if err == gobreaker.ErrOpenState {
					fmt.Println("Circuit breaker OPEN — failing fast!")
					return retry.Unrecoverable(err) // don't retry
				}

				code := status.Code(err)
				if code == codes.Unavailable || code == codes.DeadlineExceeded {
					fmt.Printf("Transient error, retrying: %v\n", err)
					return err
				}
				return retry.Unrecoverable(err)
			}

			user = result.(*generated.User)
			return nil
		},
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
		retry.MaxDelay(2*time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("Retry attempt #%d after error: %v\n", n+1, err)
		}),
	)

	return user, err
}
