# ADR-001: No Persistent Storage

**Status**: Accepted
**Date**: 2026-02-25

## Context

Video transcoding and streaming services often cache transcoded videos to reduce computational load when serving the same content to multiple users. However, caching YouTube content raises legal and infrastructural concerns.

## Decision

All transcoding happens in-memory or via temporary streams. No video caching or persistent storage of video content.

## Rationale

1. **Legal compliance**: Avoids gray areas with YouTube's Terms of Service regarding content caching and redistribution
2. **Simplified infrastructure**: No need for large storage volumes, backup systems, or cache invalidation logic
3. **Reduced operational costs**: Storage costs can be significant for video content
4. **Privacy**: No residual data stored that could raise privacy concerns
5. **Stateless design**: Aligns with ADR-006 stateless architecture principle

## Consequences

### Positive
- Simpler system architecture
- Lower storage costs
- Clearer legal position
- Better privacy posture
- Easier compliance with data retention policies

### Negative
- **Higher CPU usage**: Every request requires full transcode, can't serve same video from cache
- **Slower concurrent performance**: Can't share transcoded output across multiple users watching same video
- **Resource intensive**: System must be sized for peak concurrent transcode load

## Implementation Notes

- Use `/tmp` or in-memory filesystem for HLS segments during active streaming
- Implement aggressive cleanup (ADR-014: 5-minute inactivity timeout)
- Consider horizontal scaling to handle concurrent load rather than caching
