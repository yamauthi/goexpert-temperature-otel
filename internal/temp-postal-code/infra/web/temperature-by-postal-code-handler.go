package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/yamauthi/goexpert-temperature-otel/internal/temp-postal-code/application/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const TemperatureByPostalCodeRequest = "temperature-by-postal-code-request"

type TemperatureByPostalCodeHandler struct {
	Service    *service.TemperatureByPostalCodeService
	OTelTracer trace.Tracer
}

func NewTemperatureByPostalCodeHandler(service *service.TemperatureByPostalCodeService, tracer trace.Tracer) *TemperatureByPostalCodeHandler {
	return &TemperatureByPostalCodeHandler{
		Service:    service,
		OTelTracer: tracer,
	}
}

func (h *TemperatureByPostalCodeHandler) TemperatureByPostalCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTelTracer.Start(ctx, TemperatureByPostalCodeRequest)
	defer span.End()

	input := service.PostalCodeInput{
		PostalCode: chi.URLParam(r, "postalCode"),
	}

	output, err := h.Service.Execute(ctx, input)
	if err != nil {
		statusCodeError := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidPostalCode) {
			statusCodeError = http.StatusUnprocessableEntity
		}

		if errors.Is(err, service.ErrPostalCodeNotFound) {
			statusCodeError = http.StatusNotFound
		}

		http.Error(w, err.Error(), statusCodeError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
