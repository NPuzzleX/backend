# syntax=docker/dockerfile:1

FROM golang:1.19.0-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /npuzzlex-be-acc

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /npuzzlex-be-acc /npuzzlex-be-acc

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/npuzzlex-be-acc"]