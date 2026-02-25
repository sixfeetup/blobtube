#!/usr/bin/env bash
#
# Script to create GitHub milestones, labels, and issues for BlobTube project
# Usage: ./scripts/create-github-issues.sh
#

set -e

REPO="sixfeetup/blobtube"

echo "Creating GitHub infrastructure for $REPO..."
echo ""

# Check if gh is authenticated
if ! gh auth status &>/dev/null; then
    echo "Error: Not authenticated with GitHub CLI"
    echo "Please run: gh auth login"
    exit 1
fi

echo "=== Creating Milestones ==="

# Create milestones
gh api repos/$REPO/milestones -f title="Phase 1: Foundation" \
    -f description="Basic project infrastructure and scaffolding. Includes Go module setup, Docker configuration, HTTP server, and yt-dlp integration." \
    --jq '.number' > /tmp/milestone-1.txt
echo "✓ Created Phase 1: Foundation ($(cat /tmp/milestone-1.txt))"

gh api repos/$REPO/milestones -f title="Phase 2: Transcoding Engine" \
    -f description="Core FFmpeg transcoding with HLS output. Includes FFmpeg pipeline, multi-quality transcoding, HLS playlists, and segment serving." \
    --jq '.number' > /tmp/milestone-2.txt
echo "✓ Created Phase 2: Transcoding Engine ($(cat /tmp/milestone-2.txt))"

gh api repos/$REPO/milestones -f title="Phase 3: Stream Management" \
    -f description="Stream lifecycle and resource management. Includes state machine, cleanup, error handling, queueing, analytics, and HTTPS setup." \
    --jq '.number' > /tmp/milestone-3.txt
echo "✓ Created Phase 3: Stream Management ($(cat /tmp/milestone-3.txt))"

gh api repos/$REPO/milestones -f title="Phase 4: Web UI" \
    -f description="User-facing web interface. Includes basic web page, Video.js player, quality selector, and UX improvements." \
    --jq '.number' > /tmp/milestone-4.txt
echo "✓ Created Phase 4: Web UI ($(cat /tmp/milestone-4.txt))"

gh api repos/$REPO/milestones -f title="Phase 5: Observability" \
    -f description="Metrics, monitoring, and logging. Includes Prometheus metrics and structured logging." \
    --jq '.number' > /tmp/milestone-5.txt
echo "✓ Created Phase 5: Observability ($(cat /tmp/milestone-5.txt))"

gh api repos/$REPO/milestones -f title="Phase 6: Testing & Quality" \
    -f description="Comprehensive test coverage. Includes unit tests, integration tests, E2E tests, and documentation." \
    --jq '.number' > /tmp/milestone-6.txt
echo "✓ Created Phase 6: Testing & Quality ($(cat /tmp/milestone-6.txt))"

gh api repos/$REPO/milestones -f title="Phase 7: Production Readiness" \
    -f description="Performance, security, and deployment. Includes optimization, security hardening, deployment guide, and CI/CD." \
    --jq '.number' > /tmp/milestone-7.txt
echo "✓ Created Phase 7: Production Readiness ($(cat /tmp/milestone-7.txt))"

echo ""
echo "=== Creating Labels ==="

# Create labels for phases
gh label create "phase-1" --description "Phase 1: Foundation" --color "0e8a16" || echo "Label phase-1 already exists"
gh label create "phase-2" --description "Phase 2: Transcoding Engine" --color "1d76db" || echo "Label phase-2 already exists"
gh label create "phase-3" --description "Phase 3: Stream Management" --color "5319e7" || echo "Label phase-3 already exists"
gh label create "phase-4" --description "Phase 4: Web UI" --color "fbca04" || echo "Label phase-4 already exists"
gh label create "phase-5" --description "Phase 5: Observability" --color "d93f0b" || echo "Label phase-5 already exists"
gh label create "phase-6" --description "Phase 6: Testing & Quality" --color "006b75" || echo "Label phase-6 already exists"
gh label create "phase-7" --description "Phase 7: Production Readiness" --color "c2e0c6" || echo "Label phase-7 already exists"

