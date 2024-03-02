# syntax=docker/dockerfile:1.4
FROM golang:1.20-buster as build-env

RUN mkdir /app
WORKDIR /app

COPY --link ./ .
COPY --link go.mod .
COPY --link go.sum .

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.BuildVersion=$BUILD_VERSION" ./cmd/api/

# Deploy into alpine
FROM alpine:latest

WORKDIR /app
COPY --from=build-env --link /app/api .
COPY --from=build-env --link /app/config/config.yaml ./config/config.yaml

EXPOSE 80
CMD ["./api"]