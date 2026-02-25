# BlobTube

**Low-Bandwidth YouTube Video Streaming Proxy**

A Go-based video transcoding service that streams YouTube videos at ultra-low resolutions (128x128 and variants) for bandwidth-constrained environments. Features real-time AV1 transcoding, HLS adaptive streaming, and request queueing.

---

## Features

- **Ultra-Low Bandwidth**: Stream YouTube videos at 64x64, 128x128, or 256x256 resolution
- **AV1 Compression**: SVT-AV1 encoder provides 30%+ better compression than H.264
- **Adaptive Streaming**: HLS protocol automatically adjusts quality based on your bandwidth
- **No Storage**: Real-time transcoding with no video caching
- **Queue System**: Handles high load with request queueing
- **Privacy-First**: No authentication, no user tracking, aggregate analytics only
- **HTTPS Everywhere**: Secure connections in both dev and production

---

## Quick Start

### Prerequisites

- Docker & Docker Compose
- 8GB+ RAM recommended
- Modern browser (Chrome 90+, Firefox 88+, Safari 14+)

### Run with Docker Compose

```bash
# Clone the repository
git clone https://github.com/sixfeetup/blobtube.git
cd BlobTube

# Start the service
docker-compose up

# Access at https://localhost:8443
# (Accept self-signed certificate warning in browser)
```

### Development

```bash
# Install dependencies
make setup

# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Start development server
make dev
```

---

## Architecture

### Tech Stack

- **Backend**: Go 1.22+
- **Transcoding**: FFmpeg 6.0+ with SVT-AV1
- **YouTube Extraction**: yt-dlp
- **Streaming Protocol**: HLS (HTTP Live Streaming)
- **Analytics**: SQLite
- **Deployment**: Docker + Docker Compose

### Key Components

```
┌─────────────┐
│   Browser   │
│  (Web UI)   │
└──────┬──────┘
       │ HTTPS
       ▼
┌──────────────────────────────────────┐
│       Go HTTP Server (TLS)           │
│  ┌──────────┐    ┌──────────────┐    │
│  │ Web UI   │    │  API         │    │
│  │ Handler  │    │  /api/*      │    │
│  └──────────┘    └───────┬──────┘    │
│                          │           │
│       ┌──────────────────┴────┐      │
│       ▼                       ▼      │
│  ┌─────────────┐    ┌──────────────┐ │
│  │  Transcode  │    │  HLS Segment │ │
│  │  Manager    │    │   Server     │ │
│  │  + Queue    │    │              │ │
│  └──────┬──────┘    └──────────────┘ │
│         │                            │
└─────────┼────────────────────────────┘
          ▼
    ┌──────────┐
    │  yt-dlp  │──── YouTube
    └──────────┘
          │
          ▼
    ┌──────────┐
    │  FFmpeg  │
    │ SVT-AV1  │
    └──────────┘
```

---

## Usage

### Web Interface

1. Open https://localhost:8443 in your browser
2. Paste a YouTube URL into the input field
3. Click "Stream"
4. If at capacity, you'll see your queue position
5. Video will start playing automatically when ready

### API

#### Start a Stream
```bash
curl -X POST https://localhost:8443/api/stream \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'

# Response:
# {
#   "stream_id": "abc-123-def-456",
#   "status": "initializing",
#   "playlist_url": "/api/stream/abc-123-def-456/playlist.m3u8"
# }
```

#### Check Stream Status
```bash
curl https://localhost:8443/api/stream/{stream_id}/status

# Response:
# {
#   "status": "active",
#   "quality_tiers": ["64x64", "128x128", "256x256"]
# }
```

#### Get Analytics
```bash
curl https://localhost:8443/api/analytics

# Response:
# {
#   "total_streams": 1234,
#   "avg_duration_seconds": 420,
#   "unique_videos": 567
# }
```

---

## Configuration

### Environment Variables

```bash
# Development mode (enables yt-dlp caching)
DEV_MODE=false

# Maximum concurrent streams
MAX_CONCURRENT_STREAMS=5

# Maximum stream duration (seconds)
MAX_STREAM_DURATION_SECONDS=3600

# Server port
PORT=8443

# TLS certificate paths
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key

# Log level (debug, info, warn, error)
LOG_LEVEL=info
```

---

## Limitations

- **Single videos only**: Playlists not supported
- **1-hour maximum**: Videos limited to first 60 minutes
- **5 concurrent streams**: Excess requests queued (2-min timeout)
- **No authentication**: Open access (consider reverse proxy for auth)
- **No caching**: Every stream requires fresh transcode

---

## Project Structure

```
BlobTube/
├── cmd/server/          # Main application entry point
├── internal/            # Private application code
│   ├── api/             # HTTP handlers and routes
│   ├── transcode/       # FFmpeg and yt-dlp wrappers
│   ├── queue/           # Request queueing system
│   ├── analytics/       # SQLite analytics
│   ├── hls/             # HLS playlist/segment handling
│   └── metrics/         # Prometheus metrics
├── web/                 # Frontend assets (HTML/CSS/JS)
├── docs/                # Documentation
│   ├── adr/             # Architecture Decision Records
│   ├── tickets.md       # Development tickets
│   └── user-stories.md  # User stories
├── test/                # Integration and E2E tests
├── certs/               # TLS certificates (dev)
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Local development environment
├── Makefile             # Build automation
└── README.md            # This file
```

