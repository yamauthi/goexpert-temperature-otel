package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var ErrInvalidPostalCode = errors.New("invalid postal code")
var InvalidPostalCodeStatus = http.StatusUnprocessableEntity

type PostalCodeInputService struct {
	TemperatureService TemperatureServiceInterface
}

func NewPostalCodeInput(
	temperatureService TemperatureServiceInterface,
) *PostalCodeInputService {
	return &PostalCodeInputService{
		TemperatureService: temperatureService,
	}
}

type PostalCodeInputDTO struct {
	PostalCode string `json:"cep"`
}

func (t *PostalCodeInputService) Execute(ctx context.Context, input PostalCodeInputDTO) (TemperatureResponse, error) {
	if !isPostalCodeValid(input.PostalCode) {
		return TemperatureResponse{StatusCode: InvalidPostalCodeStatus}, ErrInvalidPostalCode
	}

	cityTempResp, err := t.TemperatureService.GetTemperature(ctx, input.PostalCode)
	if err != nil {
		if cityTempResp.StatusCode != http.StatusInternalServerError {
			return cityTempResp, err
		}

		return cityTempResp, fmt.Errorf("unexpected error on TemperatureService: %s", err)
	}

	return cityTempResp, nil
}

func isPostalCodeValid(postalCode string) bool {
	match, err := regexp.MatchString("([0-9]{8})", postalCode)

	if err == nil && len(postalCode) == 8 && match {
		return true
	}

	return false
}