# Create labels for epics
gh label create "epic-video-streaming" --description "Epic 1: Core Video Streaming" --color "e99695" || echo "Label epic-video-streaming already exists"
gh label create "epic-api" --description "Epic 2: API Access" --color "f9d0c4" || echo "Label epic-api already exists"
gh label create "epic-operations" --description "Epic 3: Operations & Reliability" --color "c5def5" || echo "Label epic-operations already exists"
gh label create "epic-quality" --description "Epic 4: Quality & Testing" --color "bfdadc" || echo "Label epic-quality already exists"

# Create labels for priority
gh label create "priority-high" --description "High priority" --color "d73a4a" || echo "Label priority-high already exists"
gh label create "priority-medium" --description "Medium priority" --color "fbca04" || echo "Label priority-medium already exists"
gh label create "priority-low" --description "Low priority" --color "0e8a16" || echo "Label priority-low already exists"

# Create labels for status
gh label create "status-blocked" --description "Blocked by dependencies" --color "b60205" || echo "Label status-blocked already exists"

echo "✓ Created all labels"

# Store milestone numbers
M1=$(cat /tmp/milestone-1.txt)
M2=$(cat /tmp/milestone-2.txt)
M3=$(cat /tmp/milestone-3.txt)
M4=$(cat /tmp/milestone-4.txt)
M5=$(cat /tmp/milestone-5.txt)
M6=$(cat /tmp/milestone-6.txt)
M7=$(cat /tmp/milestone-7.txt)

echo ""
echo "=== Creating Issues ==="

# Phase 1 Issues
gh issue create --title "TICKET-001: Project Scaffolding" \
    --milestone "$M1" \
    --label "phase-1" \
    --body "$(cat <<'EOF'
## Description
Initialize Go module, create directory structure, setup build tooling.

## Tasks
- [ ] Run `go mod init github.com/sixfeetup/blobtube`
- [ ] Create directory structure (cmd/, internal/, docs/, web/, test/)
- [ ] Setup Makefile with targets: build, test, lint, docker-build, docker-up
- [ ] Add .gitignore for Go + Docker
- [ ] Create basic README.md
- [ ] Setup golangci-lint configuration

## Dependencies
None

## Acceptance Criteria
- `make build` compiles successfully
- Directory structure matches plan
- golangci-lint configured

## Estimate
2 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-001-project-scaffolding)
- [PLAN.md](../blob/main/PLAN.md)
EOF
)"
echo "✓ Created TICKET-001"

gh issue create --title "TICKET-002: Docker Infrastructure" \
    --milestone "$M1" \
    --label "phase-1" \
    --body "$(cat <<'EOF'
## Description
Create Docker setup with FFmpeg (SVT-AV1) and yt-dlp.

## Tasks
- [ ] Create multi-stage Dockerfile (builder + runtime)
- [ ] Install FFmpeg with `--enable-libsvtav1`
- [ ] Verify SVT-AV1: `ffmpeg -encoders | grep svt`
- [ ] Install yt-dlp
- [ ] Create docker-compose.yml with HTTPS support
- [ ] Generate self-signed certificate for dev
- [ ] Setup environment variable configuration
- [ ] Document Docker usage in README

## Dependencies
- TICKET-001

## Acceptance Criteria
- `docker-compose up` starts successfully
- FFmpeg has SVT-AV1 support
- Service accessible at https://localhost:8443

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-002-docker-infrastructure)
- [ADR-007: Docker Compose](../blob/main/docs/adr/007-docker-compose.md)
- [ADR-009: AV1 with SVT-AV1](../blob/main/docs/adr/009-av1-svt.md)
- [ADR-010: HTTPS](../blob/main/docs/adr/010-https.md)
EOF
)"
echo "✓ Created TICKET-002"

