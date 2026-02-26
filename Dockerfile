# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS builder

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY web ./web

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/blobtube ./cmd/server


FROM debian:bookworm-slim AS runtime

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    ca-certificates \
    ffmpeg \
    python3 \
    python3-pip \
    curl \
  && rm -rf /var/lib/apt/lists/* \
  && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
  && chmod a+rx /usr/local/bin/yt-dlp

# Verify libx264 encoder is available
RUN ffmpeg -hide_banner -encoders | grep -q libx264

WORKDIR /app

COPY --from=builder /out/blobtube /app/blobtube
COPY web /app/web
COPY scripts/docker-entrypoint.sh /app/docker-entrypoint.sh

RUN useradd --create-home --uid 10001 app \
  && chown -R app:app /app \
  && chmod +x /app/docker-entrypoint.sh

USER app

ENV PORT=8443 \
  HTTP_PORT=8080 \
  TLS_CERT_FILE=/certs/server.crt \
  TLS_KEY_FILE=/certs/server.key \
  STATIC_DIR=/app/web \
  LOG_LEVEL=info \
  DEV_MODE=false \
  YTDLP_PATH=yt-dlp

EXPOSE 8443 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
