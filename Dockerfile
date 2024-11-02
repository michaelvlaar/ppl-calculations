FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM debian:bookworm

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
