FROM golang:1.22.3-alpine as builder
WORKDIR /app
COPY . .
RUN touch /app/cmd/temp-postal-code/.env
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o temperature-service ./cmd/temp-postal-code

FROM scratch
WORKDIR /app
COPY --from=builder /app/cmd/temp-postal-code/.env .
COPY --from=builder /app/temperature-service .

ENTRYPOINT ["/app/temperature-service"]