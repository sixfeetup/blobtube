# ADR-007: Docker Compose for Local Development

**Status**: Accepted
**Date**: 2026-02-25

## Context

Development environments need consistency across team members and similarity to production. Options include native development, Docker, or full orchestration (Kubernetes).

## Decision

Use Docker Compose for local development environment.

## Rationale

1. **Environment consistency**: All developers use identical FFmpeg/yt-dlp versions
2. **Matches production**: Production uses Docker, dev environment should match
3. **Simple onboarding**: New developers run `docker-compose up`, no complex setup
4. **Dependency isolation**: FFmpeg, yt-dlp, Go version all containerized
5. **Right complexity level**: Simpler than Kubernetes, more complete than native dev

## Consequences

### Positive
- Consistent environments across team
- Easy onboarding (README: "run docker-compose up")
- Catches deployment issues early
- Can include supporting services (analytics DB, metrics)
- Simplifies CI/CD (same Docker images)

### Negative
- **Learning curve**: Team needs Docker knowledge
- **Slower iteration**: Rebuild container for code changes (mitigated with volume mounts)
- **Resource usage**: Docker adds overhead vs native
- **Debugging**: Slightly harder to attach debuggers

## Implementation Notes

- Use docker-compose.yml with:
  - Go service with volume mounts for hot reload
  - Environment variables for configuration
  - Port mapping for HTTPS (8443)
  - Volume for TLS certificates
- Provide Makefile targets: `make docker-up`, `make docker-down`, `make docker-logs`
- Use multi-stage Dockerfile: builder + runtime
- Support native development as alternative (documented in README)