gh issue create --title "TICKET-003: Basic HTTP Server" \
    --milestone "$M1" \
    --label "phase-1" \
    --body "$(cat <<'EOF'
## Description
Implement Go HTTP server with TLS, routing, and graceful shutdown.

## Tasks
- [ ] Setup HTTP server with chi or gorilla/mux
- [ ] Implement TLS/HTTPS support
- [ ] Create health check endpoint: GET /health
- [ ] Setup static file serving for web UI
- [ ] Configure structured logging with zerolog
- [ ] Implement graceful shutdown handling (SIGTERM)

## Dependencies
- TICKET-002

## Acceptance Criteria
- Server responds to https://localhost:8443/health
- Graceful shutdown on Ctrl+C
- Logs are structured JSON

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-003-basic-http-server)
- [ADR-003: Go for Backend](../blob/main/docs/adr/003-go-backend.md)
- [ADR-010: HTTPS](../blob/main/docs/adr/010-https.md)
EOF
)"
echo "✓ Created TICKET-003"

gh issue create --title "TICKET-004: yt-dlp Integration" \
    --milestone "$M1" \
    --label "phase-1" \
    --body "$(cat <<'EOF'
## Description
Create wrapper for yt-dlp to extract YouTube video stream URLs.

## Tasks
- [ ] Create ytdlp package in internal/transcode/
- [ ] Implement Execute() function wrapping yt-dlp
- [ ] Parse JSON output (`yt-dlp -j`)
- [ ] Extract direct stream URL
- [ ] Handle errors (invalid URL, private video, region lock)
- [ ] Write unit tests with mocked yt-dlp responses
- [ ] Implement debug caching (DEV_MODE only, ADR-013)

## Dependencies
- TICKET-003

## Acceptance Criteria
- Successfully extracts stream URL from YouTube video
- Handles error cases gracefully
- Unit tests pass

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-004-yt-dlp-integration)
- [ADR-005: yt-dlp for YouTube Extraction](../blob/main/docs/adr/005-yt-dlp-extraction.md)
- [ADR-013: Debug-Only Caching](../blob/main/docs/adr/013-debug-caching.md)
EOF
)"
echo "✓ Created TICKET-004"

# Phase 2 Issues
gh issue create --title "TICKET-005: FFmpeg Pipeline Foundation" \
    --milestone "$M2" \
    --label "phase-2" \
    --body "$(cat <<'EOF'
## Description
Design and implement FFmpeg command for SVT-AV1 HLS transcoding.

## Tasks
- [ ] Design FFmpeg command: `-c:v libsvtav1 -preset 8 -crf 35`
- [ ] Test: YouTube URL → 128x128 AV1 HLS segments
- [ ] Implement exec wrapper in internal/transcode/ffmpeg.go
- [ ] Capture stdout/stderr for debugging
- [ ] Enforce 1-hour max duration: `-t 3600`
- [ ] Write unit tests

## Dependencies
- TICKET-004

## Acceptance Criteria
- FFmpeg successfully transcodes test video
- HLS segments generated in /tmp
- 1-hour limit enforced

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-005-ffmpeg-pipeline-foundation)
- [ADR-004: FFmpeg for Transcoding](../blob/main/docs/adr/004-ffmpeg-transcoding.md)
- [ADR-009: AV1 with SVT-AV1](../blob/main/docs/adr/009-av1-svt.md)
- [ADR-014: Maximum Stream Duration](../blob/main/docs/adr/014-max-duration.md)
EOF
)"
echo "✓ Created TICKET-005"

gh issue create --title "TICKET-006: Multi-Quality Transcoding" \
    --milestone "$M2" \
    --label "phase-2" \
    --body "$(cat <<'EOF'
## Description
Generate 3 quality tiers (64x64, 128x128, 256x256) in parallel.

## Tasks
- [ ] Spawn 3 parallel FFmpeg processes
- [ ] Configure bitrates: 50k, 100k, 200k
- [ ] Coordinate segment generation across qualities
- [ ] Handle individual process failures
- [ ] Test concurrent transcoding

