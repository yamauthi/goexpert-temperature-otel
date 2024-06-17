package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/application/repository"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/domain/entity"
	"go.opentelemetry.io/otel/trace"
)

const AddressRepositoryRequest = "address-repository-request"
const TemperatureRepositoryRequest = "weather-repository-request"

var ErrInvalidPostalCode = errors.New("invalid postal code")
var ErrPostalCodeNotFound = errors.New("can not find postal code")

type TemperatureByPostalCodeService struct {
	PostalAddressRepository repository.PostalAddressRepositoryInterface
	TemperatureRepository   repository.TemperatureRepositoryInterface
	OTelTracer              trace.Tracer
}

func NewTemperatureByPostalCodeService(
	postalRepository repository.PostalAddressRepositoryInterface,
	temperatureRepository repository.TemperatureRepositoryInterface,
	tracer trace.Tracer,
) *TemperatureByPostalCodeService {
	return &TemperatureByPostalCodeService{
		PostalAddressRepository: postalRepository,
		TemperatureRepository:   temperatureRepository,
		OTelTracer:              tracer,
	}
}

type PostalCodeInput struct {
	PostalCode string
}

type TemperatureOutput struct {
	Location       string  `json:"city"`
	TempCelsius    float64 `json:"temp_C"`
	TempFahrenheit float64 `json:"temp_F"`
	TempKelvin     float64 `json:"temp_K"`
}

func (t *TemperatureByPostalCodeService) Execute(ctx context.Context, input PostalCodeInput) (TemperatureOutput, error) {
	if !isPostalCodeValid(input.PostalCode) {
		return TemperatureOutput{}, ErrInvalidPostalCode
	}

	// Span for postal address api
	_, addressSpan := t.OTelTracer.Start(ctx, AddressRepositoryRequest)

	address, err := t.PostalAddressRepository.GetAddress(input.PostalCode)
	if err != nil {
		return TemperatureOutput{}, fmt.Errorf("unexpected error when trying to get postal code information: %s", err)
	}
	addressSpan.End()

	if (address == entity.PostalAddress{}) {
		return TemperatureOutput{}, ErrPostalCodeNotFound
	}

	tempUnits := []entity.TemperatureUnit{
		entity.Celsius,
		entity.Fahrenheit,
		entity.Kelvin,
	}

	// Span for postal address api
	_, temperatureSpan := t.OTelTracer.Start(ctx, TemperatureRepositoryRequest)

	temperatures, err := t.TemperatureRepository.GetCityTemperature(address.City, tempUnits)
	if err != nil {
		return TemperatureOutput{}, fmt.Errorf("unexpected error when trying to get temperature information: %s", err)
	}
	temperatureSpan.End()

	return TemperatureOutput{
		Location:       fmt.Sprintf("%s, %s", address.City, address.ProvinceState),
		TempCelsius:    temperatures[entity.Celsius].Degrees,
		TempFahrenheit: temperatures[entity.Fahrenheit].Degrees,
		TempKelvin:     temperatures[entity.Kelvin].Degrees,
	}, nil
}

func isPostalCodeValid(postalCode string) bool {
	match, err := regexp.MatchString("([0-9]{8})", postalCode)

	if err == nil && len(postalCode) == 8 && match {
		return true
	}

	return false
}
