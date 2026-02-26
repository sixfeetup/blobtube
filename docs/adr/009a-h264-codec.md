# ADR-009a: H.264 Video Codec with libx264

**Status**: Accepted
**Date**: 2026-02-26
**Supersedes**: ADR-009 (AV1 with SVT-AV1)

## Context

Initial implementation used AV1 codec for superior compression, but during testing we discovered that AV1 in HLS streams is not production-ready across browsers. Players failed with `MEDIA_ERR_DECODE` errors, even with proper fMP4 muxing. See ADR-009 for full details on why AV1 was superseded.

## Decision

Use H.264 (AVC) codec with libx264 encoder, Baseline Profile, for maximum browser compatibility.

## Rationale

1. **Universal browser support**: H.264 works in Chrome, Firefox, Safari, Edge
2. **Mature HLS implementation**: HLS + H.264 + fMP4 is well-tested across all players
3. **No muxing issues**: libx264 fMP4 output works reliably with browser MSE
4. **Still excellent compression**: At 50-200 kbps bitrates, H.264 is very efficient
5. **Lower CPU usage**: H.264 encoding is faster than AV1, enabling more concurrent streams
6. **Hardware decode support**: Better battery life on mobile devices
7. **Proven reliability**: Decade of production use in HLS streaming

## Trade-offs vs AV1

### What We Lose
- ~30% better compression (AV1 advantage)
- "Modern" codec appeal

### What We Gain
- ✅ **Actually works** in all browsers
- ✅ No decode errors or compatibility issues
- ✅ Faster encoding (better concurrency)
- ✅ Hardware decode support (better mobile battery)
- ✅ Proven, stable ecosystem

**Conclusion**: At these low bitrates (50-200 kbps), reliability and compatibility matter more than compression ratio.

## Implementation Details

### FFmpeg Parameters
```bash
-c:v libx264
-preset fast              # or medium (balance speed/quality)
-crf 28                   # Constant Rate Factor (lower = better quality)
-profile:v baseline       # Maximum compatibility
-level 3.0                # Supports up to 640x480 @ 30fps
-pix_fmt yuv420p          # Standard 4:2:0 chroma subsampling
```

### Preset Mapping
| Number | Preset Name | Use Case |
|--------|-------------|----------|
| 1 | ultrafast | Emergency speed |
| 2 | superfast | Very fast encoding |
| 3 | veryfast | Fast encoding |
| 4 | faster | Faster than default |
| 5 | **fast** | **Default for production** |
| 6 | **medium** | **Balanced (can use if CPU allows)** |
| 7 | slow | Better compression |
| 8 | slower | Much better compression |
| 9 | veryslow | Best compression |

**Default**: preset 5 (fast) for production, preset 6 (medium) if CPU allows.

### Quality Parameters
- **CRF 28**: Good quality at low bitrates
  - CRF range: 0 (lossless) to 51 (worst)
  - 18-28: Visually good quality
  - At 64x64-256x256 resolutions, visual quality differences are minimal
- **Baseline Profile**: No B-frames, simplified entropy coding
  - Best compatibility across all devices
  - Slightly less compression than Main/High profiles
  - Critical for Safari and older devices

### HLS Master Playlist CODECS String
```m3u8
#EXT-X-STREAM-INF:BANDWIDTH=50000,RESOLUTION=64x64,CODECS="avc1.42E01E,mp4a.40.2"
```

- `avc1.42E01E`: H.264 Baseline Profile, Level 3.0
  - `42`: Baseline Profile (0x42)
  - `E0`: Constraint flags
  - `1E`: Level 3.0 (0x1E = 30)
- `mp4a.40.2`: AAC-LC audio

## Consequences

### Positive
- ✅ Works universally across all browsers
- ✅ No decode errors or compatibility issues
- ✅ Faster encoding → more concurrent streams possible
- ✅ Hardware decode → better mobile battery life
- ✅ Mature ecosystem with great debugging tools
- ✅ Still excellent compression at 50-200 kbps

### Negative
- Slightly less compression than AV1 (~30% at high bitrates)
- Less "cutting-edge" technology appeal
- At our low bitrates (50-200 kbps), difference is minimal

### Neutral
- Can revisit AV1 in future when browser support matures (2027-2028)
- Can offer AV1 as opt-in quality tier for Chrome users (future enhancement)

## Testing Results

After switching to H.264:
- ✅ Video plays successfully in Chrome, Firefox, Safari
- ✅ Quality selector works (Auto/64x64/128x128/256x256)
- ✅ No MEDIA_ERR_DECODE errors
- ✅ Smooth playback with HLS adaptive bitrate
- ✅ Transcoding completes in 4-5 seconds for 19-second video

## Benchmarks (Approximate)

**Encoding Speed** (19-second source video, 3 quality tiers in parallel):
- H.264 preset=fast (5): ~4-5 seconds total
- H.264 preset=medium (6): ~6-8 seconds total
- AV1 preset=8: ~10-15 seconds total (estimated, not tested due to compatibility issues)

**File Sizes** (19-second video, 256x256, ~200 kbps):
- H.264: ~450 KB total
- AV1: ~300 KB total (estimated 30% smaller, not tested)
- **Impact**: Extra 150 KB over 19 seconds = ~8 KB/sec = negligible at these resolutions

## Future Considerations

1. **VP9 Codec**: Middle ground between H.264 and AV1
   - Better compression than H.264
   - Better browser support than AV1
   - Consider as future enhancement

2. **Dual Codec Support**: Offer both H.264 and AV1
   - Detect browser capabilities
   - Serve AV1 to Chrome users, H.264 to others
   - Requires dual-encoding (2x CPU cost)

3. **AV1 Revisit Timeline**: 2027-2028
   - Monitor browser HLS + AV1 support maturity
   - Check for fMP4 muxing improvements
   - Test MSE compatibility again

## Related ADRs

- ADR-009: Original AV1 decision (superseded)
- ADR-002: HLS for Streaming Protocol
- ADR-008: Adaptive Bitrate Streaming
- ADR-004: FFmpeg for Transcoding

## References

- [H.264 Baseline Profile Spec](https://en.wikipedia.org/wiki/Advanced_Video_Coding#Profiles)
- [HLS Codec Strings](https://developer.apple.com/documentation/http_live_streaming/hls_authoring_specification_for_apple_devices)
- [libx264 Presets Guide](https://trac.ffmpeg.org/wiki/Encode/H.264)
- [CRF Guide](https://trac.ffmpeg.org/wiki/Encode/H.264#crf)
