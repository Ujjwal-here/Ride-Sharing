package domain

import (
	"context"
	"ride-sharing/shared/types"

	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"

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
	SaveRideFare(context context.Context, f *RideFareModel) error
	GetRideFareByID(context context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(context context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(context context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(context context.Context, fares []*RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(context context.Context, fareID, userID string) (*RideFareModel, error)
}
