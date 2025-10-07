package grpcserver

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGrpcHandler(server *grpc.Server, service domain.TripService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *grpcHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetPickup()
	destination := req.GetDestination()

	pickupCoord := &types.Coordinate{
		Latitude:  pickup.GetLatitude(),
		Longitude: pickup.GetLongitude(),
	}
	destinationCoord := &types.Coordinate{
		Latitude:  destination.GetLatitude(),
		Longitude: destination.GetLongitude(),
	}
	route, err := h.service.GetRoute(ctx, pickupCoord, destinationCoord)

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "Failed to get route %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToRoute(),
		RideFares: []*pb.RideFare{},
	}, nil
}
