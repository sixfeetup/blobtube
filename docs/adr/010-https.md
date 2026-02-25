# ADR-010: HTTPS for All Environments

**Status**: Accepted
**Date**: 2026-02-25

## Context

Modern web features increasingly require HTTPS (Service Workers, Clipboard API, Geolocation). We can use HTTP for local dev and HTTPS for production, or enforce HTTPS everywhere.

## Decision

Require HTTPS for both local development and production environments.

## Rationale

1. **Modern browser requirements**: Many APIs require secure context (HTTPS)
2. **Environment parity**: Dev should match prod to catch issues early
3. **Security by default**: Better security posture
4. **Future-proof**: More features moving to HTTPS-only
5. **Best practice**: Industry trend toward HTTPS-only

## Consequences

### Positive
- Access to all browser APIs
- Dev/prod parity
- Better security posture
- Prepares for future browser requirements

### Negative
- **Slightly complex setup**: Need to generate/manage certificates locally
- **Self-signed cert warnings**: Developers must accept self-signed cert in browser
- **Certificate management**: Production needs Let's Encrypt or similar

## Implementation Notes

### Local Development
- Generate self-signed certificate with script:
  ```bash
  openssl req -x509 -newkey rsa:4096 -nodes \
    -keyout certs/server.key \
    -out certs/server.crt \
    -days 365 \
    -subj "/CN=localhost"
  ```
- Store in `certs/` directory (gitignored)
- Document in README how to trust cert locally
- Go TLS server:
  ```go
  srv := &http.Server{
    Addr: ":8443",
    TLSConfig: &tls.Config{
      MinVersion: tls.VersionTLS12,
    },
  }
  srv.ListenAndServeTLS("certs/server.crt", "certs/server.key")
  ```

### Production
- Use Let's Encrypt for free TLS certificates
- Consider reverse proxy (nginx/traefik) for cert management
- Auto-renewal setup documented in deployment guide
