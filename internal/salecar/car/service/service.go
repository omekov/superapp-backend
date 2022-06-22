package service

import carrepository "github.com/omekov/superapp-backend/internal/salecar/car/repository"

// Service ...
type Service struct {
	carrepository *carrepository.CarRepository
}

// NewService ...
func NewService(carrepository *carrepository.CarRepository) *Service {
	return &Service{
		carrepository: carrepository,
	}
}
