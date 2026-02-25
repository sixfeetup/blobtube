# ADR-009: AV1 Video Codec with SVT-AV1

**Status**: Accepted
**Date**: 2026-02-25

## Context

Video codec choice affects compression ratio, encoding speed, and client compatibility. For low-bandwidth use case, compression is critical, but encoding speed affects concurrent stream capacity.

## Decision

Use AV1 codec with SVT-AV1 encoder, prioritizing speed (preset 8).

## Rationale

1. **Superior compression**: AV1 provides ~30% better compression than H.264
2. **Critical for low-bandwidth**: At 128x128 resolution, every byte saved matters
3. **SVT-AV1 for speed**: SVT-AV1 is 2-3x faster than libaom while maintaining good quality
4. **Preset 8 for concurrency**: cpu-used 8 trades compression for speed, enables more concurrent streams
5. **Modern browser support**: Chrome 90+, Firefox 88+, Edge 90+ support AV1
6. **Royalty-free**: No licensing fees unlike H.265/HEVC

## Consequences

### Positive
- 30%+ better compression than H.264
- Smaller file sizes = better experience on slow connections
- Future-proof codec choice
- SVT-AV1 faster than libaom

### Negative
- **Slower than H.264**: Still 2-3x slower encoding even with SVT-AV1 and preset 8
- **Higher CPU usage**: Limits concurrent streams (target: 5 concurrent)
- **Limited hardware decode**: More battery usage on mobile devices
- **No fallback**: No H.264 fallback, AV1-only (per user decision)

## Implementation Notes

- FFmpeg flags: `-c:v libsvtav1 -preset 8 -crf 35`
  - `-preset 8`: Fast encoding (range 0-13, 8 is fast)
  - `-crf 35`: Constant Rate Factor, balanced quality/size
- Requires FFmpeg compiled with `--enable-libsvtav1`
- Benchmark on target hardware to determine concurrent stream capacity
- Consider preset 10 if need even more speed
- Monitor encoding time metrics to tune preset
