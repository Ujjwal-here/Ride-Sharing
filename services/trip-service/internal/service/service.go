package service

import "ride-sharing/services/trip-service/internal/domain"

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{repo: repo}
}
