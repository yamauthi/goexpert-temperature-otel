package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/application"
	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TemperatureService struct {
	ServiceUrl string
}

func NewTemperatureService(serviceUrl string) *TemperatureService {
	return &TemperatureService{
		ServiceUrl: serviceUrl,
	}
}

func (ts *TemperatureService) GetTemperature(ctx context.Context, postalCode string) (application.TemperatureResponse, error) {
	url := fmt.Sprintf("%s/%s", ts.ServiceUrl, postalCode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return application.TemperatureResponse{}, err
	}

	// inject otel header carrier
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return application.TemperatureResponse{}, fmt.Errorf("unknown error on TemperatureService.GetTemperature. Url: %s | Err: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errString, err := io.ReadAll(resp.Body)
		if err != nil {
			return application.TemperatureResponse{}, fmt.Errorf("unknown error on TemperatureService.GetTemperature. Url: %s | Err: %s", url, err)
		}

		return application.TemperatureResponse{StatusCode: resp.StatusCode}, errors.New(string(errString))
	}

	var cityTemp domain.CityTemperature
	if err := json.NewDecoder(resp.Body).Decode(&cityTemp); err != nil {
		return application.TemperatureResponse{}, fmt.Errorf("unknown error on TemperatureService.GetTemperature. Url: %s | Err: %s", url, err)
	}

	return application.TemperatureResponse{
		StatusCode:      http.StatusOK,
		CityTemperature: cityTemp,
	}, nil
}
