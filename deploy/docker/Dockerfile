FROM golang:1.24.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . /app
RUN CGO_ENABLED=0 go build -o qbittorrent_exporter cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/qbittorrent_exporter /app/qbittorrent_exporter
COPY --from=builder /app/deploy/docker/config /app/config
ENTRYPOINT ["/app/qbittorrent_exporter"]
CMD ["-config", "/app/config/config.yaml"]
