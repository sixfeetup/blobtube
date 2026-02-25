# User Stories

**Project**: BlobTube (Low-Bandwidth YouTube Streaming)
**Last Updated**: 2026-02-25

This document tracks all user stories for the project, organized by epic.

## Status Legend
- **TODO**: Not started
- **IN_PROGRESS**: Currently being worked on
- **DONE**: Completed and verified

---

## Epic 1: Core Video Streaming

### US-001: Basic Video Playback
**As a** user
**I want to** paste a YouTube URL and watch a low-bandwidth version
**So that** I can view videos on slow connections

**Status**: TODO

**Acceptance Criteria**:
- [ ] Web page has input field for YouTube URL
- [ ] Submit button triggers video playback
- [ ] Video plays in 128x128 resolution
- [ ] Playback starts within 5 seconds

**Notes**: Foundation for entire system

---

### US-002: Adaptive Quality
**As a** user
**I want** the video to adapt to my bandwidth
**So that** I get the best quality my connection can handle

**Status**: TODO

**Acceptance Criteria**:
- [ ] Player automatically switches between 64x64, 128x128, 256x256
- [ ] No buffering for more than 2 seconds on quality switch
- [ ] Quality indicator visible in player

**Notes**: Depends on ADR-008 (Adaptive Bitrate Streaming)

---

### US-003: Video Seeking
**As a** user
**I want to** seek through the video
**So that** I can skip to specific parts

**Status**: TODO

**Acceptance Criteria**:
- [ ] Seek bar functional
- [ ] Seeking works without re-transcoding entire video
- [ ] Seeking completes within 3 seconds

**Notes**: HLS segments enable efficient seeking

---

### US-003-A: Queue Status Display
**As a** user
**I want to** see my position in queue when system is busy
**So that** I know when my video will start

**Status**: TODO

**Acceptance Criteria**:
- [ ] Message displays "Position X of Y in queue"
- [ ] Estimated wait time shown
- [ ] Queue position updates in real-time
- [ ] When stream starts, automatically transitions to player

**Notes**: Related to ADR-011 (Request Queueing)

---

## Epic 2: API Access

### US-004: Streaming API
**As a** developer
**I want to** use an API to get video streams
**So that** I can integrate with other applications

**Status**: TODO

**Acceptance Criteria**:
- [ ] POST /api/stream with {"url": "youtube_url"} returns stream ID
- [ ] GET /api/stream/{id}/playlist.m3u8 returns valid HLS playlist
- [ ] API documented with OpenAPI/Swagger spec

**Notes**: RESTful API design

---

### US-005: Stream Status Query
**As a** developer
**I want to** query stream status
**So that** I can handle errors gracefully

**Status**: TODO

**Acceptance Criteria**:
- [ ] GET /api/stream/{id}/status returns state (initializing, ready, error)
- [ ] Error responses include helpful messages
- [ ] 404 for expired/invalid stream IDs

**Notes**: Enables robust client implementations

---

## Epic 3: Operations & Reliability

### US-006: Docker Deployment
**As an** operator
**I want to** deploy via Docker Compose
**So that** I can run the service consistently

**Status**: TODO

**Acceptance Criteria**:
- [ ] `docker-compose up` starts all services
- [ ] Service accessible at https://localhost:8443
- [ ] Logs visible via `docker-compose logs`

**Notes**: ADR-007 (Docker Compose), ADR-010 (HTTPS)

---

### US-007: Monitoring & Metrics
**As an** operator
**I want to** monitor active streams
**So that** I can understand system load

**Status**: TODO

**Acceptance Criteria**:
- [ ] Prometheus metrics endpoint (/metrics)
- [ ] Metrics: active streams, transcode errors, segment generation rate, queue depth
- [ ] CPU/memory usage visible

**Notes**: Essential for production operations

---

### US-008: Resource Cleanup
**As an** operator
**I want** streams to auto-cleanup
**So that** the system doesn't leak resources

**Status**: TODO

**Acceptance Criteria**:
- [ ] Streams inactive for 5 minutes are terminated
- [ ] Streams exceeding 1 hour are terminated (ADR-014)
- [ ] FFmpeg processes killed on cleanup
- [ ] Temp files removed
- [ ] Graceful shutdown on SIGTERM

**Notes**: Critical for production stability

---

### US-008-A: Analytics Dashboard
**As an** operator
**I want to** see basic analytics
**So that** I understand system usage

**Status**: TODO

**Acceptance Criteria**:
- [ ] Dashboard shows: total streams, popular videos, average duration
- [ ] Aggregate data only (no user tracking)
- [ ] Data persisted to SQLite
- [ ] Endpoint: GET /api/analytics

**Notes**: ADR-012 (Basic Analytics), privacy-preserving

---

## Epic 4: Quality & Testing

### US-009: Comprehensive Testing
**As a** developer
**I want** comprehensive tests
**So that** I can refactor confidently

**Status**: TODO

**Acceptance Criteria**:
- [ ] Unit tests for core logic (>80% coverage)
- [ ] Integration tests for API endpoints
- [ ] E2E test: full stream lifecycle
- [ ] Tests run in CI

**Notes**: Foundation for maintainability

---

### US-010: Code Quality Checks
**As a** developer
**I want** code quality checks
**So that** the codebase stays maintainable

**Status**: TODO

**Acceptance Criteria**:
- [ ] golangci-lint passes with strict config
- [ ] Go fmt, go vet enforced
- [ ] Pre-commit hooks available
- [ ] Make target: `make lint`

**Notes**: Prevents technical debt

---

## Summary Statistics

**Total Stories**: 12
**TODO**: 12
**IN_PROGRESS**: 0
**DONE**: 0

**By Epic**:
- Epic 1 (Core Video Streaming): 4 stories
- Epic 2 (API Access): 2 stories
- Epic 3 (Operations & Reliability): 4 stories
- Epic 4 (Quality & Testing): 2 stories

---

## How to Use This Document

1. **During Planning**: Add new stories as requirements emerge
2. **During Sprints**: Update status as work progresses
3. **During Review**: Check acceptance criteria completion
4. **Retrospectives**: Analyze velocity and completion patterns

**Update Frequency**: Update status whenever work begins/completes on a story
