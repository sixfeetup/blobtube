# Low-Bandwidth YouTube Video Streaming Project Plan

**Project Name**: TicketShit (working name - to be renamed)
**Status**: Planning Phase
**Last Updated**: 2026-02-25

## Executive Summary

Build a Go-based video transcoding proxy that streams YouTube videos at ultra-low resolutions (128x128 and variants) for bandwidth-constrained environments. The system will transcode videos on-the-fly using SVT-AV1 codec without storage, using HLS protocol with adaptive bitrate streaming. Features request queueing for load management, basic analytics, and HTTPS throughout. Designed for scalability with Docker/Docker Compose infrastructure.

## Core Requirements

- **Language**: Go
- **Video Codec**: SVT-AV1 (preset 8 for speed)
- **Primary Resolution**: 128x128 pixels (with adaptive: 64x64, 128x128, 256x256)
- **Authentication**: None
- **Storage**: None for videos; SQLite for analytics only
- **Streaming Protocol**: HLS (HTTP Live Streaming)
- **Transport**: HTTPS only (self-signed for dev, Let's Encrypt for prod)
- **Deployment**: Docker/Docker Compose
- **Quality**: Adaptive based on bandwidth
- **Interface**: Web UI + REST API
- **Infrastructure**: Makefile for common commands, tests, linting
- **Limitations**: Single videos only (no playlists), 1-hour max duration
- **Queueing**: Queue requests when at capacity (5 concurrent streams)
- **Analytics**: Basic aggregate stats (no user tracking)

## Architecture Decision Records (ADRs)

### ADR-001: No Persistent Storage
**Decision**: All transcoding happens in-memory or via temporary streams. No video caching.
**Rationale**: Simplifies infrastructure, avoids legal gray areas with YouTube content, reduces operational costs.
**Consequences**: Higher CPU usage per request, can't serve same video to multiple users from cache.
**Status**: Accepted

### ADR-002: HLS for Streaming Protocol
**Decision**: Use HTTP Live Streaming (HLS) with segmented transport streams.
**Rationale**: Industry standard, wide browser support, enables adaptive bitrate, supports seeking.
**Consequences**: More complex than progressive download, requires segment generation, but provides superior UX.
**Status**: Accepted

### ADR-003: Go for Backend
**Decision**: Use Go as primary language for transcoding service.
**Rationale**: Excellent concurrency model, strong stdlib, good FFmpeg bindings, single binary deployment.
**Consequences**: Team needs Go expertise, fewer video processing libraries than Python/Node.
**Status**: Accepted

### ADR-004: FFmpeg for Transcoding
**Decision**: Use FFmpeg via Go bindings or exec for video processing.
**Rationale**: Industry standard, supports all codecs, battle-tested, excellent performance.
**Consequences**: FFmpeg must be available in Docker containers, adds dependency.
**Status**: Accepted

### ADR-005: yt-dlp for YouTube Extraction
**Decision**: Use yt-dlp to extract YouTube video streams.
**Rationale**: Most reliable YouTube downloader, actively maintained, handles auth/region restrictions.
**Consequences**: External Python dependency, may break with YouTube changes.
**Status**: Accepted

### ADR-006: Stateless Architecture
**Decision**: No session state, no user tracking, ephemeral processing.
**Rationale**: Simplifies scaling, aligns with no-auth requirement, reduces privacy concerns.
**Consequences**: Can't implement user preferences or rate limiting per user.
**Status**: Accepted

### ADR-007: Docker Compose for Local Development
**Decision**: Use Docker Compose for local development environment.
**Rationale**: Consistent environments, easy onboarding, matches production setup.
**Consequences**: Requires Docker knowledge, slightly slower iteration than native development.
**Status**: Accepted

### ADR-008: Adaptive Bitrate Streaming
**Decision**: Generate multiple quality tiers (64x64, 128x128, 256x256) for HLS ABR.
**Rationale**: Better user experience, handles varying network conditions.
**Consequences**: 3x transcoding load per video, more complex playlist generation.
**Status**: Accepted

### ADR-009: AV1 Video Codec with SVT-AV1
**Decision**: Use AV1 codec with SVT-AV1 encoder, prioritizing speed (cpu-used 8).
**Rationale**: Superior compression (~30% better than H.264) is critical for low-bandwidth use case. SVT-AV1 chosen over libaom for faster encoding speed. cpu-used 8 preset prioritizes encoding speed to maximize concurrent streams while maintaining good compression.
**Consequences**:
- SVT-AV1 is 2-3x faster than libaom but still slower than H.264
- cpu-used 8 trades some compression quality for speed
- Higher CPU usage during transcode than H.264
- Limited hardware decoding support (more battery usage on mobile)
- Requires FFmpeg compiled with libsvtav1
- Modern browser support is sufficient (Chrome 90+, Firefox 88+, Edge 90+)
**Status**: Accepted
**Implementation**: FFmpeg flags: `-c:v libsvtav1 -preset 8 -crf 35`

### ADR-010: HTTPS for All Environments
**Decision**: Require HTTPS for both local development and production.
**Rationale**: Modern browsers require HTTPS for many APIs (Service Workers, clipboard, etc.). Self-signed certs acceptable for local dev.
**Consequences**: Need to generate/manage TLS certificates, slightly more complex local setup.
**Status**: Accepted

### ADR-011: Request Queueing System
**Decision**: Implement queue for stream requests when at capacity, no H.264 fallback.
**Rationale**: AV1 encoding is CPU-intensive. Queueing provides better UX than rejection and maintains codec consistency.
**Consequences**: Need queue implementation, position tracking, timeout handling. Users may wait during high load.
**Status**: Accepted

### ADR-012: Basic Analytics
**Decision**: Track basic metrics: view count, popular videos, stream duration.
**Rationale**: Useful for understanding usage patterns and optimizing system.
**Consequences**: Need lightweight storage (SQLite or in-memory with periodic export). No user tracking, only aggregate stats.
**Status**: Accepted

### ADR-013: Debug-Only Caching
**Decision**: Cache yt-dlp results briefly (5 min) for debugging only, not production.
**Rationale**: Speeds up development iteration when testing same video repeatedly.
**Consequences**: Dev environment only, clear documentation that this is debug-only.
**Status**: Accepted

### ADR-014: Maximum Stream Duration
**Decision**: Enforce 1-hour maximum stream duration.
**Rationale**: Prevents resource exhaustion, covers vast majority of use cases.
**Consequences**: Long videos will be truncated, need clear user messaging.
**Status**: Accepted

## System Architecture

### High-Level Components

```
┌─────────────┐
│   Browser   │
│  (Web UI)   │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────────────────────┐
│       Go HTTP Server                │
│  ┌──────────┐    ┌──────────────┐  │
│  │ Web UI   │    │  API Handler │  │
│  │ (static) │    │  (/api/*)    │  │
│  └──────────┘    └───────┬──────┘  │
│                          │          │
│       ┌──────────────────┴────┐    │
│       ▼                       ▼     │
│  ┌─────────────┐    ┌──────────────┐
│  │  Transcode  │    │  HLS Segment │
│  │  Manager    │    │   Server     │
│  └──────┬──────┘    └──────────────┘
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
    │ Pipeline │
    └──────────┘
```

### Component Responsibilities

1. **Web UI Handler**
   - Serve static HTML/CSS/JS
   - Simple form: YouTube URL input + play button
   - Video.js player for HLS playback

2. **API Handler**
   - `POST /api/stream` - Initialize stream, return HLS playlist URL
   - `GET /api/stream/{id}/playlist.m3u8` - Master HLS playlist
   - `GET /api/stream/{id}/{quality}/playlist.m3u8` - Quality-specific playlist
   - `GET /api/stream/{id}/{quality}/segment_{n}.ts` - Video segments

3. **Transcode Manager**
   - Accepts YouTube URL
   - Calls yt-dlp to get stream URL
   - Spawns FFmpeg processes for each quality tier
   - Manages temporary segment storage (in-memory or /tmp)
   - Tracks active streams, cleanup on completion/timeout

4. **HLS Segment Server**
   - Serves .m3u8 playlists and .ts segments
   - Handles byte-range requests
   - Implements CORS for browser access

5. **Stream Lifecycle Manager**
   - Assigns unique stream IDs
   - Tracks stream state (initializing, active, completed, error)
   - Cleanup: kills FFmpeg processes, removes temp files
   - Timeout handling (streams expire after inactivity, max 1 hour)

6. **Request Queue**
   - Queues stream requests when at capacity (5 concurrent streams)
   - Tracks queue position and estimated wait time
   - 2-minute timeout for queued requests
   - Promotes queued requests when slots available

7. **Analytics Store**
   - SQLite database for aggregate statistics
   - Records: video URL hash, view count, duration, timestamp
   - No user tracking, aggregate only
   - Provides summary and popular videos endpoints

## Technology Stack

### Core
- **Go 1.22+**: Backend service
- **FFmpeg 6.0+**: Video transcoding (must include libsvtav1 for AV1 encoding)
- **SVT-AV1**: Fast AV1 encoder (preset 8 for speed)
- **yt-dlp**: YouTube video extraction
- **SQLite** or **in-memory store**: Lightweight analytics storage

### Go Libraries (Proposed)
- `net/http`: HTTP server with TLS (stdlib)
- `gorilla/mux` or `chi`: HTTP routing
- `github.com/grafov/m3u8`: HLS playlist generation
- `github.com/rs/zerolog`: Structured logging
- `github.com/mattn/go-sqlite3`: SQLite driver for analytics
- `crypto/tls`: HTTPS/TLS support (stdlib)
- Testing: `testing` stdlib + `testify/assert`

### Frontend
- **Video.js**: HLS video player
- Vanilla JS (no framework for simplicity)
- Minimal CSS (responsive, mobile-friendly)

### Infrastructure
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration
- **Make**: Build automation
- **golangci-lint**: Go linting
- **GitHub Actions** (optional): CI/CD

## User Stories

### Epic 1: Core Video Streaming

**US-001**: As a user, I want to paste a YouTube URL and watch a low-bandwidth version, so I can view videos on slow connections.
- Acceptance Criteria:
  - Web page has input field for YouTube URL
  - Submit button triggers video playback
  - Video plays in 128x128 resolution
  - Playback starts within 5 seconds

**US-002**: As a user, I want the video to adapt to my bandwidth, so I get the best quality my connection can handle.
- Acceptance Criteria:
  - Player automatically switches between 64x64, 128x128, 256x256
  - No buffering for more than 2 seconds on quality switch
  - Quality indicator visible in player

**US-003**: As a user, I want to seek through the video, so I can skip to specific parts.
- Acceptance Criteria:
  - Seek bar functional
  - Seeking works without re-transcoding entire video
  - Seeking completes within 3 seconds

### Epic 2: API Access

**US-004**: As a developer, I want to use an API to get video streams, so I can integrate with other applications.
- Acceptance Criteria:
  - POST /api/stream with {"url": "youtube_url"} returns stream ID
  - GET /api/stream/{id}/playlist.m3u8 returns valid HLS playlist
  - API documented with OpenAPI/Swagger spec

**US-005**: As a developer, I want to query stream status, so I can handle errors gracefully.
- Acceptance Criteria:
  - GET /api/stream/{id}/status returns state (initializing, ready, error)
  - Error responses include helpful messages
  - 404 for expired/invalid stream IDs

**US-003-A**: As a user, I want to see my position in queue when system is busy, so I know when my video will start.
- Acceptance Criteria:
  - Message displays "Position X of Y in queue"
  - Estimated wait time shown
  - Queue position updates in real-time
  - When stream starts, automatically transitions to player

### Epic 3: Operations & Reliability

**US-006**: As an operator, I want to deploy via Docker Compose, so I can run the service consistently.
- Acceptance Criteria:
  - `docker-compose up` starts all services
  - Service accessible at http://localhost:8080
  - Logs visible via `docker-compose logs`

**US-007**: As an operator, I want to monitor active streams, so I can understand system load.
- Acceptance Criteria:
  - Prometheus metrics endpoint (/metrics)
  - Metrics: active streams, transcode errors, segment generation rate, queue depth
  - CPU/memory usage visible

**US-008**: As an operator, I want streams to auto-cleanup, so the system doesn't leak resources.
- Acceptance Criteria:
  - Streams inactive for 5 minutes are terminated
  - Streams exceeding 1 hour are terminated
  - FFmpeg processes killed on cleanup
  - Temp files removed
  - Graceful shutdown on SIGTERM

**US-008-A**: As an operator, I want to see basic analytics, so I understand system usage.
- Acceptance Criteria:
  - Dashboard shows: total streams, popular videos, average duration
  - Aggregate data only (no user tracking)
  - Data persisted to SQLite
  - Endpoint: GET /api/analytics
- Acceptance Criteria:
  - Streams inactive for 5 minutes are terminated
  - FFmpeg processes killed on cleanup
  - Temp files removed
  - Graceful shutdown on SIGTERM

### Epic 4: Quality & Testing

**US-009**: As a developer, I want comprehensive tests, so I can refactor confidently.
- Acceptance Criteria:
  - Unit tests for core logic (>80% coverage)
  - Integration tests for API endpoints
  - E2E test: full stream lifecycle
  - Tests run in CI

**US-010**: As a developer, I want code quality checks, so the codebase stays maintainable.
- Acceptance Criteria:
  - golangci-lint passes with strict config
  - Go fmt, go vet enforced
  - Pre-commit hooks available
  - Make target: `make lint`

## Development Tickets

### Phase 1: Foundation (Milestone 1)

**TICKET-001**: Project scaffolding
- Initialize Go module
- Create directory structure
- Setup Makefile with targets: build, test, lint, docker-build, docker-up
- Add .gitignore, README.md
- Setup golangci-lint configuration
- **Estimate**: 2 hours

**TICKET-002**: Docker infrastructure
- Create Dockerfile for Go service
- Multi-stage build (builder + runtime)
- Install FFmpeg with SVT-AV1 support and yt-dlp in container
- Verify SVT-AV1 encoding capability (`ffmpeg -encoders | grep svt`)
- Create docker-compose.yml with HTTPS support (self-signed cert for dev)
- Environment variable configuration
- **Estimate**: 5 hours (increased due to SVT-AV1 + HTTPS)

**TICKET-003**: Basic HTTP server
- Implement HTTP server with chi/gorilla
- Health check endpoint: GET /health
- Static file serving for web UI
- Structured logging with zerolog
- Graceful shutdown handling
- **Estimate**: 3 hours

**TICKET-004**: yt-dlp integration
- Create wrapper for yt-dlp execution
- Extract direct video stream URL from YouTube link
- Handle errors (invalid URL, private video, region lock)
- Unit tests with mocked yt-dlp responses
- **Estimate**: 4 hours

### Phase 2: Transcoding Engine (Milestone 2)

**TICKET-005**: FFmpeg pipeline foundation
- Design FFmpeg command for HLS output with SVT-AV1 encoding
- Use SVT-AV1 preset 8 (speed priority): `-c:v libsvtav1 -preset 8 -crf 35`
- Test commands: YouTube URL → 128x128 AV1 HLS segments
- Implement exec wrapper for FFmpeg
- Capture stdout/stderr for debugging
- Enforce 1-hour max duration with `-t 3600`
- **Estimate**: 4 hours

**TICKET-006**: Multi-quality transcoding
- Generate 3 quality tiers: 64x64, 128x128, 256x256
- Spawn parallel FFmpeg processes
- Coordinate segment generation across qualities
- **Estimate**: 5 hours

**TICKET-007**: HLS playlist generation
- Generate master playlist (lists quality variants)
- Generate media playlists (per quality, lists segments)
- Update playlists as segments are produced
- Use github.com/grafov/m3u8 library
- **Estimate**: 4 hours

**TICKET-008**: Segment serving
- Implement segment endpoint: GET /api/stream/{id}/{quality}/segment_{n}.ts
- Stream segments as they're generated (don't wait for completion)
- Handle missing segments (404 vs 202 "not ready yet")
- CORS headers for browser access
- **Estimate**: 3 hours

