package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/domain/entity"
)

type ViaCepApiRepository struct {
	Config ApiConfig
}

type viaCepInfoDTO struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

func NewViaCepApiRepository(config ApiConfig) *ViaCepApiRepository {
	return &ViaCepApiRepository{
		Config: config,
	}
}

func (api *ViaCepApiRepository) GetAddress(postalCode string) (entity.PostalAddress, error) {
	url := fmt.Sprintf("%s/%s/json", api.Config.BaseURL, postalCode)
	resp, err := http.Get(url)
	if err != nil {
		return entity.PostalAddress{}, fmt.Errorf("unknown error on ViaCepAPI.GetAddress. Url: %s | Err: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var cepInfo viaCepInfoDTO
		if err := json.NewDecoder(resp.Body).Decode(&cepInfo); err != nil {
			return entity.PostalAddress{}, err
		}

		if (cepInfo == viaCepInfoDTO{}) {
			return entity.PostalAddress{}, nil
		}

		return entity.PostalAddress{
			PostalCode:    cepInfo.Cep,
			Address:       fmt.Sprintf("%s, %s", cepInfo.Logradouro, cepInfo.Bairro),
			City:          cepInfo.Localidade,
			ProvinceState: cepInfo.Uf,
		}, nil
	}

	return entity.PostalAddress{}, fmt.Errorf("unknown error on ViaCepAPI.GetAddress. Url: %s", url)
}
