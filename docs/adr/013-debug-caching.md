# ADR-013: Debug-Only Caching

**Status**: Accepted
**Date**: 2026-02-25

## Context

During development, repeatedly testing the same YouTube video requires multiple yt-dlp calls, which is slow. Caching could speed iteration, but conflicts with ADR-001 (no storage) for production.

## Decision

Cache yt-dlp results briefly (5 min TTL) for debugging only, not production.

## Rationale

1. **Faster iteration**: Developers can test same video repeatedly without waiting for yt-dlp
2. **Dev-only**: Clear separation from production behavior
3. **Short TTL**: 5 minutes limits staleness
4. **In-memory**: No persistent storage, aligns with ADR-001
5. **Explicitly disabled in prod**: Environment variable gates this feature

## Consequences

### Positive
- Much faster development iteration
- Reduced YouTube API calls during testing
- Simpler debugging (consistent metadata)

### Negative
- **Dev/prod difference**: Behavior differs between environments
- **Potential confusion**: Developers must understand this is dev-only
- **Stale data risk**: Cached results may become stale (5 min helps)

## Implementation Notes

### Configuration
```go
type Config struct {
  DevMode bool  // Set via DEV_MODE=true env var
  CacheTTL time.Duration  // 5 minutes if DevMode
}
```

### Cache Structure
```go
type YtDlpCache struct {
  mu sync.RWMutex
  entries map[string]*CacheEntry
}

type CacheEntry struct {
  Data *VideoMetadata
  ExpiresAt time.Time
}
```

### Usage
```go
if config.DevMode {
  if cached, ok := cache.Get(videoURL); ok {
    log.Debug().Msg("yt-dlp cache hit (dev mode)")
    return cached
  }
}
result := ytdlp.Extract(videoURL)
if config.DevMode {
  cache.Set(videoURL, result, 5*time.Minute)
}
```

### Documentation
- README must clearly state this is dev-only
- Log warning on startup if DevMode enabled
- Default: DevMode=false (production safe)
