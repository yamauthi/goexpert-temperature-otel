FROM golang:1.22.3-alpine as builder
WORKDIR /app
COPY . .
RUN touch /app/cmd/input-service/.env
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o input-service ./cmd/input-service

FROM scratch
WORKDIR /app
COPY --from=builder /app/cmd/input-service/.env .
COPY --from=builder /app/input-service .

ENTRYPOINT ["/app/input-service"]