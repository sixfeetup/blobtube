# API Documentation

**Status**: Placeholder - OpenAPI specification coming soon

This directory will contain the OpenAPI/Swagger specification for the BlobTube API.

## Planned Endpoints

### Stream Management

- `POST /api/stream` - Initialize a new stream
- `GET /api/stream/{id}/status` - Get stream status
- `GET /api/stream/{id}/playlist.m3u8` - Master HLS playlist
- `GET /api/stream/{id}/{quality}/playlist.m3u8` - Quality-specific playlist
- `GET /api/stream/{id}/{quality}/segment_{n}.ts` - Video segment

### Queue Management

- `GET /api/queue/{id}/status` - Get queue position and estimated wait time

### Analytics

- `GET /api/analytics` - Get aggregate statistics
- `GET /api/analytics/popular` - Get popular videos

### Health & Metrics

- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## Coming Soon

- Full OpenAPI 3.0 specification (openapi.yml)
- Interactive API documentation (Swagger UI)
- Example requests and responses
- Authentication documentation (if/when added)

## Development

See [TICKET-020](../tickets.md) for API documentation task.
