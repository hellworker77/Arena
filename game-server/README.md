# MMORPG 2D — Step 23: clock sync + interpolation buffer (client-side)

Builds on Step 22 (reliable UDP + congestion control) and adds:

## Added
- Gateway tags every replication line with `serverTick`.
- Client maintains a **serverTick ↔ local time** mapping (simple sync).
- Client keeps a per-entity **sample buffer** and renders **interpolated** positions
  at `now - interpDelay` (default 150ms) to smooth jitter.

This is intentionally minimal:
- No fancy drift correction (PLL), no packet timestamp echo, no client-side prediction.
- But it demonstrates the correct MMO pattern: **buffer + interpolate**, not "render immediately".

## Run
Same as Step 22.

Client flags:
```bash
go run ./cmd/client -addr 127.0.0.1:7777 -proto 1 -char 1 -interpMs 150 -tickHz 20
```

In client:
- `m dx dy` movement (unreliable)
- `a skill targetEID` action (reliable)
- `q` quit

Output:
- Raw lines are still printed (EV/STAT etc).
- Interpolated positions print periodically as `RENDER tick=<t> eid=<id> x=<..> y=<..>`.
