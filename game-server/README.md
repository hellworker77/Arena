# MMORPG 2D Step 10 — Replication channels + per-session scheduler (no legacy)

Builds on Step 9:
- Replication is split into **channels**:
  - Movement (frequent, small)
  - State (less frequent)
  - Event (immediate)
- Zone has a strict **per-session budget scheduler** per tick:
  - Spawn/Despawn always highest priority
  - Then movement updates
  - Then state updates (only on `STATE_EVERY_TICKS`)
  - Event queue flushes first (within budget)

Gateway still forwards as plaintext demo messages.

## Run locally

Terminal A (zone):
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1 -aoi 25 -cell 8 -budget 900 -stateEvery 5
```

Terminal B (gateway):
```bash
go run ./cmd/gateway -udp :7777 -zone 127.0.0.1:4000
```

Client -> gateway (UDP plaintext demo):
- `HELLO <charID>`
- `IN <tick> <mx> <my>`

Gateway -> client (UDP plaintext demo):
- Movement:
  - `SPAWN <eid> <x> <y>`
  - `DESPAWN <eid>`
  - `MOV <eid> <x> <y>`
- State (toy):
  - `STAT <eid> hp=<hp>`
- Events (toy):
  - `EV <text>`

This is still a skeleton. Your real crypto handshake stays at the gateway.
