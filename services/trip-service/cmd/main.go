package cmd

import (
	"context"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	fare := &domain.RideFareModel{
		ID:                primitive.NewObjectID(),
		UserID:            "123",
		PackageSlug:       "123",
		TotalPriceInCents: 100,
		ExpiresAt:         time.Now().Add(time.Hour * 24),
	}
	trip, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		log.Fatalf("Failed to create trip: %v", err)
	}
	fmt.Println(trip)
}
