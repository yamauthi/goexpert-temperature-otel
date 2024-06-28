package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/application"
	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/configs"
	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/infra/web/handler"
	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/infra/web/service"
	"github.com/yamauthi/goexpert-temperature-otel/internal/pkg/infra/observability"
	"github.com/yamauthi/goexpert-temperature-otel/internal/pkg/infra/web"
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
	shutdown, err := observability.InitOtelProvider(context.Background(), conf.ServiceName, conf.OtelExporterEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer(conf.ServiceName)
	temperatureService := service.NewTemperatureService(conf.TemperatureServiceUrl)
	postalCodeInputService := application.NewPostalCodeInput(temperatureService)

	inputHandler := handler.NewPostalCodeInputHandler(postalCodeInputService, tracer)

	webserver := web.NewWebServer(
		conf.WebServerUrl,
		conf.ServiceName,
		[]web.Route{
			{
				Path:    "/temperature",
				Method:  "post",
				Handler: inputHandler.TemperatureByPostalCode,
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
