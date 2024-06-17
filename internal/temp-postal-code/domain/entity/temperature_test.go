package entity_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/domain/entity"
)

type TemperatureTestSuite struct {
	suite.Suite
}

func TestTemperatureTestSuite(t *testing.T) {
	suite.Run(t, new(TemperatureTestSuite))
}

func (suite *TemperatureTestSuite) TestConvertTemperatureTo_FromUnknowUnitOrConversionMethod() {
	units := []entity.TemperatureUnit{
		entity.Celsius,
		entity.Fahrenheit,
		entity.Kelvin,
	}

	testCases := []entity.Temperature{
		{
			Degrees: 25,
			Unit:    entity.UnknownUnit,
		},
		{
			Degrees: 100,
			Unit:    -1, // Unknown Unit
		},
		{
			Degrees: 10,
			Unit:    6516465, // Unknown Unit
		},
	}

	suite.Run("convert temperature from uknown unit to existing unit", func() {
		for _, input := range testCases {
			for _, unit := range units {
				result, err := entity.ConvertTemperatureTo(input, unit)
				suite.Equal(entity.ErrUnitOrConversionUnknown, err)
				suite.Equal(entity.Temperature{}, result)
			}
		}
	})

	suite.Run("convert temperature from existing unit to unknown unit", func() {
		for _, input := range testCases {
			for _, unit := range units {
				temp := entity.Temperature{
					Degrees: input.Degrees,
					Unit:    unit,
				}
				// input.Unit = Unknown unit
				result, err := entity.ConvertTemperatureTo(temp, input.Unit)
				suite.Equal(entity.ErrUnitOrConversionUnknown, err)
				suite.Equal(entity.Temperature{}, result)
			}
		}
	})
}

func (suite *TemperatureTestSuite) TestConvertTemperatureTo() {
	type TemperatureConversion struct {
		Celsius    float64
		Fahrenheit float64
		Kelvin     float64
	}

	type TestCase struct {
		Input    entity.Temperature
		Expected TemperatureConversion
	}

	testCases := []TestCase{
		{
			Input: entity.Temperature{
				Degrees: 100,
				Unit:    entity.Celsius,
			},
			Expected: TemperatureConversion{
				Celsius:    100,
				Fahrenheit: 212,
				Kelvin:     373,
			},
		},
		{
			Input: entity.Temperature{
				Degrees: 25,
				Unit:    entity.Celsius,
			},
			Expected: TemperatureConversion{
				Celsius:    25,
				Fahrenheit: 77,
				Kelvin:     298,
			},
		},
		{
			Input: entity.Temperature{
				Degrees: -5,
				Unit:    entity.Celsius,
			},
			Expected: TemperatureConversion{
				Celsius:    -5,
				Fahrenheit: 23,
				Kelvin:     268,
			},
		},
		{
			Input: entity.Temperature{
				Degrees: -20,
				Unit:    entity.Celsius,
			},
			Expected: TemperatureConversion{
				Celsius:    -20,
				Fahrenheit: -4,
				Kelvin:     253,
			},
		},
	}

	for _, tc := range testCases {
		// from Celsius
		cToC, err := entity.ConvertTemperatureTo(tc.Input, entity.Celsius)
		suite.Nil(err)
		suite.Equal(tc.Expected.Celsius, cToC.Degrees)
		suite.Equal(entity.Celsius, cToC.Unit)

		cToF, err := entity.ConvertTemperatureTo(tc.Input, entity.Fahrenheit)
		suite.Nil(err)
		suite.Equal(tc.Expected.Fahrenheit, cToF.Degrees)
		suite.Equal(entity.Fahrenheit, cToF.Unit)

		cToK, err := entity.ConvertTemperatureTo(tc.Input, entity.Kelvin)
		suite.Nil(err)
		suite.Equal(tc.Expected.Kelvin, cToK.Degrees)
		suite.Equal(entity.Kelvin, cToK.Unit)

		// from Fahrenheit
		fahrenheit := cToF
		fToC, err := entity.ConvertTemperatureTo(fahrenheit, entity.Celsius)
		suite.Nil(err)
		suite.Equal(tc.Expected.Celsius, fToC.Degrees)
		suite.Equal(entity.Celsius, fToC.Unit)

		fToF, err := entity.ConvertTemperatureTo(cToF, entity.Fahrenheit)
		suite.Nil(err)
		suite.Equal(tc.Expected.Fahrenheit, fToF.Degrees)
		suite.Equal(entity.Fahrenheit, fToF.Unit)

		fToK, err := entity.ConvertTemperatureTo(cToF, entity.Kelvin)
		suite.Nil(err)
		suite.Equal(tc.Expected.Kelvin, fToK.Degrees)
		suite.Equal(entity.Kelvin, fToK.Unit)

		// from Kelvin
		kelvin := cToK
		kToC, err := entity.ConvertTemperatureTo(kelvin, entity.Celsius)
		suite.Nil(err)
		suite.Equal(tc.Expected.Celsius, kToC.Degrees)
		suite.Equal(entity.Celsius, kToC.Unit)

		kToF, err := entity.ConvertTemperatureTo(kelvin, entity.Fahrenheit)
		suite.Nil(err)
		suite.Equal(tc.Expected.Fahrenheit, kToF.Degrees)
		suite.Equal(entity.Fahrenheit, kToF.Unit)

		kToK, err := entity.ConvertTemperatureTo(kelvin, entity.Kelvin)
		suite.Nil(err)
		suite.Equal(tc.Expected.Kelvin, kToK.Degrees)
		suite.Equal(entity.Kelvin, kToK.Unit)
	}
}

func (suite *TemperatureTestSuite) TestTemperatureUnit_String_UnknownUnit() {
	type TestCase struct {
		Input    entity.TemperatureUnit
		Expected string
	}

	testCases := []TestCase{
		{
			Input:    entity.UnknownUnit,
			Expected: "Unknown",
		},
		{
			Input:    -1,
			Expected: "Unknown",
		},
		{
			Input:    49549219,
			Expected: "Unknown",
		},
	}

	for _, tc := range testCases {
		suite.Equal(tc.Expected, tc.Input.String())
	}
}

func (suite *TemperatureTestSuite) TestTemperatureUnit_String() {
	type TestCase struct {
		Input    entity.TemperatureUnit
		Expected string
	}

	testCases := []TestCase{
		{
			Input:    entity.Celsius,
			Expected: "Celsius",
		},
		{
			Input:    entity.Fahrenheit,
			Expected: "Fahrenheit",
		},
		{
			Input:    entity.Kelvin,
			Expected: "Kelvin",
		},
	}

	for _, tc := range testCases {
		suite.Equal(tc.Expected, tc.Input.String())
	}
}
