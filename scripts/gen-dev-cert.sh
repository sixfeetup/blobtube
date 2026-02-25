#!/bin/sh
set -eu

OUT_DIR="${1:-./certs}"
mkdir -p "$OUT_DIR"

CRT="$OUT_DIR/server.crt"
KEY="$OUT_DIR/server.key"

if [ -f "$CRT" ] && [ -f "$KEY" ]; then
  echo "certs already exist: $CRT $KEY"
  exit 0
fi

openssl req -x509 -newkey rsa:2048 -sha256 -days 365 -nodes \
  -keyout "$KEY" \
  -out "$CRT" \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "generated: $CRT $KEY"
