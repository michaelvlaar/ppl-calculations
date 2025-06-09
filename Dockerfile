FROM golang:1.24-alpine AS build

ARG VERSION=dev
ENV VERSION=$VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache upx && \
    go build -ldflags="-w -s -X 'main.version=${VERSION}'" -o main . && \
    upx --best main

FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/michaelvlaar/ppl-calculations"
LABEL org.opencontainers.image.description="PPL Calculations is an open source project designed to simplify the calculation and management of weight and balance for an aeroclub fleet. It combines a reliable Go-based HTTP backend with a straightforward HTMX frontend, enabling pilots to perform quick and accurate weight and balance calculations. The project aims to improve operational safety and efficiency within aeroclubs and is open to community contributions."
LABEL org.opencontainers.image.licenses=MIT

COPY --from=build /app/main /app/main

RUN apk add --no-cache rsvg-convert && \
    addgroup -g 30001 appuser && \
    adduser -D -u 10001 -G appuser appuser

WORKDIR /app
USER appuser

EXPOSE 80

CMD ["./main"]
