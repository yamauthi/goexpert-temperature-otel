package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/domain/entity"
)

type WeatherAPIRepository struct {
	Config ApiConfig
}

type weatherApiResponseDTO struct {
	Temperature weatherApiTemperature `json:"current"`
}

type weatherApiTemperature struct {
	TempCelsius    float64 `json:"temp_c"`
	TempFahrenheit float64 `json:"temp_f"`
}

func NewWeatherAPIRepository(config ApiConfig) *WeatherAPIRepository {
	return &WeatherAPIRepository{
		Config: config,
	}
}

func (api *WeatherAPIRepository) GetCityTemperature(
	city string,
	units []entity.TemperatureUnit,
) (map[entity.TemperatureUnit]entity.Temperature, error) {
	escapedCity := url.QueryEscape(city)

	url := fmt.Sprintf("%s?key=%s&q=%s", api.Config.BaseURL, api.Config.ApiKey, escapedCity)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var weatherResponse weatherApiResponseDTO
		if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
			return nil, err
		}

		celsiusTemp := entity.Temperature{
			Degrees: weatherResponse.Temperature.TempCelsius,
			Unit:    entity.Celsius,
		}

		result := make(map[entity.TemperatureUnit]entity.Temperature)
		for _, unit := range units {
			temp, err := entity.ConvertTemperatureTo(celsiusTemp, unit)
			if err != nil {
				return nil, err
			}

			result[unit] = temp
		}

		return result, nil
	}

	return nil, fmt.Errorf("unknown error on WeatherAPIRepository.GetCityTemperature. Url: %s", url)
}
