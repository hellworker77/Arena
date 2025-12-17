# MMORPG 2D Step 8 — Gateway + Zone split (no legacy)

This repo is a *strict* skeleton for a 2D MMO architecture split:

- `cmd/gateway`: edge server (UDP sessions) + routes to zones over an internal TCP link
- `cmd/zone`: authoritative simulation host (no knowledge of UDP addresses)

No backward-compat, no legacy protocol parsing.

## Run locally

Terminal A (zone):
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1
```

Terminal B (gateway):
```bash
go run ./cmd/gateway -udp :7777 -zone 127.0.0.1:4000
```

Quick UDP test (sends plaintext demo messages):
- `HELLO <charID>`
- `IN <tick> <mx> <my>`

Example using `nc -u`:
```bash
echo "HELLO 1" | nc -u -w1 127.0.0.1 7777
echo "IN 1 1 0" | nc -u -w1 127.0.0.1 7777
```

This plaintext UDP format is *only* a placeholder to demonstrate the split.
Replace it with your real gateway handshake/crypto.
