# syntax=docker/dockerfile:1

# Build environment
# -----------------
FROM golang:1.23-alpine as build-env
WORKDIR /incrate

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api

# Deployment environment
# ----------------------
FROM alpine

COPY --from=build-env /incrate/bin/api /incrate/

EXPOSE 8080
CMD ["/incrate/api"]