# TODO - Current Session

**Project**: BlobTube (Low-Bandwidth YouTube Streaming)
**Session Date**: 2026-02-26
**Status**: Core Implementation Complete

## Current Session Goals

This session focused on completing the core streaming implementation and resolving critical codec/playback issues.

### Completed ✅

#### Core Implementation
- [x] Implement POST /api/stream endpoint with async orchestration
- [x] Wire up yt-dlp integration for YouTube URL extraction
- [x] Connect FFmpeg transcoding to stream creation flow
- [x] Build multi-quality HLS transcoding (64x64, 128x128, 256x256)
- [x] Create functional web UI with Video.js player
- [x] Add quality selector UI (Auto/Manual quality switching)
- [x] Test end-to-end flow with real YouTube videos

#### Critical Bug Fixes
- [x] Fix YouTube 403 Forbidden errors (yt-dlp → FFmpeg pipe)
- [x] Update yt-dlp to latest version (2026.02.21)
- [x] Fix video playback issue (black screen with audio)
- [x] Fix master playlist URL resolution (relative → absolute URLs)
- [x] Fix race condition (wait for transcoding completion)
- [x] **MAJOR: Switch from AV1 to H.264 codec for browser compatibility**

### Key Lessons Learned 📚

#### 1. AV1 Codec in HLS is Not Browser-Ready
**Problem**: AV1 (libsvtav1) in HLS streams caused MEDIA_ERR_DECODE errors
**Root Cause**:
- AV1 codec support in HLS via Media Source Extensions is incomplete
- Even Chrome (best AV1 support) fails with fMP4 segments
- Firefox and Safari have no/limited AV1 HLS support
- fMP4 muxing issues: "trun track id unknown, no tfhd was found"

**Solution**: Switch to H.264 (libx264) Baseline Profile
- ✅ Universal browser support (Chrome, Firefox, Safari)
- ✅ Well-tested HLS implementation
- ✅ No muxing issues with fMP4 segments
- ✅ Excellent compression at low bitrates

**Lesson**: For production HLS, H.264 is still the safe choice in 2026

**Files Changed**:
- `internal/transcode/ffmpeg.go` - Switch to libx264, add preset mapping
- `internal/api/stream_hls.go` - Update CODECS to avc1.42E01E
- `Dockerfile` - Verify libx264 instead of libsvtav1

#### 2. YouTube Stream URLs Require yt-dlp Piping
**Problem**: Direct FFmpeg access to yt-dlp extracted URLs failed with 403 Forbidden
**Root Cause**: YouTube's stream URLs are IP/client-validated

**Solution**: Pipe yt-dlp stdout → FFmpeg stdin
```go
ytdlpCmd := exec.Command("yt-dlp", "--output", "-", youtubeURL)
ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", ...)
pipe, _ := ytdlpCmd.StdoutPipe()
ffmpegCmd.Stdin = pipe
```

**Lesson**: YouTube URLs can't be passed directly to FFmpeg; use piping

**Files Changed**:
- `internal/transcode/ffmpeg.go` - Add TranscodeHLSFromYtDlpPipe()
- `internal/transcode/multi_quality.go` - Use pipe-based transcoding

#### 3. Async Transcoding Needs State Management
**Problem**: Web UI loaded player before transcoding completed, causing 404 errors

**Root Cause**: Race condition - UI triggered on "active" state (transcoding in progress)

**Solution**: Only load player on "completed" state
```javascript
// Before: if (status.state === 'active' || status.state === 'completed')
// After: if (status.state === 'completed')
```

**Lesson**: Distinguish between "processing" and "ready" states for async operations

**Files Changed**:
- `web/index.html` - Update status polling logic

#### 4. HLS Playlist URLs Must Be Absolute
**Problem**: Quality switching failed with "playlist request error"

**Root Cause**: Relative URLs (e.g., `64x64/index.m3u8`) resolved from wrong base path

**Solution**: Use absolute URLs in master playlist
```
# Before: 64x64/index.m3u8
# After: /api/stream/{id}/64x64/index.m3u8
```

