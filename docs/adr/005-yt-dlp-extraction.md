# ADR-005: yt-dlp for YouTube Extraction

**Status**: Accepted
**Date**: 2026-02-25

## Context

To transcode YouTube videos, we need to extract direct video stream URLs from YouTube. YouTube's structure changes frequently, breaking simple scraping approaches.

## Decision

Use yt-dlp to extract direct YouTube video stream URLs.

## Rationale

1. **Most reliable**: Fork of youtube-dl, actively maintained, handles YouTube's frequent changes
2. **Comprehensive features**: Handles authentication, region restrictions, age-gated content
3. **JSON output**: Machine-readable output for programmatic use
4. **Community support**: Large community, fast updates when YouTube changes
5. **Format selection**: Can select optimal source quality/format

## Consequences

### Positive
- Robust handling of YouTube's complexity
- Active maintenance and quick updates
- Handles edge cases (private videos, region locks, etc.)
- Well-documented JSON output
- Can extract metadata (title, duration, thumbnail)

### Negative
- **External dependency**: Python tool, adds runtime dependency
- **Breakage risk**: May break when YouTube updates, requires yt-dlp updates
- **Process overhead**: Must spawn yt-dlp process for each request
- **Version pinning needed**: Updates may introduce breaking changes

## Implementation Notes

- Install yt-dlp in Docker container
- Use JSON output mode: `yt-dlp -j <url>`
- Extract direct stream URL from JSON response
- Implement error handling for:
  - Invalid URLs
  - Private/deleted videos
  - Region-locked content
  - Rate limiting
- Consider caching results briefly for dev (ADR-013)
- Pin yt-dlp version in production for stability
