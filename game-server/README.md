# MMORPG 2D — Step 22: congestion control + pacing + adaptive RTO

Builds on Step 21 (reliable UDP seq/ack/ackBits) and adds **production-critical guards**:

## Added
- **Per-session token bucket** (bytes/sec + burst) for *all* outbound UDP.
- **Unreliable drop on budget**: replicate lines are dropped when over budget.
- **Reliable backlog cap**: pending reliable bytes are capped; on overflow the session is dropped (strict).
- **Adaptive RTO** using RTT samples from ACKed reliable packets (Jacobson/Karels).
- **Pacing**: retransmits and new reliable sends also respect the token bucket.

## Run
Same as Step 21, plus optional gateway knobs:

Gateway:
```bash
go run ./cmd/gateway -udp :7777 -zone 1=127.0.0.1:4000 -zone 2=127.0.0.1:4001 -proto 1 \
  -rateBps 20000 -burst 40000 -maxReliableBytes 65536
```

Client:
```bash
go run ./cmd/client -addr 127.0.0.1:7777 -proto 1 -char 1
```

Notes:
- This is still a skeleton: no CC like BBR/CUBIC, no encryption, no compression.
- But these guards prevent the classic "reliable UDP melts under loss" failure mode.
