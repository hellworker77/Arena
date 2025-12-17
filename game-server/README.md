# MMORPG 2D — Step 21: reliable channel over UDP (seq/ack/ackbits)

This step replaces the old plaintext UDP demo with a **binary, versioned, reliable UDP layer**
between **Client <-> Gateway**.

Why: for MMO-like gameplay you need reliability for session control and actions:
- HELLO / attach
- actions (combat, interactions)
- optional critical events
while keeping movement/input **unreliable** for latency.

## What is reliable here
- Client->Gateway: HELLO and ACT are sent on Reliable channel (retransmit until ack).
- Client->Gateway: IN (movement) is Unreliable (no resend).
- Gateway->Client: critical TEXT events are Reliable (demo); replicate streams remain Unreliable.

## Run (2 zones + gateway + client)

Zone 1:
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1 -store ./data1 -http :9101
```

Zone 2:
```bash
go run ./cmd/zone -listen 127.0.0.1:4001 -zone 2 -store ./data2 -http :9102
```

Gateway:
```bash
go run ./cmd/gateway -udp :7777 -zone 1=127.0.0.1:4000 -zone 2=127.0.0.1:4001 -proto 1
```

Client:
```bash
go run ./cmd/client -addr 127.0.0.1:7777 -proto 1 -char 1
```

In client:
- WASD-like: type `m dx dy` (e.g. `m 2 0`)
- action: `a skill targetEID` (e.g. `a 1 1`)
- quit: `q`

Notes:
- The reliable layer is a standard UDP pattern: `seq`, `ack`, `ackBits` (32-bit window).
- This is still a skeleton: no encryption, no congestion control, no compression.
