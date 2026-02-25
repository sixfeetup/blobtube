# ADR-012: Basic Analytics

**Status**: Accepted
**Date**: 2026-02-25

## Context

Understanding usage patterns helps optimize the system. Options range from no analytics to comprehensive tracking. Balance needed between useful metrics and privacy.

## Decision

Track basic aggregate metrics: view count, popular videos, stream duration. No user tracking.

## Rationale

1. **Useful insights**: Understand which videos are popular, typical durations
2. **System optimization**: Identify patterns for capacity planning
3. **Privacy-preserving**: Aggregate only, no user tracking (aligns with ADR-006)
4. **Lightweight**: SQLite sufficient for analytics storage
5. **Optional**: Can be disabled without affecting core functionality

## Consequences

### Positive
- Understand usage patterns
- Inform capacity planning
- Identify popular content types
- Simple implementation (SQLite)
- Privacy-friendly (no user tracking)

### Negative
- **Adds storage dependency**: Need SQLite (lightweight, acceptable)
- **Slight overhead**: Database write per stream
- **Limited insights**: Can't analyze user behavior (intentional)

## Implementation Notes

### Schema
```sql
CREATE TABLE stream_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  video_url_hash TEXT NOT NULL,  -- SHA256 of URL, not URL itself
  duration_seconds INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  quality_tiers TEXT,  -- JSON array: ["64x64", "128x128"]
  completed BOOLEAN
);

CREATE INDEX idx_url_hash ON stream_events(video_url_hash);
CREATE INDEX idx_created_at ON stream_events(created_at);
```

### API Endpoints
- `GET /api/analytics` - Summary stats
  ```json
  {
    "total_streams": 1234,
    "avg_duration_seconds": 420,
    "unique_videos": 567
  }
  ```
- `GET /api/analytics/popular` - Top 10 videos (by hash)
  ```json
  {
    "popular": [
      {"url_hash": "abc123...", "view_count": 45},
      ...
    ]
  }
  ```

### Privacy
- Store SHA256(video_url), not actual URL
- No IP addresses, user agents, or identifying info
- Aggregate metrics only
- Document in README/privacy policy