**Lesson**: Always use absolute URLs in HLS master playlists

**Files Changed**:
- `internal/api/stream_hls.go` - Generate absolute URLs

#### 5. yt-dlp Versions Matter
**Problem**: yt-dlp getting 403 Forbidden from YouTube

**Root Cause**: Container had yt-dlp 2023.03.04 (nearly 2 years old)

**Solution**: Download latest yt-dlp directly from GitHub releases
```dockerfile
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
    -o /usr/local/bin/yt-dlp && chmod a+rx /usr/local/bin/yt-dlp
```

**Lesson**: YouTube frequently changes; keep yt-dlp updated

**Files Changed**:
- `Dockerfile` - Download from GitHub instead of apt

### Updated ADRs 📝

#### ADR-009: Video Codec Choice
**Status**: SUPERSEDED

**Original Decision**: AV1 with SVT-AV1 encoder
**Rationale**: Modern codec with excellent compression

**New Decision**: H.264 (AVC) with libx264 encoder
**Rationale**:
- AV1 in HLS has insufficient browser support (2026)
- H.264 Baseline Profile works universally
- Better ecosystem support and debugging tools
- Still provides excellent compression at low bitrates (50-200 kbps)

**Parameters**:
- Encoder: libx264
- Profile: baseline
- Level: 3.0
- Preset: fast (5) to medium (6)
- CRF: 28
- Pixel Format: yuv420p

### In Progress 🔄

None - Core functionality is complete and working!

### Next Session 📋

#### Phase 2: Production Hardening
- [ ] Add request queueing system (ADR-011)
  - Max 5 concurrent streams
  - Queue excess requests with position feedback
  - Timeout handling for queued requests

- [ ] Add basic analytics (ADR-012)
  - SQLite for aggregate stats
  - Track: streams created, completed, errors, avg duration
  - Privacy-preserving (no URLs stored)

- [ ] Add Prometheus metrics endpoint
  - Active streams count
  - Queue length
  - Transcoding duration histogram
  - Error rates

#### Phase 3: Testing & Polish
- [ ] Add integration tests
  - Stream creation flow
  - HLS playlist generation
  - Cleanup lifecycle

- [ ] Add E2E tests
  - Full YouTube → transcode → playback flow
  - Quality switching
  - Error handling

- [ ] Performance optimization
  - Benchmark transcoding speeds
  - Optimize preset selection
  - Test concurrent streams

## Architecture Notes

### Current Stack
- **Backend**: Go 1.22, chi router, zerolog
- **Transcoding**: FFmpeg (libx264), yt-dlp 2026.02.21
- **Streaming**: HLS with fMP4 segments
- **Frontend**: Vanilla JS, Video.js 8.10.0
- **Deployment**: Docker Compose, HTTPS (self-signed dev certs)

### Stream Lifecycle
1. POST /api/stream with YouTube URL
2. Create stream entry (state: initializing)
3. yt-dlp extracts video metadata
4. Set state: active
5. Spawn 3 parallel FFmpeg jobs (64x64, 128x128, 256x256)
6. yt-dlp pipes video → FFmpeg stdin
7. FFmpeg outputs HLS playlists + segments
8. Set state: completed
9. Client polls status, loads player when ready
10. Cleanup after 5 minutes of inactivity

### Quality Tiers
- **64x64**: 50 kbps, ~2.3 MB/min (ultra-low bandwidth)
- **128x128**: 100 kbps, ~4.6 MB/min (very low bandwidth)
- **256x256**: 200 kbps, ~9.2 MB/min (low bandwidth)

All with AAC audio @ 32 kbps

## Blockers

None currently

## Questions for Next Session

1. Should we add VP9 codec support as fallback? (Better compression than H.264, broader support than AV1)
2. Do we need to support higher resolutions? (480p, 720p for better connections)
3. Should queue system be in-memory or Redis-backed?

---

**How to Use This File**:
- Update at start of each session with goals
- Mark tasks complete as you go
- Document lessons learned with root causes and solutions
- Keep ADRs up to date when decisions change
- Review at end of session to plan next steps
