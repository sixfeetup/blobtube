#!/bin/sh
set -eu

CERT_FILE="${TLS_CERT_FILE:-/certs/server.crt}"
KEY_FILE="${TLS_KEY_FILE:-/certs/server.key}"

i=0
while [ $i -lt 50 ]; do
  if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    break
  fi
  i=$((i+1))
  sleep 0.1
done

if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
  echo "TLS certs not found: $CERT_FILE / $KEY_FILE" >&2
  exit 1
fi

exec /app/blobtube
