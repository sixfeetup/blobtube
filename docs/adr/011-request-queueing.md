# ADR-011: Request Queueing System

**Status**: Accepted
**Date**: 2026-02-25

## Context

AV1 encoding is CPU-intensive. System can handle limited concurrent streams (target: 5). When capacity reached, we can either reject requests or queue them.

## Decision

Implement queue for stream requests when at capacity. No H.264 fallback.

## Rationale

1. **Better UX**: Waiting in queue better than rejection
2. **Codec consistency**: Maintain AV1-only, no fallback complexity
3. **Load smoothing**: Queue absorbs traffic spikes
4. **Transparent to user**: Queue position shown, automatic start when ready
5. **Simpler than rejection**: Avoids user retry logic

## Consequences

### Positive
- Better user experience during high load
- Automatic handling of capacity limits
- Smooth load distribution
- Maintains codec consistency (AV1-only)

### Negative
- **Wait time**: Users may wait during high load
- **Queue management complexity**: Need timeout handling, position tracking
- **Memory usage**: Queue state held in memory
- **Lost on restart**: In-memory queue lost on server restart (acceptable per ADR-006)

## Implementation Notes

- Max concurrent streams: 5 (conservative for AV1 encoding)
- Queue timeout: 2 minutes (remove from queue if not started)
- Implementation: Go channels or priority queue
  ```go
  type StreamQueue struct {
    requests chan *StreamRequest
    capacity int
  }
  ```
- API: `GET /api/queue/{id}/status` returns:
  ```json
  {
    "position": 3,
    "estimated_wait_seconds": 45,
    "status": "queued"
  }
  ```
- WebSocket or polling for queue position updates
- Automatic transition to playback when stream starts
