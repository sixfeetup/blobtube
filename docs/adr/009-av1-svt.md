# ADR-009: AV1 Video Codec with SVT-AV1

**Status**: SUPERSEDED by ADR-009a (H.264)
**Date**: 2026-02-25
**Superseded Date**: 2026-02-26

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

---

## Why This Decision Was Superseded (2026-02-26)

During implementation and testing, we discovered that **AV1 codec support in HLS is not production-ready** across browsers:

### Issues Encountered

1. **Browser Decode Errors**:
   - Chrome: `MEDIA_ERR_DECODE: video append failed for segment #0`
   - Error: "media used features your browser did not support"
   - Even with proper fMP4 muxing and movflags

2. **fMP4 Muxing Issues**:
   - Segments had malformed track fragment headers
   - Error: "trun track id unknown, no tfhd was found"
   - FFmpeg fMP4 output not compatible with browser MSE for AV1

3. **Limited Browser Support**:
   - Chrome: Spotty AV1 HLS support via MSE (2026)
   - Firefox: Limited/no AV1 HLS support
   - Safari: No AV1 support at all
   - Even with codec string `av01.0.00M.08` in master playlist

4. **Compression Not Critical**:
   - At 50-200 kbps bitrates, H.264 provides adequate compression
   - Quality difference at these low bitrates is minimal
   - Reliability > 30% compression gain

### Conclusion

While AV1 is technically superior for compression, **browser HLS implementation is incomplete** in 2026. H.264 Baseline Profile provides:
- ✅ Universal browser support
- ✅ Well-tested HLS + fMP4 pipeline
- ✅ No muxing or decode issues
- ✅ Still excellent compression at low bitrates

**Recommendation**: Revisit AV1 in 2027-2028 when browser HLS support matures.

See ADR-009a for the new H.264 implementation.
