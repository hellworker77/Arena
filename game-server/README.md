# MMORPG 2D Step 12 — Zone transfer / handoff (gateway orchestrated)

This step adds **cross-zone transfer** (MMO-style sharding/zone boundaries) on top of Step 11:
- Zone is authoritative and decides when a player must transfer (e.g. crossing a boundary).
- Zone sends `MsgTransfer` to Gateway with the **character state snapshot**.
- Gateway detaches from old zone and attaches to target zone using `MsgAttachWithState`.
- UDP client remains connected to the gateway; only the authoritative zone changes.

No legacy.

## Run locally (2 zones)

Terminal A (zone 1):
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1 -targetZone 2 -xferBoundary 100 -store ./data1
```

Terminal B (zone 2):
```bash
go run ./cmd/zone -listen 127.0.0.1:4001 -zone 2 -targetZone 1 -xferBoundary -100 -store ./data2
```

Terminal C (gateway):
```bash
go run ./cmd/gateway -udp :7777 -zone 1=127.0.0.1:4000 -zone 2=127.0.0.1:4001
```

Client (toy plaintext):
```bash
echo "HELLO 1" | nc -u -w1 127.0.0.1 7777
# push right until crossing boundary (mx=2)
for i in $(seq 1 80); do echo "IN $i 2 0" | nc -u -w1 127.0.0.1 7777; done
```

You should see the gateway logs indicate a transfer, and the client keeps receiving updates.

Notes:
- This skeleton does not implement cryptographic session migration; keep that at gateway.
- Real MMO transfer needs ACK/rollback, session token re-issue, and anti-duplication; TODO markers included.
