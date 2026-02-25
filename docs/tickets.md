# Development Tickets

**Project**: TicketShit (Low-Bandwidth YouTube Streaming)
**Last Updated**: 2026-02-25

This document tracks all development tickets organized by phase and milestone.

## Status Legend
- **TODO**: Not started
- **IN_PROGRESS**: Currently being worked on
- **BLOCKED**: Waiting on dependency
- **DONE**: Completed and merged

## Summary

**Total Tickets**: 28
**Estimated Total**: ~100 hours
**Completed**: 0
**In Progress**: 0

---

# Phase 1: Foundation (Milestone 1)

**Goal**: Basic project infrastructure and scaffolding

## TICKET-001: Project Scaffolding
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 2h

### Description
Initialize Go module, create directory structure, setup build tooling.

### Tasks
- [ ] Run `go mod init github.com/yourusername/TicketShit`
- [ ] Create directory structure (cmd/, internal/, docs/, web/, test/)
- [ ] Setup Makefile with targets: build, test, lint, docker-build, docker-up
- [ ] Add .gitignore for Go + Docker
- [ ] Create basic README.md
- [ ] Setup golangci-lint configuration

### Dependencies
None

### Acceptance Criteria
- `make build` compiles successfully
- Directory structure matches plan
- golangci-lint configured

---

## TICKET-002: Docker Infrastructure
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Create Docker setup with FFmpeg (SVT-AV1) and yt-dlp.

### Tasks
- [ ] Create multi-stage Dockerfile (builder + runtime)
- [ ] Install FFmpeg with `--enable-libsvtav1`
- [ ] Verify SVT-AV1: `ffmpeg -encoders | grep svt`
- [ ] Install yt-dlp
- [ ] Create docker-compose.yml with HTTPS support
- [ ] Generate self-signed certificate for dev
- [ ] Setup environment variable configuration
- [ ] Document Docker usage in README

### Dependencies
- TICKET-001

### Acceptance Criteria
- `docker-compose up` starts successfully
- FFmpeg has SVT-AV1 support
- Service accessible at https://localhost:8443

---

## TICKET-003: Basic HTTP Server
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Implement Go HTTP server with TLS, routing, and graceful shutdown.

### Tasks
- [ ] Setup HTTP server with chi or gorilla/mux
- [ ] Implement TLS/HTTPS support
- [ ] Create health check endpoint: GET /health
- [ ] Setup static file serving for web UI
- [ ] Configure structured logging with zerolog
- [ ] Implement graceful shutdown handling (SIGTERM)

### Dependencies
- TICKET-002

### Acceptance Criteria
- Server responds to https://localhost:8443/health
- Graceful shutdown on Ctrl+C
- Logs are structured JSON

---

## TICKET-004: yt-dlp Integration
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Create wrapper for yt-dlp to extract YouTube video stream URLs.

### Tasks
- [ ] Create ytdlp package in internal/transcode/
- [ ] Implement Execute() function wrapping yt-dlp
- [ ] Parse JSON output (`yt-dlp -j`)
- [ ] Extract direct stream URL
- [ ] Handle errors (invalid URL, private video, region lock)
- [ ] Write unit tests with mocked yt-dlp responses
- [ ] Implement debug caching (DEV_MODE only, ADR-013)

### Dependencies
- TICKET-003

### Acceptance Criteria
- Successfully extracts stream URL from YouTube video
- Handles error cases gracefully
- Unit tests pass

---

# Phase 2: Transcoding Engine (Milestone 2)

**Goal**: Core FFmpeg transcoding with HLS output

## TICKET-005: FFmpeg Pipeline Foundation
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Design and implement FFmpeg command for SVT-AV1 HLS transcoding.

### Tasks
- [ ] Design FFmpeg command: `-c:v libsvtav1 -preset 8 -crf 35`
- [ ] Test: YouTube URL → 128x128 AV1 HLS segments
- [ ] Implement exec wrapper in internal/transcode/ffmpeg.go
- [ ] Capture stdout/stderr for debugging
- [ ] Enforce 1-hour max duration: `-t 3600`
- [ ] Write unit tests

### Dependencies
- TICKET-004

### Acceptance Criteria
- FFmpeg successfully transcodes test video
- HLS segments generated in /tmp
- 1-hour limit enforced

---

## TICKET-006: Multi-Quality Transcoding
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Generate 3 quality tiers (64x64, 128x128, 256x256) in parallel.

### Tasks
- [ ] Spawn 3 parallel FFmpeg processes
- [ ] Configure bitrates: 50k, 100k, 200k
- [ ] Coordinate segment generation across qualities
- [ ] Handle individual process failures
- [ ] Test concurrent transcoding