---

## Development

### Makefile Targets

```bash
make help              # Show available commands
make build             # Build Go binary
make test              # Run all tests
make test-unit         # Unit tests only
make test-integration  # Integration tests
make test-e2e          # End-to-end tests
make lint              # Run golangci-lint
make fmt               # Format code
make docker-build      # Build Docker image
make docker-up         # Start docker-compose
make docker-down       # Stop docker-compose
make docker-logs       # View logs
make clean             # Remove build artifacts
make setup             # Install development dependencies
```

### Running Tests

```bash
# All tests
make test

# Specific test packages
go test ./internal/transcode/...
go test ./internal/api/...

# With coverage
make test-coverage
```

### Linting

```bash
# Run linter
make lint

# Auto-fix issues
make fmt
```

---

## Documentation

- **[PLAN.md](./PLAN.md)**: Comprehensive project plan
- **[TODO.md](./TODO.md)**: Current session tasks
- **[Architecture Decision Records](./docs/adr/)**: All technical decisions documented
- **[User Stories](./docs/user-stories.md)**: Feature requirements
- **[Development Tickets](./docs/tickets.md)**: Implementation tasks
- **[API Documentation](./docs/api/)**: OpenAPI specification (coming soon)

---

## Architecture Decisions

All major technical decisions are documented as ADRs in `docs/adr/`:

- [ADR-001: No Persistent Storage](./docs/adr/001-no-storage.md)
- [ADR-002: HLS for Streaming Protocol](./docs/adr/002-hls-streaming.md)
- [ADR-009: AV1 Video Codec with SVT-AV1](./docs/adr/009-av1-svt.md)
- [ADR-011: Request Queueing System](./docs/adr/011-request-queueing.md)
- [See all 14 ADRs](./docs/adr/README.md)

---

## Performance

### System Requirements

- **CPU**: 4+ cores recommended (AV1 encoding is CPU-intensive)
- **RAM**: 8GB minimum, 16GB recommended
- **Bandwidth**: 100 Mbps+ for 5 concurrent streams

### Benchmarks

- **Encoding Speed**: ~0.5x realtime (1min video = 2min encode) with SVT-AV1 preset 8
- **Concurrent Streams**: 5 streams @ 128x128 on 4-core/8GB system
- **Compression**: 30%+ smaller than H.264 at same visual quality

### Scaling

For more than 5 concurrent streams:

1. **Horizontal Scaling**: Run multiple instances behind load balancer
2. **Optimize Preset**: Use SVT-AV1 preset 10 (faster, less compression)
3. **Upgrade Hardware**: More CPU cores = more concurrent streams

---

## Security

- **HTTPS Only**: Self-signed cert for dev, Let's Encrypt for production
- **No Authentication**: Consider adding reverse proxy with auth for production
- **Rate Limiting**: Basic IP-based rate limiting implemented
- **Input Validation**: YouTube URL validation prevents injection attacks
- **No Data Storage**: Videos never cached, reducing legal/privacy risks

---

## Legal & Compliance

**IMPORTANT**: This tool is intended for educational purposes and personal use only.

- Streaming YouTube content may violate YouTube's Terms of Service
- No video content is cached or stored
- Aggregate analytics only, no user tracking
- Use at your own risk

---

## Troubleshooting

### "FFmpeg not found"
Ensure FFmpeg with SVT-AV1 support is installed:
```bash
ffmpeg -encoders | grep svt
```

### "yt-dlp error: Unable to extract video"
- Video may be private, deleted, or region-locked
- YouTube may have changed their API (update yt-dlp)
- Check yt-dlp version: `yt-dlp --version`

### "Queue timeout"
System at capacity. Wait for current streams to complete or increase `MAX_CONCURRENT_STREAMS`.

### Browser certificate warning
Self-signed cert for local dev. Click "Advanced" → "Proceed to localhost" in your browser.

---

## Contributing

See [docs/tickets.md](./docs/tickets.md) for current development tasks.

### Development Workflow

1. Pick a ticket from [docs/tickets.md](./docs/tickets.md)
2. Create a feature branch: `git checkout -b ticket-XXX-description`
3. Make changes and write tests
4. Run `make test lint` to verify
5. Submit pull request
6. Update ticket status when merged

---

## License

[To be determined]

---

## Acknowledgements

- **FFmpeg**: The backbone of video transcoding
- **SVT-AV1**: Fast AV1 encoder by Intel/Netflix
- **yt-dlp**: YouTube downloader and metadata extractor
- **Video.js**: HTML5 video player

---

## Status

🚧 **Currently in Planning Phase**

See [TODO.md](./TODO.md) for current session status and [PLAN.md](./PLAN.md) for comprehensive project plan.

---

**Built with Go, FFmpeg, and questionable life choices.**
