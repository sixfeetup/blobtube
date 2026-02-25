# ADR-014: Maximum Stream Duration

**Status**: Accepted
**Date**: 2026-02-25

## Context

Video streams could run indefinitely, consuming resources. Need to balance allowing reasonable content against resource exhaustion.

## Decision

Enforce 1-hour maximum stream duration.

## Rationale

1. **Prevents resource exhaustion**: Limits compute/memory per stream
2. **Covers majority of use cases**: Most YouTube videos under 1 hour
3. **Predictable load**: Operators can calculate max resource usage
4. **Reasonable limit**: Users rarely watch >1 hour at 128x128 resolution
5. **Aligned with use case**: Low-bandwidth streaming typically for shorter content

## Consequences

### Positive
- Bounded resource usage per stream
- Predictable system load
- Prevents abuse (users leaving streams running)
- Simplifies capacity planning

### Negative
- **Long videos truncated**: 2-hour videos cut off at 1 hour
- **User frustration**: May not be obvious until hit limit
- **Arbitrary limit**: Some legitimate use cases excluded

## Implementation Notes

### FFmpeg Flag
```bash
ffmpeg -i <input> -t 3600 ...  # -t limits duration to 3600 seconds (1 hour)
```

### User Messaging
- Show warning before starting stream if video >1 hour
- API response includes:
  ```json
  {
    "stream_id": "abc123",
    "max_duration_seconds": 3600,
    "video_duration_seconds": 7200,
    "truncated": true
  }
  ```
- Player shows warning: "Video limited to first 60 minutes"

### Configuration
- Make configurable via environment variable: `MAX_STREAM_DURATION_SECONDS=3600`
- Default: 3600 (1 hour)
- Allow override for special deployments

### Monitoring
- Track how often videos hit limit
- Consider adjusting limit based on usage patterns
- Prometheus metric: `streams_truncated_total`

### Future Considerations
- Could implement "continuation" feature (restart stream from 1-hour mark)
- Or allow authenticated users longer durations (when auth added)
