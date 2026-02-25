# TODO - Current Session

**Project**: BlobTube (Low-Bandwidth YouTube Streaming)
**Session Date**: 2026-02-25
**Status**: Planning & Scaffolding

## Current Session Goals

This session focused on comprehensive planning and initial project structure creation.

### Completed ✅
- [x] Define project requirements and scope
- [x] Create 14 Architecture Decision Records (ADRs)
  - ADR-001: No Persistent Storage
  - ADR-002: HLS for Streaming Protocol
  - ADR-003: Go for Backend
  - ADR-004: FFmpeg for Transcoding
  - ADR-005: yt-dlp for YouTube Extraction
  - ADR-006: Stateless Architecture
  - ADR-007: Docker Compose for Local Development
  - ADR-008: Adaptive Bitrate Streaming
  - ADR-009: AV1 Video Codec with SVT-AV1
  - ADR-010: HTTPS for All Environments
  - ADR-011: Request Queueing System
  - ADR-012: Basic Analytics
  - ADR-013: Debug-Only Caching
  - ADR-014: Maximum Stream Duration
- [x] Document 12 user stories across 4 epics
- [x] Create 28 development tickets with estimates
- [x] Set up docs/ directory structure
- [x] Create comprehensive plan document
- [x] Create user stories tracking document
- [x] Create tickets tracking document
- [x] Create this TODO file

### In Progress 🔄
- [ ] Create README.md with quick start guide
- [ ] Create .gitignore for Go + Docker
- [ ] Create PLAN.md (comprehensive plan snapshot)
- [ ] Create docs/api/README.md placeholder

### Next Session 📋
- [ ] Review and refine plan with team
- [ ] Begin TICKET-001: Project Scaffolding
  - Initialize Go module
  - Create directory structure (cmd/, internal/, etc.)
  - Setup Makefile
  - Configure golangci-lint
- [ ] Begin TICKET-002: Docker Infrastructure
  - Create Dockerfile with FFmpeg + SVT-AV1
  - Setup docker-compose.yml
  - Generate dev TLS certificates

## Key Decisions This Session

1. **Codec Choice**: AV1 with SVT-AV1 encoder (preset 8) for speed
2. **No Playlist Support**: Single videos only
3. **HTTPS Everywhere**: Both dev and production
4. **Queueing**: Max 5 concurrent streams, queue excess requests
5. **Analytics**: Basic aggregate stats only (privacy-preserving)
6. **Max Duration**: 1 hour limit per stream
7. **Caching**: Debug-only, 5-minute TTL in dev mode

## Notes

- All planning artifacts saved in docs/ directory
- ADRs provide rationale for all major technical decisions
- Tickets organized into 7 phases with clear dependencies
- Estimated total effort: ~100 hours across all phases
- Ready to begin implementation in next session

## Blockers

None currently

## Questions for Team

1. ~~Should we rename from "TicketShit" to something more professional?~~ **RESOLVED: Renamed to BlobTube**
2. Any additional ADRs needed before starting implementation?
3. Should we set up GitHub project board now or after scaffolding?

---

**How to Use This File**:
- Update at start of each session with goals
- Mark tasks complete as you go
- Add notes, decisions, blockers as they arise
- Review at end of session to plan next steps
- Keep this file as the "source of truth" for current work