## Dependencies
- TICKET-005

## Acceptance Criteria
- 3 quality tiers generated simultaneously
- Segments sync across qualities
- One quality failing doesn't break others

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-006-multi-quality-transcoding)
- [ADR-008: Adaptive Bitrate Streaming](../blob/main/docs/adr/008-adaptive-bitrate.md)
EOF
)"
echo "✓ Created TICKET-006"

gh issue create --title "TICKET-007: HLS Playlist Generation" \
    --milestone "$M2" \
    --label "phase-2" \
    --body "$(cat <<'EOF'
## Description
Generate master and media playlists for HLS streaming.

## Tasks
- [ ] Integrate github.com/grafov/m3u8 library
- [ ] Generate master playlist listing 3 variants
- [ ] Generate media playlists per quality
- [ ] Update playlists as segments are produced
- [ ] Handle playlist serving via HTTP

## Dependencies
- TICKET-006

## Acceptance Criteria
- Valid HLS master playlist
- Media playlists list segments correctly
- Playlists update dynamically

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-007-hls-playlist-generation)
- [ADR-002: HLS for Streaming Protocol](../blob/main/docs/adr/002-hls-streaming.md)
EOF
)"
echo "✓ Created TICKET-007"

gh issue create --title "TICKET-008: Segment Serving" \
    --milestone "$M2" \
    --label "phase-2" \
    --body "$(cat <<'EOF'
## Description
Implement HTTP endpoints to serve HLS segments and playlists.

## Tasks
- [ ] Endpoint: GET /api/stream/{id}/{quality}/segment_{n}.ts
- [ ] Stream segments as generated (don't wait for completion)
- [ ] Handle missing segments (404 vs 202 "not ready yet")
- [ ] Add CORS headers for browser access
- [ ] Implement byte-range request support

## Dependencies
- TICKET-007

## Acceptance Criteria
- Segments served via HTTP
- Browser can fetch and play segments
- Handles missing segments gracefully

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-008-segment-serving)
- [ADR-002: HLS for Streaming Protocol](../blob/main/docs/adr/002-hls-streaming.md)
EOF
)"
echo "✓ Created TICKET-008"

# Phase 3 Issues
gh issue create --title "TICKET-009: Stream Lifecycle" \
    --milestone "$M3" \
    --label "phase-3" \
    --body "$(cat <<'EOF'
## Description
Implement stream state machine and lifecycle management.

## Tasks
- [ ] Generate unique stream IDs (UUIDv4)
- [ ] Implement state machine: initializing → active → completed/error
- [ ] Track stream metadata (start time, quality, status)
- [ ] Implement 5-minute inactivity timeout
- [ ] Create API: GET /api/stream/{id}/status

## Dependencies
- TICKET-008

## Acceptance Criteria
- Stream transitions through states correctly
- Inactive streams timeout after 5 minutes
- Status API returns accurate state

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-009-stream-lifecycle)
- [ADR-006: Stateless Architecture](../blob/main/docs/adr/006-stateless-architecture.md)
- [User Story US-005](../blob/main/docs/user-stories.md#us-005-stream-status-query)
EOF
)"
echo "✓ Created TICKET-009"

gh issue create --title "TICKET-010: Resource Cleanup" \
    --milestone "$M3" \
    --label "phase-3" \
    --body "$(cat <<'EOF'
## Description
Implement cleanup of FFmpeg processes and temporary files.

## Tasks
- [ ] Kill FFmpeg processes on stream completion/timeout
- [ ] Remove temporary segment files
- [ ] Handle cleanup on server shutdown (SIGTERM)
- [ ] Prevent zombie processes
- [ ] Log cleanup operations

## Dependencies
- TICKET-009

