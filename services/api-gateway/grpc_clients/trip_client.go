package grpcclients

import (
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)

type TripServiceClient struct {
	Client pb.TripServiceClient
	Conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {
	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")

	if tripServiceUrl == "" {
		tripServiceUrl = "trip-service:9093"
	}

	conn, err := grpc.NewClient(tripServiceUrl)
	if err != nil {
		return nil, err
	}

	return &TripServiceClient{
		Client: pb.NewTripServiceClient(conn),
		Conn:   conn,
	}, nil
}

func (c *TripServiceClient) Close() error {
	return c.Conn.Close()
}
