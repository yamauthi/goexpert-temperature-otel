package entity

import "errors"

var ErrUnitOrConversionUnknown = errors.New("unit or conversion method unknown")

type TemperatureUnit int

const (
	UnknownUnit TemperatureUnit = iota
	Celsius
	Fahrenheit
	Kelvin
)

func (tu TemperatureUnit) String() string {
	switch tu {
	case Celsius:
		return "Celsius"
	case Fahrenheit:
		return "Fahrenheit"
	case Kelvin:
		return "Kelvin"
	default:
		return "Unknown"
	}
}

type Temperature struct {
	Degrees float64
	Unit    TemperatureUnit
}

func ConvertTemperatureTo(from Temperature, to TemperatureUnit) (Temperature, error) {
	switch from.Unit {
	case Celsius:
		switch to {
		case Celsius:
			return from, nil
		case Fahrenheit:
			return Temperature{
				Degrees: celsiusToFahrenheit(from.Degrees),
				Unit:    Fahrenheit,
			}, nil
		case Kelvin:
			return Temperature{
				Degrees: celsiusToKelvin(from.Degrees),
				Unit:    Kelvin,
			}, nil
		}
	case Fahrenheit:
		switch to {
		case Celsius:
			return Temperature{
				Degrees: fahrenheitToCelsius(from.Degrees),
				Unit:    Celsius,
			}, nil
		case Fahrenheit:
			return from, nil
		case Kelvin:
			return Temperature{
				Degrees: fahrenheitToKelvin(from.Degrees),
				Unit:    Kelvin,
			}, nil
		}
	case Kelvin:
		switch to {
		case Celsius:
			return Temperature{
				Degrees: kelvinToCelsius(from.Degrees),
				Unit:    Celsius,
			}, nil
		case Fahrenheit:
			return Temperature{
				Degrees: kelvinToFahrenheit(from.Degrees),
				Unit:    Fahrenheit,
			}, nil
		case Kelvin:
			return from, nil
		}
	}

	return Temperature{}, ErrUnitOrConversionUnknown
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func fahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

func fahrenheitToKelvin(fahrenheit float64) float64 {
	return fahrenheitToCelsius(fahrenheit) + 273
}

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273
}

func kelvinToFahrenheit(kelvin float64) float64 {
	return celsiusToFahrenheit(
		kelvinToCelsius(kelvin),
	)
}