## Acceptance Criteria
- No orphaned FFmpeg processes
- Temp files removed after stream ends
- Graceful cleanup on shutdown

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-010-resource-cleanup)
- [ADR-001: No Persistent Storage](../blob/main/docs/adr/001-no-storage.md)
- [User Story US-008](../blob/main/docs/user-stories.md#us-008-resource-cleanup)
EOF
)"
echo "✓ Created TICKET-010"

gh issue create --title "TICKET-011: Error Handling" \
    --milestone "$M3" \
    --label "phase-3" \
    --body "$(cat <<'EOF'
## Description
Standardized error handling and rate limiting.

## Tasks
- [ ] Standardized JSON error responses
- [ ] Handle YouTube extraction failures
- [ ] Handle FFmpeg transcode failures
- [ ] Basic rate limiting (max N concurrent streams)
- [ ] Error logging and metrics

## Dependencies
- TICKET-010

## Acceptance Criteria
- Consistent error response format
- Rate limiting enforced
- Errors logged with context

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-011-error-handling)
- [User Story US-005](../blob/main/docs/user-stories.md#us-005-stream-status-query)
EOF
)"
echo "✓ Created TICKET-011"

gh issue create --title "TICKET-011-A: Request Queueing System" \
    --milestone "$M3" \
    --label "phase-3,epic-operations" \
    --body "$(cat <<'EOF'
## Description
Implement queue for requests when at capacity (ADR-011).

## Tasks
- [ ] Implement queue data structure (channel-based)
- [ ] Queue when concurrent streams reach limit (5 streams)
- [ ] Track queue position and estimated wait time
- [ ] 2-minute timeout for queued requests
- [ ] API: GET /api/queue/{id}/status
- [ ] Automatic promotion when slot available

## Dependencies
- TICKET-011

## Acceptance Criteria
- Requests queued when at capacity
- Queue position tracked accurately
- Auto-start when slot available
- Timeout removes stale requests

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-011-a-request-queueing-system)
- [ADR-011: Request Queueing System](../blob/main/docs/adr/011-request-queueing.md)
- [User Story US-003-A](../blob/main/docs/user-stories.md#us-003-a-queue-status-display)
EOF
)"
echo "✓ Created TICKET-011-A"

gh issue create --title "TICKET-011-B: Analytics System" \
    --milestone "$M3" \
    --label "phase-3,epic-operations" \
    --body "$(cat <<'EOF'
## Description
SQLite-based analytics for aggregate stats (ADR-012).

## Tasks
- [ ] Create SQLite schema for stream_events
- [ ] Track: video URL hash, view count, duration, timestamp
- [ ] Implement endpoints: GET /api/analytics, GET /api/analytics/popular
- [ ] Create simple dashboard page
- [ ] Ensure no user tracking (aggregate only)
- [ ] Write to DB on stream completion

## Dependencies
- TICKET-010

## Acceptance Criteria
- Analytics data persisted to SQLite
- API returns summary and popular videos
- Privacy-preserving (hashed URLs only)

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-011-b-analytics-system)
- [ADR-012: Basic Analytics](../blob/main/docs/adr/012-basic-analytics.md)
- [User Story US-008-A](../blob/main/docs/user-stories.md#us-008-a-analytics-dashboard)
EOF
)"
echo "✓ Created TICKET-011-B"

gh issue create --title "TICKET-011-C: HTTPS Setup" \
    --milestone "$M3" \
    --label "phase-3" \
    --body "$(cat <<'EOF'
## Description
Configure TLS/HTTPS for development and production (ADR-010).

## Tasks
- [ ] Create script to generate self-signed cert for dev
- [ ] Configure Go HTTP server for TLS
- [ ] Update docker-compose with cert volumes
- [ ] Document production cert setup (Let's Encrypt)
- [ ] Implement HTTP → HTTPS redirect
- [ ] Test with browser (accept self-signed cert)

## Dependencies
- TICKET-002

## Acceptance Criteria
- Dev server runs on https://localhost:8443
- Self-signed cert generation automated
- Production cert setup documented

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-011-c-https-setup)
- [ADR-010: HTTPS for All Environments](../blob/main/docs/adr/010-https.md)
- [User Story US-006](../blob/main/docs/user-stories.md#us-006-docker-deployment)
EOF
)"
echo "✓ Created TICKET-011-C"

