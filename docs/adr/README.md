# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records for the TicketShit project (low-bandwidth YouTube streaming).

## Format

Each ADR follows this template:

```markdown
# ADR-XXX: Title

**Status**: [Proposed | Accepted | Deprecated | Superseded]

## Context
Background and problem statement

## Decision
What we decided to do

## Rationale
Why we made this decision

## Consequences
Positive and negative impacts of this decision

## Implementation Notes
Specific technical details (optional)
```

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [001](./001-no-storage.md) | No Persistent Storage | Accepted |
| [002](./002-hls-streaming.md) | HLS for Streaming Protocol | Accepted |
| [003](./003-go-backend.md) | Go for Backend | Accepted |
| [004](./004-ffmpeg-transcoding.md) | FFmpeg for Transcoding | Accepted |
| [005](./005-yt-dlp-extraction.md) | yt-dlp for YouTube Extraction | Accepted |
| [006](./006-stateless-architecture.md) | Stateless Architecture | Accepted |
| [007](./007-docker-compose.md) | Docker Compose for Local Development | Accepted |
| [008](./008-adaptive-bitrate.md) | Adaptive Bitrate Streaming | Accepted |
| [009](./009-av1-svt.md) | AV1 Video Codec with SVT-AV1 | Accepted |
| [010](./010-https.md) | HTTPS for All Environments | Accepted |
| [011](./011-request-queueing.md) | Request Queueing System | Accepted |
| [012](./012-basic-analytics.md) | Basic Analytics | Accepted |
| [013](./013-debug-caching.md) | Debug-Only Caching | Accepted |
| [014](./014-max-duration.md) | Maximum Stream Duration | Accepted |

## Making Changes

When proposing a new ADR:
1. Copy this template
2. Assign next sequential number
3. Set status to "Proposed"
4. Submit for team review
5. Update status to "Accepted" after approval

When superseding an ADR:
1. Update old ADR status to "Superseded by ADR-XXX"
2. Create new ADR with reference to old one
