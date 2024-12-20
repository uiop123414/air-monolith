FROM golang:1.23

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./airApp ./cmd/api

FROM alpine:latest

COPY --from=0 /app/airApp /app/airApp

COPY --from=0 /app/internal/schemas /app/schemas

CMD ["/app/airApp"]