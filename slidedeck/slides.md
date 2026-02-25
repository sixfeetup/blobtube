---
marp: true
theme: gaia
paginate: true
backgroundColor: #150b4b
color: #f4e8b2
---

<!-- _class: lead -->

![w:400](assets/blobtube-logo.png)

# Announcing BlobTube

A New Web Application by

**The Dark Factory**

February 25, 2026

---

## The Problem

**Bandwidth Scarcity in Key Markets**

- YouTube's standard quality too high for constrained environments
- Emerging markets with limited internet infrastructure
- Educational institutions with bandwidth constraints
- Remote and rural areas with slow connections
- Emergency response scenarios requiring efficient data use

---

## The Solution: BlobTube

**Ultra-Low Bandwidth Video Streaming Proxy**

- Reduces bandwidth requirements by 90%+
- Real-time transcoding with no storage costs
- Privacy-first design with no user tracking
- Open access architecture
- Simple, scalable infrastructure

---

## Technical Innovation

**Next-Generation Compression Technology**

- **AV1 Codec**: 30%+ better compression than H.264
- **Adaptive Streaming**: 64x64 to 256x256 resolution tiers
- **HLS Protocol**: Broad browser support, seamless playback
- **Real-Time Processing**: No caching, no storage overhead

---

## Key Features & Benefits

**Designed for Bandwidth-Constrained Environments**

- Stream videos at ultra-low resolutions (128x128 primary)
- Smart queueing system handles demand spikes
- HTTPS everywhere for security
- No authentication required
- Privacy-first: aggregate analytics only

---

## Use Cases & Market Opportunity

**Serving Underserved Markets**

- Educational institutions with limited bandwidth
- Developing markets with expensive data
- Remote/rural areas with slow internet
- Emergency and disaster response scenarios
- Mobile data conservation for cost-sensitive users

---

## Architecture Highlights

**Simple, Scalable Design**

- Go backend for performance and reliability
- Docker-based deployment for consistency
- Horizontal scaling capability
- No storage infrastructure required
- FFmpeg transcoding pipeline

---

## Competitive Advantages

**Strategic Differentiators**

- No video caching: legal and privacy benefits
- Open source model for transparency
- Superior AV1 compression technology
- Queue-based load management
- Lower operational costs vs. traditional CDNs

---

## Risk Mitigation

**Addressing Key Challenges**

- YouTube ToS: educational use focus, no caching
- AV1 encoding performance: queueing strategy in place
- Monitoring and observability built-in
- Resource limits prevent system exhaustion
- Transparent, documented architecture

---

## Status & Next Steps

**From Planning to Production**

- **Current Phase**: Planning complete, ready for development
- **Development Phases**: 7 phases identified with clear milestones
- **Infrastructure**: Docker environment ready
- **Path to MVP**: Defined with 24 tickets across phases
- **Timeline**: Foundation through production readiness

---

<style scoped>
ul { font-size: 0.75em; }
</style>

## Current Progress: Milestone 1 Complete

**Foundation Phase Delivered**

- **TICKET-001**: Project scaffolding - Go module, directory structure, Makefile, linting
- **TICKET-002**: Docker infrastructure - Multi-stage build, FFmpeg with SVT-AV1, HTTPS support
- **TICKET-003**: Basic HTTP server - TLS-enabled server, health checks, graceful shutdown
- **TICKET-004**: yt-dlp integration - YouTube URL extraction, error handling, unit tests

**Ready for Phase 2: Transcoding Engine**

---

## Questions?

**BlobTube**

Ultra-Low Bandwidth Video Streaming

Built by The Dark Factory

