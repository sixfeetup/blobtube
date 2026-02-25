# ADR-006: Stateless Architecture

**Status**: Accepted
**Date**: 2026-02-25

## Context

Web services can store state server-side (sessions, user data) or be stateless. Our requirements include no authentication and ephemeral video processing.

## Decision

Implement a completely stateless architecture with no session state, user tracking, or persistent user data.

## Rationale

1. **Aligns with no-auth requirement**: No users = no user state needed
2. **Simplifies scaling**: Any instance can handle any request, easy horizontal scaling
3. **Privacy by design**: No user tracking or data retention reduces privacy concerns
4. **Operational simplicity**: No state synchronization between instances
5. **Fault tolerance**: Failed instances don't lose critical state

## Consequences

### Positive
- Trivial horizontal scaling (load balance to any instance)
- No session storage/management needed
- Better privacy posture
- Simpler disaster recovery
- No state synchronization complexity

### Negative
- **Can't implement user features**: No user preferences, history, or personalization
- **Limited rate limiting**: Can only rate limit by IP, not user
- **No queue persistence**: Queue state lost on restart (acceptable per ADR-011)

## Implementation Notes

- Use unique stream IDs (UUIDv4) in URLs for tracking active streams
- Stream state stored in-memory only (lost on restart)
- Analytics use aggregate data only (ADR-012)
- Rate limiting implemented via IP address
- For queue state, acceptable to lose on restart (users retry)