### Dependencies
- TICKET-005

### Acceptance Criteria
- 3 quality tiers generated simultaneously
- Segments sync across qualities
- One quality failing doesn't break others

---

## TICKET-007: HLS Playlist Generation
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Generate master and media playlists for HLS streaming.

### Tasks
- [ ] Integrate github.com/grafov/m3u8 library
- [ ] Generate master playlist listing 3 variants
- [ ] Generate media playlists per quality
- [ ] Update playlists as segments are produced
- [ ] Handle playlist serving via HTTP

### Dependencies
- TICKET-006

### Acceptance Criteria
- Valid HLS master playlist
- Media playlists list segments correctly
- Playlists update dynamically

---

## TICKET-008: Segment Serving
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Implement HTTP endpoints to serve HLS segments and playlists.

### Tasks
- [ ] Endpoint: GET /api/stream/{id}/{quality}/segment_{n}.ts
- [ ] Stream segments as generated (don't wait for completion)
- [ ] Handle missing segments (404 vs 202 "not ready yet")
- [ ] Add CORS headers for browser access
- [ ] Implement byte-range request support

### Dependencies
- TICKET-007

### Acceptance Criteria
- Segments served via HTTP
- Browser can fetch and play segments
- Handles missing segments gracefully

---

# Phase 3: Stream Management (Milestone 3)

**Goal**: Stream lifecycle and resource management

## TICKET-009: Stream Lifecycle
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Implement stream state machine and lifecycle management.

### Tasks
- [ ] Generate unique stream IDs (UUIDv4)
- [ ] Implement state machine: initializing → active → completed/error
- [ ] Track stream metadata (start time, quality, status)
- [ ] Implement 5-minute inactivity timeout
- [ ] Create API: GET /api/stream/{id}/status

### Dependencies
- TICKET-008

### Acceptance Criteria
- Stream transitions through states correctly
- Inactive streams timeout after 5 minutes
- Status API returns accurate state

---

## TICKET-010: Resource Cleanup
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Implement cleanup of FFmpeg processes and temporary files.

### Tasks
- [ ] Kill FFmpeg processes on stream completion/timeout
- [ ] Remove temporary segment files
- [ ] Handle cleanup on server shutdown (SIGTERM)
- [ ] Prevent zombie processes
- [ ] Log cleanup operations

### Dependencies
- TICKET-009

### Acceptance Criteria
- No orphaned FFmpeg processes
- Temp files removed after stream ends
- Graceful cleanup on shutdown

---

## TICKET-011: Error Handling
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Standardized error handling and rate limiting.

### Tasks
- [ ] Standardized JSON error responses
- [ ] Handle YouTube extraction failures
- [ ] Handle FFmpeg transcode failures
- [ ] Basic rate limiting (max N concurrent streams)
- [ ] Error logging and metrics

### Dependencies
- TICKET-010

### Acceptance Criteria
- Consistent error response format
- Rate limiting enforced
- Errors logged with context

---

## TICKET-011-A: Request Queueing System
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Implement queue for requests when at capacity (ADR-011).

### Tasks
- [ ] Implement queue data structure (channel-based)
- [ ] Queue when concurrent streams reach limit (5 streams)
- [ ] Track queue position and estimated wait time
- [ ] 2-minute timeout for queued requests
- [ ] API: GET /api/queue/{id}/status
- [ ] Automatic promotion when slot available

### Dependencies
- TICKET-011

### Acceptance Criteria
- Requests queued when at capacity
- Queue position tracked accurately
- Auto-start when slot available
- Timeout removes stale requests

---

## TICKET-011-B: Analytics System
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
SQLite-based analytics for aggregate stats (ADR-012).

### Tasks
- [ ] Create SQLite schema for stream_events
- [ ] Track: video URL hash, view count, duration, timestamp
- [ ] Implement endpoints: GET /api/analytics, GET /api/analytics/popular
- [ ] Create simple dashboard page
- [ ] Ensure no user tracking (aggregate only)
- [ ] Write to DB on stream completion

### Dependencies
- TICKET-010

### Acceptance Criteria
- Analytics data persisted to SQLite
- API returns summary and popular videos
- Privacy-preserving (hashed URLs only)

---

## TICKET-011-C: HTTPS Setup
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Configure TLS/HTTPS for development and production (ADR-010).

### Tasks
- [ ] Create script to generate self-signed cert for dev
- [ ] Configure Go HTTP server for TLS
- [ ] Update docker-compose with cert volumes
- [ ] Document production cert setup (Let's Encrypt)
- [ ] Implement HTTP → HTTPS redirect
- [ ] Test with browser (accept self-signed cert)

### Dependencies
- TICKET-002

### Acceptance Criteria
- Dev server runs on https://localhost:8443
- Self-signed cert generation automated
- Production cert setup documented

---

## TICKET-011-D: Debug Caching (Dev Only)
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 2h

### Description
Implement in-memory caching of yt-dlp results for dev (ADR-013).

### Tasks
- [ ] Cache yt-dlp results in memory (5 min TTL)
- [ ] Only enabled when DEV_MODE=true
- [ ] Log cache hits for debugging
- [ ] Document in README that this is dev-only
- [ ] Implement cache expiration

### Dependencies
- TICKET-004

### Acceptance Criteria
- Caching works in dev mode only
- 5-minute TTL enforced
- Cache disabled in production

---

# Phase 4: Web UI (Milestone 4)

**Goal**: User-facing web interface

## TICKET-012: Basic Web Interface
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Create web UI with Video.js player and queue status.

### Tasks
- [ ] Create HTML page with YouTube URL input
- [ ] Integrate Video.js HLS player
- [ ] Handle form submit → POST /api/stream
- [ ] Display queue status (position, wait time)
- [ ] Show loading states and error messages
- [ ] Ensure HTTPS-only (no HTTP fallback)

### Dependencies
- TICKET-011-A

### Acceptance Criteria
- User can paste URL and watch video
- Queue status displayed when at capacity
- Auto-transition from queue to playback

---

## TICKET-013: Player Enhancements
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Enhance video player with quality selector and controls.

### Tasks
- [ ] Add quality selector UI
- [ ] Display current quality indicator
- [ ] Customize playback controls (play/pause/seek/volume)
- [ ] Mobile-responsive design
- [ ] Test on multiple browsers

### Dependencies
- TICKET-012

### Acceptance Criteria
- Quality selector functional
- Mobile-friendly UI
- Works on Chrome, Firefox, Safari

---

## TICKET-014: UX Improvements
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Polish user experience with thumbnails, messaging, sharing.

### Tasks
- [ ] Display preview thumbnail (if available from yt-dlp)
- [ ] Show estimated wait time before playback
- [ ] Error messages for common issues
- [ ] Share/copy stream URL functionality
- [ ] Add favicon and page title

### Dependencies
- TICKET-013

### Acceptance Criteria
- Thumbnail shown when available
- Clear error messages
- Shareable stream URLs

---

# Phase 5: Observability (Milestone 5)

**Goal**: Metrics, monitoring, logging

## TICKET-015: Metrics & Monitoring
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Implement Prometheus metrics endpoint.

### Tasks
- [ ] Create /metrics endpoint
- [ ] Metrics: active_streams, transcode_errors, segments_served, request_duration, queue_depth
- [ ] Health check includes FFmpeg/yt-dlp availability
- [ ] Document metrics in README

### Dependencies
- TICKET-011-B

### Acceptance Criteria
- /metrics endpoint returns Prometheus format
- Key metrics tracked
- Health check comprehensive

---

## TICKET-016: Logging & Debugging
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 2h

### Description
Structured logging and debug support.

### Tasks
- [ ] Structured logging for all operations (zerolog)
- [ ] Log correlation with stream ID
- [ ] Debug mode with verbose FFmpeg output
- [ ] Log rotation in Docker
- [ ] Log levels configurable via env var

### Dependencies
- TICKET-003

### Acceptance Criteria
- All logs structured JSON
- Stream ID in all related logs
- Debug mode shows FFmpeg output

---

# Phase 6: Testing & Quality (Milestone 6)

**Goal**: Comprehensive test coverage

## TICKET-017: Unit Tests
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 6h

### Description
Write unit tests for core components.

### Tasks
- [ ] Test yt-dlp wrapper (mocked)
- [ ] Test stream lifecycle manager
- [ ] Test HLS playlist generation
- [ ] Test cleanup logic
- [ ] Test queue implementation
- [ ] Achieve 80%+ coverage

### Dependencies
None (can start anytime)

### Acceptance Criteria
- 80%+ test coverage
- All tests pass
- Mocking used appropriately

---

## TICKET-018: Integration Tests
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Write integration tests for API endpoints.

### Tasks
- [ ] Test API endpoints end-to-end
- [ ] Mock YouTube/FFmpeg where appropriate
- [ ] Test concurrent streams
- [ ] Test cleanup and timeouts
- [ ] Test queue behavior under load

### Dependencies
- TICKET-012

### Acceptance Criteria
- API integration tests pass
- Concurrent scenarios tested
- Cleanup verified

---

## TICKET-019: E2E Test
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
End-to-end test of full stream lifecycle.

### Tasks
- [ ] Full flow: submit URL → transcode → serve → playback
- [ ] Use stable test YouTube video
- [ ] Verify HLS segments are valid
- [ ] Automated E2E test script
- [ ] Run in CI

### Dependencies
- TICKET-018

### Acceptance Criteria
- E2E test passes reliably
- Tests real YouTube video
- Automated in CI

---

## TICKET-020: Documentation
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Comprehensive project documentation.

### Tasks
- [ ] Update README with setup, usage, architecture
- [ ] Create OpenAPI spec for API
- [ ] Ensure all ADRs are complete
- [ ] Write contributing guide
- [ ] Document deployment process

### Dependencies
- TICKET-019

### Acceptance Criteria
- README comprehensive
- API documented (OpenAPI)
- Contributing guide available

---

# Phase 7: Production Readiness (Milestone 7)

**Goal**: Performance, security, deployment

## TICKET-021: Performance Optimization
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 5h

### Description
Benchmark and optimize for production.

### Tasks
- [ ] Benchmark concurrent stream capacity
- [ ] Optimize FFmpeg flags for speed
- [ ] Memory profiling and leak detection
- [ ] CPU profiling
- [ ] Load testing

### Dependencies
- TICKET-019

### Acceptance Criteria
- Know max concurrent streams
- No memory leaks
- Optimized FFmpeg settings

---

## TICKET-022: Security Hardening
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 4h

### Description
Security review and hardening.

### Tasks
- [ ] Input validation (YouTube URL format)
- [ ] Rate limiting per IP
- [ ] Resource limits (max stream duration enforced)
- [ ] Dependency vulnerability scan
- [ ] Security headers (CSP, HSTS, etc.)

### Dependencies
- TICKET-021

### Acceptance Criteria
- Input validation comprehensive
- Rate limiting per IP works
- No critical vulnerabilities

---

## TICKET-023: Deployment Guide
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Production deployment documentation.

### Tasks
- [ ] Create production docker-compose.yml
- [ ] Document reverse proxy setup (nginx/traefik)
- [ ] SSL/TLS configuration (Let's Encrypt)
- [ ] Scaling guide (multiple instances)
- [ ] Monitoring setup guide

### Dependencies
- TICKET-022

### Acceptance Criteria
- Production deployment documented
- Scaling strategy documented
- Reverse proxy example provided

---

## TICKET-024: CI/CD Pipeline
**Status**: TODO | **Assignee**: Unassigned | **Estimate**: 3h

### Description
Automated CI/CD with GitHub Actions.

### Tasks
- [ ] Create GitHub Actions workflow
- [ ] Automated tests on PR
- [ ] Docker image build and push
- [ ] Linting enforcement
- [ ] Deployment automation (optional)

### Dependencies
- TICKET-023

### Acceptance Criteria
- Tests run on every PR
- Docker images built automatically
- Linting enforced

---

## Ticket Dependencies Graph

```
TICKET-001 (Scaffolding)
    └─> TICKET-002 (Docker)
        ├─> TICKET-003 (HTTP Server)
        │   └─> TICKET-004 (yt-dlp)
        │       └─> TICKET-005 (FFmpeg Pipeline)
        │           └─> TICKET-006 (Multi-Quality)
        │               └─> TICKET-007 (HLS Playlists)
        │                   └─> TICKET-008 (Segment Serving)
        │                       └─> TICKET-009 (Stream Lifecycle)
        │                           └─> TICKET-010 (Cleanup)
        │                               ├─> TICKET-011 (Error Handling)
        │                               │   └─> TICKET-011-A (Queue)
        │                               │       └─> TICKET-012 (Web UI)
        │                               └─> TICKET-011-B (Analytics)
        └─> TICKET-011-C (HTTPS)

TICKET-004 ─> TICKET-011-D (Debug Cache)

TICKET-012 ─> TICKET-013 (Player) ─> TICKET-014 (UX)
TICKET-011-B ─> TICKET-015 (Metrics)
TICKET-003 ─> TICKET-016 (Logging)

TICKET-012 ─> TICKET-018 (Integration Tests)
TICKET-018 ─> TICKET-019 (E2E) ─> TICKET-020 (Docs)
TICKET-019 ─> TICKET-021 (Perf) ─> TICKET-022 (Security) ─> TICKET-023 (Deploy) ─> TICKET-024 (CI/CD)

TICKET-017 (Unit Tests) - No dependencies, can start anytime
```

---

## Notes for Developers

- **Update this document** when starting/completing tickets
- **Add notes** to tickets as you work
- **Update estimates** if they were significantly off
- **Link commits/PRs** to tickets for traceability
