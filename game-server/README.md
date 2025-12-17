# MMORPG 2D Step 11 — Persistence without killing the tick (async save queue)

Builds on Step 10 (channels + budget scheduler) and adds **strict persistence**:
- Zone never writes to disk/DB from the tick.
- Dirty character state is queued into a **SaveQueue** worker.
- Periodic autosave + forced save on detach.
- Simple JSON-file store provided as a placeholder (replace with DB).

No legacy. The queue has bounded memory (drops/merges by CharacterID, never blocks the tick).

## Run locally

Terminal A (zone):
```bash
go run ./cmd/zone -listen 127.0.0.1:4000 -zone 1 -aoi 25 -cell 8 -budget 900 -stateEvery 5 -saveEvery 20 -store ./data
```

Terminal B (gateway):
```bash
go run ./cmd/gateway -udp :7777 -zone 127.0.0.1:4000
```

- `-saveEvery 20` means every 20 server ticks (at 20Hz => 1s) we enqueue saves for dirty characters.
- Storage files end up in `./data/char_<id>.json`.