gh issue create --title "TICKET-011-D: Debug Caching (Dev Only)" \
    --milestone "$M3" \
    --label "phase-3" \
    --body "$(cat <<'EOF'
## Description
Implement in-memory caching of yt-dlp results for dev (ADR-013).

## Tasks
- [ ] Cache yt-dlp results in memory (5 min TTL)
- [ ] Only enabled when DEV_MODE=true
- [ ] Log cache hits for debugging
- [ ] Document in README that this is dev-only
- [ ] Implement cache expiration

## Dependencies
- TICKET-004

## Acceptance Criteria
- Caching works in dev mode only
- 5-minute TTL enforced
- Cache disabled in production

## Estimate
2 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-011-d-debug-caching-dev-only)
- [ADR-013: Debug-Only Caching](../blob/main/docs/adr/013-debug-caching.md)
EOF
)"
echo "✓ Created TICKET-011-D"

# Phase 4 Issues
gh issue create --title "TICKET-012: Basic Web Interface" \
    --milestone "$M4" \
    --label "phase-4,epic-video-streaming" \
    --body "$(cat <<'EOF'
## Description
Create web UI with Video.js player and queue status.

## Tasks
- [ ] Create HTML page with YouTube URL input
- [ ] Integrate Video.js HLS player
- [ ] Handle form submit → POST /api/stream
- [ ] Display queue status (position, wait time)
- [ ] Show loading states and error messages
- [ ] Ensure HTTPS-only (no HTTP fallback)

## Dependencies
- TICKET-011-A

## Acceptance Criteria
- User can paste URL and watch video
- Queue status displayed when at capacity
- Auto-transition from queue to playback

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-012-basic-web-interface)
- [User Story US-001](../blob/main/docs/user-stories.md#us-001-basic-video-playback)
- [User Story US-003-A](../blob/main/docs/user-stories.md#us-003-a-queue-status-display)
EOF
)"
echo "✓ Created TICKET-012"

gh issue create --title "TICKET-013: Player Enhancements" \
    --milestone "$M4" \
    --label "phase-4,epic-video-streaming" \
    --body "$(cat <<'EOF'
## Description
Enhance video player with quality selector and controls.

## Tasks
- [ ] Add quality selector UI
- [ ] Display current quality indicator
- [ ] Customize playback controls (play/pause/seek/volume)
- [ ] Mobile-responsive design
- [ ] Test on multiple browsers

## Dependencies
- TICKET-012

## Acceptance Criteria
- Quality selector functional
- Mobile-friendly UI
- Works on Chrome, Firefox, Safari

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-013-player-enhancements)
- [User Story US-002](../blob/main/docs/user-stories.md#us-002-adaptive-quality)
- [User Story US-003](../blob/main/docs/user-stories.md#us-003-video-seeking)
EOF
)"
echo "✓ Created TICKET-013"

gh issue create --title "TICKET-014: UX Improvements" \
    --milestone "$M4" \
    --label "phase-4,epic-video-streaming" \
    --body "$(cat <<'EOF'
## Description
Polish user experience with thumbnails, messaging, sharing.

## Tasks
- [ ] Display preview thumbnail (if available from yt-dlp)
- [ ] Show estimated wait time before playback
- [ ] Error messages for common issues
- [ ] Share/copy stream URL functionality
- [ ] Add favicon and page title

## Dependencies
- TICKET-013

## Acceptance Criteria
- Thumbnail shown when available
- Clear error messages
- Shareable stream URLs

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-014-ux-improvements)
EOF
)"
echo "✓ Created TICKET-014"

# Phase 5 Issues
gh issue create --title "TICKET-015: Metrics & Monitoring" \
    --milestone "$M5" \
    --label "phase-5,epic-operations" \
    --body "$(cat <<'EOF'
