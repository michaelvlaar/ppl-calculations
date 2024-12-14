FROM golang:1.23-alpine AS build

ARG VERSION=dev
ENV VERSION=$VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-X 'main.version=${VERSION}'" -o main .

FROM debian:bookworm

LABEL org.opencontainers.image.source="https://github.com/michaelvlaar/ppl-calculations"
LABEL org.opencontainers.image.description="PPL Calculations is an open source project designed to simplify the calculation and management of weight and balance for an aeroclub fleet. It combines a reliable Go-based HTTP backend with a straightforward HTMX frontend, enabling pilots to perform quick and accurate weight and balance calculations. The project aims to improve operational safety and efficiency within aeroclubs and is open to community contributions."
LABEL org.opencontainers.image.licenses=MIT

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    librsvg2-bin \
    texlive-xetex \
    texlive-fonts-recommended \
    texlive-fonts-extra \
    fonts-roboto \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r appuser && useradd -r -g appuser appuser
RUN mkdir -p /tmp && chown -R appuser:appuser /tmp
COPY --from=build /app/main /app/main

WORKDIR /app
USER appuser

EXPOSE 80

CMD ["./main"]
