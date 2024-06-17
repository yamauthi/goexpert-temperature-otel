package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/yamauthi/goexpert-temperature-otel/internal/pkg/infra/observability"
	"github.com/yamauthi/goexpert-temperature-otel/internal/pkg/infra/web"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/application/service"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/configs"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/infra/api"
	handler "github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/infra/web"
	"go.opentelemetry.io/otel"
)

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Graceful shutdown setup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// otel provider
	shutdown, err := observability.InitOtelProvider(context.Background(), conf.ServiceName, conf.OtelEporterEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer(conf.ServiceName)
	postalCodeRepository := api.NewViaCepApiRepository(
		api.ApiConfig{
			BaseURL: conf.ViaCepApiUrl,
		},
	)

	temperatureRepository := api.NewWeatherAPIRepository(
		api.ApiConfig{
			BaseURL: conf.WeatherApiUrl,
			ApiKey:  conf.WeatherApiKey,
		},
	)

	temperatureService := service.NewTemperatureByPostalCodeService(
		postalCodeRepository,
		temperatureRepository,
		tracer,
	)

	temperatureHandler := handler.NewTemperatureByPostalCodeHandler(temperatureService, tracer)

	webserver := web.NewWebServer(
		conf.WebServerUrl,
		conf.ServiceName,
		[]web.Route{
			{
				Path:    "/temperature/{postalCode}",
				Method:  "get",
				Handler: temperatureHandler.TemperatureByPostalCode,
			},
		},
	)

	go webserver.Start()

	// Graceful shutdown handler
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