## Description
Implement Prometheus metrics endpoint.

## Tasks
- [ ] Create /metrics endpoint
- [ ] Metrics: active_streams, transcode_errors, segments_served, request_duration, queue_depth
- [ ] Health check includes FFmpeg/yt-dlp availability
- [ ] Document metrics in README

## Dependencies
- TICKET-011-B

## Acceptance Criteria
- /metrics endpoint returns Prometheus format
- Key metrics tracked
- Health check comprehensive

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-015-metrics--monitoring)
- [User Story US-007](../blob/main/docs/user-stories.md#us-007-monitoring--metrics)
EOF
)"
echo "✓ Created TICKET-015"

gh issue create --title "TICKET-016: Logging & Debugging" \
    --milestone "$M5" \
    --label "phase-5,epic-operations" \
    --body "$(cat <<'EOF'
## Description
Structured logging and debug support.

## Tasks
- [ ] Structured logging for all operations (zerolog)
- [ ] Log correlation with stream ID
- [ ] Debug mode with verbose FFmpeg output
- [ ] Log rotation in Docker
- [ ] Log levels configurable via env var

## Dependencies
- TICKET-003

## Acceptance Criteria
- All logs structured JSON
- Stream ID in all related logs
- Debug mode shows FFmpeg output

## Estimate
2 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-016-logging--debugging)
EOF
)"
echo "✓ Created TICKET-016"

# Phase 6 Issues
gh issue create --title "TICKET-017: Unit Tests" \
    --milestone "$M6" \
    --label "phase-6,epic-quality" \
    --body "$(cat <<'EOF'
## Description
Write unit tests for core components.

## Tasks
- [ ] Test yt-dlp wrapper (mocked)
- [ ] Test stream lifecycle manager
- [ ] Test HLS playlist generation
- [ ] Test cleanup logic
- [ ] Test queue implementation
- [ ] Achieve 80%+ coverage

## Dependencies
None (can start anytime)

## Acceptance Criteria
- 80%+ test coverage
- All tests pass
- Mocking used appropriately

## Estimate
6 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-017-unit-tests)
- [User Story US-009](../blob/main/docs/user-stories.md#us-009-comprehensive-testing)
EOF
)"
echo "✓ Created TICKET-017"

gh issue create --title "TICKET-018: Integration Tests" \
    --milestone "$M6" \
    --label "phase-6,epic-quality" \
    --body "$(cat <<'EOF'
## Description
Write integration tests for API endpoints.

## Tasks
- [ ] Test API endpoints end-to-end
- [ ] Mock YouTube/FFmpeg where appropriate
- [ ] Test concurrent streams
- [ ] Test cleanup and timeouts
- [ ] Test queue behavior under load

## Dependencies
- TICKET-012

## Acceptance Criteria
- API integration tests pass
- Concurrent scenarios tested
- Cleanup verified

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-018-integration-tests)
- [User Story US-009](../blob/main/docs/user-stories.md#us-009-comprehensive-testing)
EOF
)"
echo "✓ Created TICKET-018"

gh issue create --title "TICKET-019: E2E Test" \
    --milestone "$M6" \
    --label "phase-6,epic-quality" \
    --body "$(cat <<'EOF'
## Description
End-to-end test of full stream lifecycle.

## Tasks
- [ ] Full flow: submit URL → transcode → serve → playback
- [ ] Use stable test YouTube video
- [ ] Verify HLS segments are valid
- [ ] Automated E2E test script
- [ ] Run in CI

## Dependencies
- TICKET-018

## Acceptance Criteria
- E2E test passes reliably
- Tests real YouTube video
- Automated in CI

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-019-e2e-test)
- [User Story US-009](../blob/main/docs/user-stories.md#us-009-comprehensive-testing)
EOF
)"
echo "✓ Created TICKET-019"

