# MMORPG 2D — Step 24: server rewind / lag compensation (combat)

Builds on Step 23 (clock sync + interpolation) and adds **lag compensated combat**.

## Added
- Zone keeps a **position history buffer** per entity (ring buffer by server tick).
- Client sends `actionTick` as an **estimated server tick** (based on sync).
- On `ACT`, the zone validates `actionTick` is within a strict window (default 250ms)
  and then evaluates melee range using **rewound positions** at that tick.
- Damage is applied server-side (still authoritative), only the hit test is rewound.

## Why this is mandatory
Without rewind, the client aims at interpolated (past) positions and the server checks current positions,
so "I hit" becomes "server says miss" under normal MMO latency.

## Run
Same as Step 23.

Client:
```bash
go run ./cmd/client -addr 127.0.0.1:7777 -proto 1 -char 1 -interpMs 150 -tickHz 20
```

Try:
- move until you see NPCs, then `a 1 <npcEID>`
- you should see reliable `EV hit` and `STAT ... hp=...`

Notes:
- This is a simplified rewind (positions only). Extending to projectiles/abilities needs more state history.
