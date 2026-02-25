# ADR-002: HLS for Streaming Protocol

**Status**: Accepted
**Date**: 2026-02-25

## Context

Multiple video streaming protocols exist: Progressive Download (HTTP), HLS, DASH, WebRTC. We need a protocol that works in browsers, supports adaptive bitrate, and enables seeking without re-transcoding.

## Decision

Use HTTP Live Streaming (HLS) with segmented MPEG-TS (.ts) files and M3U8 playlists.

## Rationale

1. **Industry standard**: HLS is the de facto standard for HTTP-based video streaming
2. **Universal browser support**: Native support in Safari, works via Video.js in Chrome/Firefox
3. **Adaptive bitrate**: Supports multiple quality tiers for bandwidth adaptation (ADR-008)
4. **Seeking support**: Segments enable efficient seeking without full re-transcode
5. **Simple server requirements**: Just needs HTTP server, no special streaming server
6. **Battle-tested**: Mature protocol with extensive tooling and documentation

## Consequences

### Positive
- Excellent browser compatibility via Video.js
- Built-in adaptive bitrate streaming
- Efficient seeking and scrubbing
- Well-supported by FFmpeg
- Can start playback before full transcode completes

### Negative
- **More complex than progressive**: Requires playlist generation and segment management
- **Latency**: Typical 10-30 second latency due to segment buffering (acceptable for our use case)
- **Segment coordination**: Must coordinate segment generation across multiple quality tiers

## Implementation Notes

- Use FFmpeg HLS muxer: `-f hls -hls_time 4 -hls_list_size 0`
- Generate master playlist listing quality variants
- Generate media playlist per quality tier
- Segment duration: 4 seconds (balance between latency and overhead)
- Use github.com/grafov/m3u8 Go library for playlist generation
