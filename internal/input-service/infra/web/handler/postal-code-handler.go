package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yamauthi/goexpert-temperature-otel/internal/input-service/application"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const PostalCodeInputRequest = "postal-code-input-request"

type PostalCodeInputHandler struct {
	Service    *application.PostalCodeInputService
	OTelTracer trace.Tracer
}

func NewPostalCodeInputHandler(service *application.PostalCodeInputService, oTelTracer trace.Tracer) *PostalCodeInputHandler {
	return &PostalCodeInputHandler{
		Service:    service,
		OTelTracer: oTelTracer,
	}
}

func (h *PostalCodeInputHandler) TemperatureByPostalCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, inputSpan := h.OTelTracer.Start(ctx, PostalCodeInputRequest)
	defer inputSpan.End()

	var input application.PostalCodeInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("unexpected error on input service handler: %s", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	output, err := h.Service.Execute(ctx, input)
	if err != nil {
		http.Error(w, err.Error(), output.StatusCode)
		return
	}

	err = json.NewEncoder(w).Encode(output.CityTemperature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
