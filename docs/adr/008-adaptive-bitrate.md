# ADR-008: Adaptive Bitrate Streaming

**Status**: Accepted
**Date**: 2026-02-25

## Context

Users have varying bandwidth conditions. We can either provide a single fixed quality or generate multiple quality tiers for adaptive streaming.

## Decision

Generate multiple quality tiers (64x64, 128x128, 256x256) for HLS Adaptive Bitrate (ABR) streaming.

## Rationale

1. **Better user experience**: Player adapts to network conditions automatically
2. **Handles variable bandwidth**: Users on fluctuating connections get optimal quality
3. **Standard HLS feature**: HLS designed for ABR, players expect it
4. **Future-proof**: Supports users upgrading from slower to faster connections
5. **Low-bandwidth focus**: Even 64x64 is usable for understanding content

## Consequences

### Positive
- Excellent UX across different network conditions
- Player handles quality switching automatically
- Supports very slow connections (64x64) and moderate ones (256x256)
- Standard Video.js feature, minimal custom code

### Negative
- **3x transcoding load**: Must transcode same video at 3 resolutions simultaneously
- **More complex coordination**: Must generate master playlist + 3 media playlists
- **Higher CPU usage**: Impacts concurrent stream capacity
- **More segment files**: 3x storage in /tmp during active streams

## Implementation Notes

- Spawn 3 parallel FFmpeg processes per stream:
  - 64x64: `-s 64x64 -b:v 50k`
  - 128x128: `-s 128x128 -b:v 100k`
  - 256x256: `-s 256x256 -b:v 200k`
- Generate master playlist listing all variants
- Each quality generates its own media playlist
- Use Video.js ABR plugin for automatic quality switching
- Consider reducing to 2 qualities (64x64, 128x128) if CPU-bound