### Phase 3: Stream Management (Milestone 3)

**TICKET-009**: Stream lifecycle
- Generate unique stream IDs (UUIDv4)
- Track stream state machine: initializing → active → completed/error
- Implement timeout mechanism (5 min inactivity)
- **Estimate**: 4 hours

**TICKET-010**: Resource cleanup
- Kill FFmpeg processes on stream completion/timeout
- Remove temporary segment files
- Handle cleanup on server shutdown
- Prevent zombie processes
- **Estimate**: 3 hours

**TICKET-011**: Error handling
- Standardized error responses (JSON)
- YouTube extraction failures
- FFmpeg transcode failures
- Rate limiting (basic: max N concurrent streams)
- **Estimate**: 3 hours

**TICKET-011-A**: Request queueing system
- Implement queue data structure (channel-based or priority queue)
- Queue when concurrent streams reach limit (5 streams)
- Track queue position and estimated wait time
- Timeout: remove from queue after 2 minutes of waiting
- API: GET /api/queue/{id}/status returns position
- **Estimate**: 5 hours

**TICKET-011-B**: Analytics system
- SQLite database for aggregate stats
- Track: video URL hash, view count, total duration, timestamp
- Endpoints: GET /api/analytics (summary), GET /api/analytics/popular
- Simple dashboard page showing top videos and stats
- No user tracking, aggregate only
- **Estimate**: 4 hours

