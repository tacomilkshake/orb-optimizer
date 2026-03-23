FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /orb-optimizer .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl procps && rm -rf /var/lib/apt/lists/*
COPY --from=build /orb-optimizer /usr/local/bin/orb-optimizer
WORKDIR /data
ENTRYPOINT ["orb-optimizer"]
CMD ["collect"]
