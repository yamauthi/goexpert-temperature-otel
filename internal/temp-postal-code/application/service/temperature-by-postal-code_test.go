package service_test

// import (
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// 	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/application/service"
// 	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/domain/entity"
// )

// type MockPostalAddressRepository struct {
// 	mock.Mock
// }

// func (r *MockPostalAddressRepository) GetAddress(postalCode string) (entity.PostalAddress, error) {
// 	args := r.Called(postalCode)
// 	return args.Get(0).(entity.PostalAddress), args.Error(1)
// }

// type MockTemperatureRepository struct {
// 	mock.Mock
// }

// func (r *MockTemperatureRepository) GetCityTemperature(city string, units []entity.TemperatureUnit) (map[entity.TemperatureUnit]entity.Temperature, error) {
// 	args := r.Called(city, units)
// 	return args.Get(0).(map[entity.TemperatureUnit]entity.Temperature), args.Error(1)
// }

// type TemperatureByPostalCodeServiceTestSuite struct {
// 	suite.Suite
// 	Service                     *service.TemperatureByPostalCodeService
// 	MockPostalAddressRepository *MockPostalAddressRepository
// 	MockTemperatureRepository   *MockTemperatureRepository
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) SetupTest() {
// 	suite.MockPostalAddressRepository = &MockPostalAddressRepository{}
// 	suite.MockTemperatureRepository = &MockTemperatureRepository{}

// 	suite.Service = service.NewTemperatureByPostalCodeService(
// 		suite.MockPostalAddressRepository,
// 		suite.MockTemperatureRepository,
// 	)
// }

// func TestTemperatureByPostalCodeServiceTestSuite(t *testing.T) {
// 	suite.Run(t, new(TemperatureByPostalCodeServiceTestSuite))
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) TestTemperatureByPostalCodeService_Execute_InvalidPostalCode() {

// 	testsCase := []string{
// 		"",
// 		"123456789",
// 		"12345-678",
// 		"12 345-678",
// 		"1234-567",
// 		"12345",
// 	}

// 	for _, input := range testsCase {
// 		output, err := suite.Service.Execute(service.PostalCodeInput{
// 			PostalCode: input,
// 		})

// 		suite.Equal(service.ErrInvalidPostalCode, err)
// 		suite.Empty(output)
// 		suite.MockPostalAddressRepository.AssertNumberOfCalls(suite.T(), "GetAddress", 0)
// 		suite.MockTemperatureRepository.AssertNumberOfCalls(suite.T(), "GetCityTemperature", 0)
// 	}
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) TestTemperatureByPostalCodeService_Execute_GetAddressError() {
// 	mockError := errors.New("could not connect to api")
// 	suite.MockPostalAddressRepository.On("GetAddress", mock.Anything).Return(
// 		entity.PostalAddress{},
// 		mockError,
// 	)
// 	expectedError := fmt.Errorf("unexpected error when trying to get postal code information: %s", mockError)
// 	output, err := suite.Service.Execute(service.PostalCodeInput{
// 		PostalCode: "31270901",
// 	})

// 	suite.Equal(expectedError, err)
// 	suite.Empty(output)
// 	suite.MockPostalAddressRepository.AssertNumberOfCalls(suite.T(), "GetAddress", 1)
// 	suite.MockTemperatureRepository.AssertNumberOfCalls(suite.T(), "GetCityTemperature", 0)
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) TestTemperatureByPostalCodeService_Execute_PostalCodeNotFound() {
// 	input := "12345678"
// 	suite.MockPostalAddressRepository.On("GetAddress", input).Return(
// 		entity.PostalAddress{},
// 		nil,
// 	)

// 	output, err := suite.Service.Execute(service.PostalCodeInput{
// 		PostalCode: input,
// 	})

// 	suite.Equal(service.ErrPostalCodeNotFound, err)
// 	suite.Empty(output)
// 	suite.MockPostalAddressRepository.AssertNumberOfCalls(suite.T(), "GetAddress", 1)
// 	suite.MockTemperatureRepository.AssertNumberOfCalls(suite.T(), "GetCityTemperature", 0)
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) TestTemperatureByPostalCodeService_Execute_GetCityTemperatureError() {
// 	input := "31270901"
// 	postalAddress := entity.PostalAddress{
// 		PostalCode:    "31270901",
// 		Address:       "Av. Pres. Antônio Carlos, 6627 - Pampulha MG, 31270-901, Brasil",
// 		City:          "Belo Horizonte",
// 		ProvinceState: "MG",
// 	}
// 	suite.MockPostalAddressRepository.On("GetAddress", input).Return(
// 		postalAddress,
// 		nil,
// 	)

// 	tempUnits := []entity.TemperatureUnit{
// 		entity.Celsius,
// 		entity.Fahrenheit,
// 		entity.Kelvin,
// 	}
// 	mockError := errors.New("could not connect to api")
// 	suite.MockTemperatureRepository.On("GetCityTemperature", postalAddress.City, tempUnits).Return(
// 		map[entity.TemperatureUnit]entity.Temperature{},
// 		mockError,
// 	)

// 	expectedError := fmt.Errorf("unexpected error when trying to get temperature information: %s", mockError)

// 	output, err := suite.Service.Execute(service.PostalCodeInput{
// 		PostalCode: input,
// 	})

// 	suite.Equal(expectedError, err)
// 	suite.Empty(output)
// 	suite.MockPostalAddressRepository.AssertNumberOfCalls(suite.T(), "GetAddress", 1)
// 	suite.MockTemperatureRepository.AssertNumberOfCalls(suite.T(), "GetCityTemperature", 1)
// }

// func (suite *TemperatureByPostalCodeServiceTestSuite) TestTemperatureByPostalCodeService_Execute() {
// 	input := "31270901"
// 	postalAddress := entity.PostalAddress{
// 		PostalCode:    "31270901",
// 		Address:       "Av. Pres. Antônio Carlos, 6627 - Pampulha MG, 31270-901, Brasil",
// 		City:          "Belo Horizonte",
// 		ProvinceState: "MG",
// 	}
// 	suite.MockPostalAddressRepository.On("GetAddress", input).Return(
// 		postalAddress,
// 		nil,
// 	)

// 	tempUnits := []entity.TemperatureUnit{
// 		entity.Celsius,
// 		entity.Fahrenheit,
// 		entity.Kelvin,
// 	}
// 	suite.MockTemperatureRepository.On("GetCityTemperature", postalAddress.City, tempUnits).Return(
// 		map[entity.TemperatureUnit]entity.Temperature{
// 			entity.Celsius: {
// 				Degrees: 25,
// 				Unit:    entity.Celsius,
// 			},
// 			entity.Fahrenheit: {
// 				Degrees: 77,
// 				Unit:    entity.Fahrenheit,
// 			},
// 			entity.Kelvin: {
// 				Degrees: 298,
// 				Unit:    entity.Kelvin,
// 			},
// 		},
// 		nil,
// 	)

// 	expectedOutput := service.TemperatureOutput{
// 		Location:       fmt.Sprintf("%s, %s", postalAddress.City, postalAddress.ProvinceState),
// 		TempCelsius:    25,
// 		TempFahrenheit: 77,
// 		TempKelvin:     298,
// 	}

// 	output, err := suite.Service.Execute(service.PostalCodeInput{
// 		PostalCode: input,
// 	})

// 	suite.Nil(err)
// 	suite.Equal(expectedOutput, output)
// 	suite.MockPostalAddressRepository.AssertNumberOfCalls(suite.T(), "GetAddress", 1)
// 	suite.MockTemperatureRepository.AssertNumberOfCalls(suite.T(), "GetCityTemperature", 1)
// }
