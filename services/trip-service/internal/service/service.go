package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	pkgtypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, trip)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*pkgtypes.OSRMApiResponse, error) {
	url := fmt.Sprintf("'http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson'", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch route from OSRM API: %v", err)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read the response: %v", err)
	}

	defer res.Body.Close()

	var routeResp pkgtypes.OSRMApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %v", err)
	}

	return &routeResp, nil
}

func (s *service) EstimatePackagePriceWithRoute(route *pkgtypes.OSRMApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoutes(f, route)
	}

	return estimatedFares
}

func (s *service) GenerateTripFares(context context.Context, rideFares []*domain.RideFareModel, userid string) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, f := range rideFares {
		id := primitive.NewObjectID()

		fare := &domain.RideFareModel{
			UserID:            userid,
			ID:                id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}

		if err := s.repo.SaveRideFare(context, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare %w", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func estimateFareRoutes(f *domain.RideFareModel, route *pkgtypes.OSRMApiResponse) *domain.RideFareModel {
	pricingCfg := pkgtypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Route[0].Distance
	durationInMinutes := route.Route[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute

	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
