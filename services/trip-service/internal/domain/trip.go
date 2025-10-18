package domain

import (
	"context"
	pkgtypes "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

type TripRepository interface {
	CreateTrip(context context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(context.Context, *RideFareModel) error
	GetRideFareByID(context context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(context context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(context context.Context, pickup, destination *types.Coordinate) (*pkgtypes.OSRMApiResponse, error)
	EstimatePackagePriceWithRoute(*pkgtypes.OSRMApiResponse) []*RideFareModel
	GenerateTripFares(context.Context, []*RideFareModel, string, *pkgtypes.OSRMApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(context context.Context, fareID, userID string) (*RideFareModel, error)
}