gh issue create --title "TICKET-020: Documentation" \
    --milestone "$M6" \
    --label "phase-6,epic-quality" \
    --body "$(cat <<'EOF'
## Description
Comprehensive project documentation.

## Tasks
- [ ] Update README with setup, usage, architecture
- [ ] Create OpenAPI spec for API
- [ ] Ensure all ADRs are complete
- [ ] Write contributing guide
- [ ] Document deployment process

## Dependencies
- TICKET-019

## Acceptance Criteria
- README comprehensive
- API documented (OpenAPI)
- Contributing guide available

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-020-documentation)
EOF
)"
echo "✓ Created TICKET-020"

# Phase 7 Issues
gh issue create --title "TICKET-021: Performance Optimization" \
    --milestone "$M7" \
    --label "phase-7" \
    --body "$(cat <<'EOF'
## Description
Benchmark and optimize for production.

## Tasks
- [ ] Benchmark concurrent stream capacity
- [ ] Optimize FFmpeg flags for speed
- [ ] Memory profiling and leak detection
- [ ] CPU profiling
- [ ] Load testing

## Dependencies
- TICKET-019

## Acceptance Criteria
- Know max concurrent streams
- No memory leaks
- Optimized FFmpeg settings

## Estimate
5 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-021-performance-optimization)
EOF
)"
echo "✓ Created TICKET-021"

gh issue create --title "TICKET-022: Security Hardening" \
    --milestone "$M7" \
    --label "phase-7" \
    --body "$(cat <<'EOF'
## Description
Security review and hardening.

## Tasks
- [ ] Input validation (YouTube URL format)
- [ ] Rate limiting per IP
- [ ] Resource limits (max stream duration enforced)
- [ ] Dependency vulnerability scan
- [ ] Security headers (CSP, HSTS, etc.)

## Dependencies
- TICKET-021

## Acceptance Criteria
- Input validation comprehensive
- Rate limiting per IP works
- No critical vulnerabilities

## Estimate
4 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-022-security-hardening)
EOF
)"
echo "✓ Created TICKET-022"

gh issue create --title "TICKET-023: Deployment Guide" \
    --milestone "$M7" \
    --label "phase-7,epic-operations" \
    --body "$(cat <<'EOF'
## Description
Production deployment documentation.

## Tasks
- [ ] Create production docker-compose.yml
- [ ] Document reverse proxy setup (nginx/traefik)
- [ ] SSL/TLS configuration (Let's Encrypt)
- [ ] Scaling guide (multiple instances)
- [ ] Monitoring setup guide

## Dependencies
- TICKET-022

## Acceptance Criteria
- Production deployment documented
- Scaling strategy documented
- Reverse proxy example provided

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-023-deployment-guide)
EOF
)"
echo "✓ Created TICKET-023"

gh issue create --title "TICKET-024: CI/CD Pipeline" \
    --milestone "$M7" \
    --label "phase-7,epic-quality" \
    --body "$(cat <<'EOF'
## Description
Automated CI/CD with GitHub Actions.

## Tasks
- [ ] Create GitHub Actions workflow
- [ ] Automated tests on PR
- [ ] Docker image build and push
- [ ] Linting enforcement
- [ ] Deployment automation (optional)

## Dependencies
- TICKET-023

## Acceptance Criteria
- Tests run on every PR
- Docker images built automatically
- Linting enforced

## Estimate
3 hours

## Related Documentation
- [docs/tickets.md](../blob/main/docs/tickets.md#ticket-024-cicd-pipeline)
EOF
)"
echo "✓ Created TICKET-024"

echo ""
echo "=== Summary ==="
echo "✓ Created 7 milestones"
echo "✓ Created 15 labels"
echo "✓ Created 28 issues"
echo ""
echo "View all issues: https://github.com/$REPO/issues"
echo "View milestones: https://github.com/$REPO/milestones"

# Cleanup temp files
rm -f /tmp/milestone-*.txt

echo ""
echo "Done! 🎉"
