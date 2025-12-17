# MMORPG 2D — Steps 13 → 20 (sequential, no legacy)

This bundle implements Steps **13–20** in one coherent codebase (so each step builds on the previous cleanly).
Everything remains a **skeleton** (toy data + plaintext UDP), but the architecture follows modern MMO server practice.

## How to run (2 zones + gateway)

Terminal A (zone 1):
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1 -store ./data1 -http :9101
```

Terminal B (zone 2):
```bash
go run ./cmd/zone -listen 127.0.0.1:4001 -zone 2 -store ./data2 -http :9102
```

Terminal C (gateway):
```bash
go run ./cmd/gateway -udp :7777 -zone 1=127.0.0.1:4000 -zone 2=127.0.0.1:4001 -proto 1
```

Client (toy plaintext):
- `HELLO <proto> <charID>`  (proto must match gateway `-proto`)
- `IN <tick> <mx> <my>`
- `ACT <tick> <skill> <targetEID>` (skill=1 is a toy melee hit)

Example:
```bash
echo "HELLO 1 1" | nc -u -w1 127.0.0.1 7777
echo "IN 1 2 0" | nc -u -w1 127.0.0.1 7777
echo "ACT 2 1 1" | nc -u -w1 127.0.0.1 7777
```

## Step index

### Step 13 — **2‑phase zone transfer (prepare/commit/abort)**
- Zone never deletes the player immediately.
- Gateway orchestrates:
  1) `TRANSFER_PREPARE` from old zone
  2) `ATTACH_WITH_STATE` to target zone
  3) upon target `ATTACH_ACK` → `TRANSFER_COMMIT` to old zone
  4) if timeout/failure → `TRANSFER_ABORT` to old zone

### Step 14 — **Ownership + interest layers**
- Entities have `Kind` and `InterestMask`.
- Player has `InterestMask`.
- AOI replication filters by interest layers.

### Step 15 — **Server‑authoritative combat (anti‑cheat)**
- Client sends `ACT` as an *intent*.
- Zone validates cooldown + range and applies damage server-side.
- Results replicated via Event/State (no client-authoritative damage).

### Step 16 — **ECS-ish component split**
- Position, Velocity, Health, Owner, Kind are stored in separate component stores.
- Removes the “god struct entity”.

### Step 17 — **AI + simulation budgets**
- NPCs exist and wander.
- AI budget per tick, and LOD:
  - only NPCs near *any* player get updated.

### Step 18 — **Fault tolerance via snapshots**
- Zone periodically snapshots world + players to disk **asynchronously** (never in the tick).
- On start, zone loads its snapshot if present.

### Step 19 — **Metrics**
- Per-zone `/metrics` (Prometheus text format) and `/debug/vars` (expvar).
- Tick durations, entity counts, replication bytes.

### Step 20 — **Protocol contract**
- UDP has a strict proto version (`HELLO <proto> <charID>`).
- Internal wire has an explicit `WireVersion` constant in code + documented message formats.