**TICKET-011-C**: HTTPS setup
- Generate self-signed certificate for local dev
- Configure Go HTTP server for TLS
- Update docker-compose with cert volumes
- Document production cert setup (Let's Encrypt)
- Redirect HTTP → HTTPS
- **Estimate**: 3 hours

**TICKET-011-D**: Debug caching (dev only)
- Cache yt-dlp results in memory (5 min TTL)
- Only enabled when DEV_MODE=true environment variable set
- Log cache hits for debugging
- Clear documentation that this is dev-only
- **Estimate**: 2 hours

### Phase 4: Web UI (Milestone 4)

**TICKET-012**: Basic web interface
- HTML page with YouTube URL input
- Integrate Video.js player
- Submit form → POST /api/stream → handle queue or play HLS
- Queue status display: position, estimated wait time
- Loading states, error messages
- HTTPS-only (no HTTP fallback)
- **Estimate**: 5 hours

**TICKET-013**: Player enhancements
- Quality selector UI
- Display current quality
- Playback controls (play/pause/seek/volume)
- Mobile-responsive design
- **Estimate**: 3 hours

**TICKET-014**: UX improvements
- Preview thumbnail (if possible)
- Estimated wait time before playback
- Error messages for common issues
- Share/copy stream URL
- **Estimate**: 3 hours

### Phase 5: Observability (Milestone 5)

**TICKET-015**: Metrics & monitoring
- Prometheus metrics endpoint
- Metrics: active_streams, transcode_errors, segments_served, request_duration
- Health check includes FFmpeg/yt-dlp availability
- **Estimate**: 3 hours

**TICKET-016**: Logging & debugging
- Structured logging for all operations
- Log correlation with stream ID
- Debug mode with verbose FFmpeg output
- Log rotation in Docker
- **Estimate**: 2 hours

### Phase 6: Testing & Quality (Milestone 6)

**TICKET-017**: Unit tests
- Test yt-dlp wrapper
- Test stream lifecycle manager
- Test playlist generation
- Test cleanup logic
- Coverage target: 80%+
- **Estimate**: 6 hours

**TICKET-018**: Integration tests
- Test API endpoints end-to-end
- Mock YouTube/FFmpeg where appropriate
- Test concurrent streams
- Test cleanup and timeouts
- **Estimate**: 5 hours

**TICKET-019**: E2E test
- Full flow: submit URL → transcode → serve → playback
- Use real YouTube video (stable test video)
- Verify HLS segments are valid
- **Estimate**: 4 hours

**TICKET-020**: Documentation
- README: setup, usage, architecture
- API documentation (OpenAPI spec)
- ADR directory with all decisions
- Contributing guide
- **Estimate**: 4 hours

### Phase 7: Production Readiness (Milestone 7)

**TICKET-021**: Performance optimization
- Benchmark concurrent stream capacity
- Optimize FFmpeg flags for speed
- Memory profiling and leak detection
- **Estimate**: 5 hours

**TICKET-022**: Security hardening
- Input validation (YouTube URL format)
- Rate limiting per IP
- Resource limits (max stream duration)
- Dependency vulnerability scan
- **Estimate**: 4 hours

**TICKET-023**: Deployment guide
- Production docker-compose.yml
- Reverse proxy setup (nginx/traefik)
- SSL/TLS configuration
- Scaling guide (multiple instances)
- **Estimate**: 3 hours

**TICKET-024**: CI/CD pipeline
- GitHub Actions workflow
- Automated tests on PR
- Docker image build and push
- Linting enforcement
- **Estimate**: 3 hours

## Project Structure

```
TicketShit/
├── Makefile                 # Build automation
├── docker-compose.yml       # Local development environment
├── Dockerfile               # Multi-stage build
├── go.mod                   # Go dependencies
├── go.sum
├── README.md                # Project overview
├── .gitignore
├── .golangci.yml            # Linter configuration
│
├── docs/                    # Documentation
│   ├── adr/                 # Architecture Decision Records
│   │   ├── 001-no-storage.md
│   │   ├── 002-hls-streaming.md
│   │   └── ...
│   ├── api/                 # API documentation
│   │   └── openapi.yml
│   └── architecture.md      # System architecture
│
├── cmd/
│   └── server/              # Main application
│       └── main.go
│
├── internal/                # Private application code
│   ├── api/                 # HTTP handlers
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── routes.go
│   │
│   ├── transcode/           # Transcoding logic
│   │   ├── manager.go       # Orchestrates transcoding
│   │   ├── ffmpeg.go        # FFmpeg wrapper
│   │   ├── ytdlp.go         # yt-dlp wrapper
│   │   └── stream.go        # Stream lifecycle
│   │
│   ├── queue/               # Request queueing
│   │   ├── queue.go         # Queue implementation
│   │   └── manager.go       # Queue management
│   │
│   ├── analytics/           # Analytics storage
│   │   ├── store.go         # SQLite interface
│   │   └── models.go        # Data models
│   │
│   ├── hls/                 # HLS playlist/segment handling
│   │   ├── playlist.go
│   │   ├── segment.go
│   │   └── server.go
│   │
│   ├── cleanup/             # Resource cleanup
│   │   └── cleanup.go
│   │
│   └── metrics/             # Prometheus metrics
│       └── metrics.go
│
├── certs/                   # TLS certificates (dev)
│   ├── server.crt
│   └── server.key
│
├── pkg/                     # Public libraries (if any)
│
├── web/                     # Frontend assets
│   ├── index.html
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js
│
├── scripts/                 # Utility scripts
│   ├── install-deps.sh      # Install FFmpeg, yt-dlp
│   └── test-e2e.sh          # E2E test script
│
└── test/                    # Integration and E2E tests
    ├── integration/
    └── e2e/
```

## Development Workflow

### Initial Setup
```bash
# Clone and setup
git clone <repo>
cd TicketShit
make setup          # Install dependencies
make docker-build   # Build Docker image
```

### Development Loop
```bash
make dev            # Start with hot reload (air/realize)
make test           # Run tests
make lint           # Run linters
```

### Docker Development
```bash
make docker-up      # Start all services
make docker-down    # Stop services
make docker-logs    # View logs
```

### Testing
```bash
make test-unit      # Unit tests only
make test-integration  # Integration tests
make test-e2e       # End-to-end tests
make test-coverage  # Coverage report
```

## Makefile Targets (Proposed)

```makefile
.PHONY: help build test lint docker-build docker-up docker-down

help:              # Show this help
build:             # Build Go binary
test:              # Run all tests
test-unit:         # Unit tests only
test-integration:  # Integration tests
test-e2e:          # E2E tests
test-coverage:     # Generate coverage report
lint:              # Run golangci-lint
fmt:               # Format code
docker-build:      # Build Docker image
docker-up:         # Start docker-compose
docker-down:       # Stop docker-compose
docker-logs:       # Tail logs
clean:             # Remove build artifacts
setup:             # Install dev dependencies
```

## Risk Assessment

### High Risk
1. **YouTube API changes**: yt-dlp may break with YouTube updates
   - Mitigation: Pin yt-dlp version, monitor updates, have fallback strategies

2. **AV1 encoding performance**: AV1 is 2-3x slower than H.264, limiting concurrent streams
   - Mitigation: Use faster encoding preset (cpu-used 6-8), implement queueing, monitor CPU usage, consider H.264 fallback option

3. **FFmpeg resource exhaustion**: Multiple concurrent transcodes could OOM, especially with AV1
   - Mitigation: Limit concurrent streams (start with 3-5), set memory limits, implement queueing

4. **Legal/ToS concerns**: Proxying YouTube content may violate ToS
   - Mitigation: Add disclaimer, ensure no caching, consider educational use only

### Medium Risk
4. **HLS segment timing**: Coordinating segment generation across qualities is complex
   - Mitigation: Start with single quality, add ABR incrementally

5. **Browser compatibility**: HLS support varies across browsers
   - Mitigation: Use Video.js which handles fallbacks

### Low Risk
6. **Docker complexity**: Team may lack Docker experience
   - Mitigation: Provide detailed docs, support both Docker and native development

## Success Metrics

- **Functional**: User can stream any public YouTube video at 128x128 with AV1 encoding
- **Performance**: Playback starts within 10 seconds (allows for AV1 encoding startup)
- **Reliability**: No crashes for 24h with 5 concurrent streams (conservative due to AV1 CPU usage)
- **Quality**: 80%+ test coverage, no linter errors
- **Compression**: Achieve 30%+ smaller file size vs H.264 baseline
- **Scalability**: System can handle 20+ concurrent streams with horizontal scaling (future)

## Next Steps

1. **Review this plan** - Gather feedback, adjust priorities
2. **Setup ADR directory structure** - Start documenting decisions
3. **Create project board** - Convert tickets to GitHub Issues/Jira
4. **Begin Phase 1** - Start with TICKET-001 (scaffolding)
5. **Establish sprint cadence** - Define iteration length and ceremonies

## Immediate Actions (Once Plan Approved)

Upon exiting plan mode, the following project structure will be created:

### Files to Create:
1. **docs/adr/** - Directory for Architecture Decision Records
   - `docs/adr/README.md` - ADR index and template
   - `docs/adr/001-no-storage.md` through `docs/adr/014-max-duration.md` (14 ADRs total)

2. **docs/api/** - API documentation directory
   - `docs/api/README.md` - Placeholder for OpenAPI spec

3. **Root documentation:**
   - `README.md` - Project overview and setup instructions
   - `PLAN.md` - Copy of this comprehensive plan
   - `TODO.md` - Current sprint/session todos (living document, updated each session)
   - `.gitignore` - Standard Go + Docker ignores

4. **Project tracking (living documents):**
   - `docs/user-stories.md` - All user stories organized by epic, with status tracking (TODO/IN_PROGRESS/DONE)
   - `docs/tickets.md` - All tickets organized by phase/milestone, with status and assignee fields
   - Both files will be updated as work progresses across multiple sessions

### Directory Structure to Create:
```
TicketShit/
├── docs/
│   ├── adr/              # Architecture Decision Records
│   │   ├── README.md
│   │   └── 001-no-storage.md through 009-av1-codec.md
│   ├── api/              # API documentation
│   │   └── README.md
│   ├── user-stories.md   # User stories with status tracking
│   ├── tickets.md        # Development tickets with phase organization
│   └── architecture.md   # System architecture diagrams
├── README.md             # Project overview and setup
├── PLAN.md               # This comprehensive plan (snapshot)
├── TODO.md               # Current session/sprint todos (living)
└── .gitignore            # Go + Docker ignores
```

**File Format Details:**
- `docs/user-stories.md`: Each story with ID, epic, description, acceptance criteria, status (TODO/IN_PROGRESS/DONE), notes
- `docs/tickets.md`: Each ticket with ID, title, description, estimate, status, assignee, dependencies, phase/milestone
- `TODO.md`: Session-specific tasks, updated at start/end of each work session
- ADRs: Follow standard ADR template (Context, Decision, Rationale, Consequences, Status)

This structure will be created in the first action after plan approval, before any code is written.

## Resolved Questions

1. ✅ **Playlists**: Single videos only, no playlist support
2. ✅ **HTTPS**: Required for both local dev and production
3. ✅ **Queueing**: Implement queueing system for high load
4. ✅ **Analytics**: Basic analytics (view count, popular videos)
5. ✅ **Caching**: Cache yt-dlp results briefly for debugging purposes only
6. ✅ **Max duration**: 1 hour maximum stream duration
7. ✅ **AV1 preset**: Prioritize speed (cpu-used 8) for better concurrency
8. ✅ **Overload strategy**: Queue requests when at capacity, no H.264 fallback
9. ✅ **AV1 encoder**: Use SVT-AV1 (faster encoding)

---

**Plan Status**: DRAFT - Ready for review and brainstorming session
**Next Review**: After initial feedback from team
