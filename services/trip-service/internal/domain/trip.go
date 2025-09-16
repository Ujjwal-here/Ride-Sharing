package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
	ExpiresAt         time.Time
}

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare RideFareModel
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip TripModel) (*TripModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare RideFareModel) (*TripModel, error)
}
