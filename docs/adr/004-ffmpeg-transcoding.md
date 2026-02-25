# ADR-004: FFmpeg for Transcoding

**Status**: Accepted
**Date**: 2026-02-25

## Context

Video transcoding requires converting YouTube video streams to lower resolutions with AV1 encoding and HLS output format. Multiple options exist: native Go libraries, cloud APIs, or FFmpeg.

## Decision

Use FFmpeg via Go's os/exec for all video transcoding operations.

## Rationale

1. **Industry standard**: FFmpeg is the de facto standard for video processing
2. **Comprehensive codec support**: Supports all codecs including SVT-AV1 (ADR-009)
3. **HLS support**: Native HLS muxer with segment generation
4. **Battle-tested**: Extremely mature, used in production by major streaming platforms
5. **Performance**: Highly optimized, hardware acceleration support
6. **Flexibility**: Powerful CLI allows fine-tuning of encoding parameters

## Consequences

### Positive
- Proven reliability and performance
- Excellent documentation and community support
- One tool handles: decode, resize, encode, and HLS muxing
- Can leverage hardware acceleration if available
- Easy to prototype and test commands independently

### Negative
- **External dependency**: FFmpeg must be installed in Docker containers
- **Process overhead**: Spawning processes has overhead vs native library
- **Error handling**: Must parse stdout/stderr for debugging
- **Version sensitivity**: FFmpeg versions may behave differently

## Implementation Notes

- Use FFmpeg 6.0+ for latest codec support
- Compile with `--enable-libsvtav1` for AV1 encoding
- Spawn via os/exec.CommandContext for timeout/cancellation
- Capture stderr for debugging transcode issues
- Example command:
  ```bash
  ffmpeg -i <input_url> \
    -c:v libsvtav1 -preset 8 -crf 35 \
    -s 128x128 \
    -f hls -hls_time 4 \
    -t 3600 \
    output.m3u8
  ```
