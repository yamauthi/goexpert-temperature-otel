package application

import (
	"context"

	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/domain"
)

type TemperatureResponse struct {
	StatusCode      int
	CityTemperature domain.CityTemperature
}

type TemperatureServiceInterface interface {
	GetTemperature(ctx context.Context, postalCode string) (TemperatureResponse, error)
}
