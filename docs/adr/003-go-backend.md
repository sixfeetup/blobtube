# ADR-003: Go for Backend

**Status**: Accepted
**Date**: 2026-02-25

## Context

Backend language choice affects development velocity, performance, deployment simplicity, and team productivity. Primary requirements: handle concurrent transcoding, serve HTTP, manage lifecycle of external processes (FFmpeg).

## Decision

Use Go (Golang) as the primary backend language for the transcoding service.

## Rationale

1. **Excellent concurrency**: Goroutines and channels make managing multiple concurrent transcodes natural
2. **Strong standard library**: net/http, os/exec, and other stdlib packages cover core needs
3. **Single binary deployment**: Compiles to single statically-linked binary, simplifies Docker images
4. **Performance**: Fast execution, low memory overhead compared to Python/Node.js
5. **Process management**: Strong support for spawning and managing child processes (FFmpeg)
6. **Production-ready**: Used widely for streaming/media services (YouTube, Netflix infrastructure)

## Consequences

### Positive
- Simple deployment (single binary)
- Excellent performance for I/O-bound workloads
- Built-in concurrency primitives ideal for managing multiple streams
- Fast compilation for quick iteration
- Strong typing catches errors at compile time

### Negative
- **Team expertise**: Team needs Go knowledge (learning curve if unfamiliar)
- **Fewer media libraries**: Not as many video processing libraries as Python
- **Verbose error handling**: Can be more verbose than dynamic languages
- **FFmpeg integration**: Must use exec or CGo bindings, no native video processing

## Implementation Notes

- Use Go 1.22+ for improved performance and features
- Use os/exec for FFmpeg process management
- Consider gorilla/mux or chi for HTTP routing
- Use zerolog for structured logging
- Leverage context.Context for timeout and cancellation propagation
