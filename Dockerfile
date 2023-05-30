# syntax=docker/dockerfile:1
# A sample microservice in Go packaged into a container image.
# from https://docs.docker.com/language/golang/build-images/

# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/api_server/*.go ./cmd/api_server/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api_server cmd/api_server/*.go
# TODO test
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /api_server /api_server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api_server"]


