# Object Store V2

A content-addressed, append-only object store with WAL recovery, deduplication, and segment-based compaction.

This project is designed as a production-like internal storage daemon:
- CAS (Content Addressable Storage) with SHA-256 identity
- Key -> hash mapping (latest version semantics)
- Append-only segments for data
- WAL for crash safety
- Background GC / compaction for space reclamation
- REST API with Range GET + ETag semantics
- Prometheus-style metrics
- Kubernetes probes (`/livez`, `/readyz`)
- Graceful shutdown + draining

---

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [On-Disk Layout](#on-disk-layout)
- [Startup and Recovery](#startup-and-recovery)
- [REST API](#rest-api)
    - [PUT](#put-apiv1objectskey)
    - [GET](#get-apiv1objectskey)
    - [Range GET](#range-get-partial-read)
    - [HEAD](#head-apiv1objectskey)
    - [DELETE](#delete-apiv1objectskey)
    - [ETag, If-Match, If-None-Match](#etag-if-match-if-none-match)
- [Health and Probes](#health-and-probes)
- [Metrics](#metrics)
- [Graceful Shutdown](#graceful-shutdown)
- [GC and Compaction](#gc-and-compaction)
- [Operational Playbook](#operational-playbook)
- [Limitations](#limitations)
- [Roadmap](#roadmap)

---

## Architecture Overview

### Data Model

- **CAS object identity**: `hash = SHA256(plaintext)`
- **Key index**: maps `key -> latest(hash)`; deletes are tombstones
- **CAS index**: maps `hash -> (segment_id, offset, size, refcount)`
- **Segments**: immutable append-only files storing objects
- **Manifest**: append-only log describing the lifecycle of segments and index SSTables
- **WAL**: write-ahead log for commit-aware recovery

### Write Path (PUT)
1. Compute SHA-256 of plaintext.
2. If hash exists in CAS: increment refcount, update key index.
3. Otherwise append to active segment.
4. Append WAL commit record.
5. Periodically checkpoint indexes into SSTables.

### Read Path (GET/HEAD)
1. Lookup latest KeyIndex record.
2. Resolve CAS entry by hash.
3. Read object from segment at `(segment_id, offset)`.
4. Range GET reads only the requested byte range.

---

## On-Disk Layout

Repository directory example:

repo/
wal/
00000001.wal
meta/
manifest.log
segments/
seg-00000.seg
seg-00001.seg
index/
key-00001.sst
cas-00001.sst
tmp/
...


### Segment format (conceptual)

Each segment stores a sequence of objects. Objects are immutable.

[SEG1][object_count]
[hash(32)][nonce(12)][size_plain(u64)][size_cipher(u64)][cipher bytes...]
...


---

## Startup and Recovery

Startup uses:
- Manifest to discover segments and SSTables
- WAL to replay only committed operations

Recovery is commit-aware:
- WAL records are applied only after a `Commit` marker
- This prevents partially-written operations from corrupting indexes

---

## REST API

Base prefix: `/api/v1`

### PUT `/api/v1/objects/:key`

Stores the request body as the value for `key`.

Example:
```bash
curl -X PUT "http://localhost:8080/api/v1/objects/foo" --data-binary "hello world"

Response:

    201 Created on success

Notes:

    Deduplication occurs automatically when multiple keys store identical content.

GET /api/v1/objects/:key

Returns the stored object.

curl "http://localhost:8080/api/v1/objects/foo"

Response:

    200 OK with body bytes

    404 Not Found if key does not exist

    410 Gone if key is a tombstone (deleted)

Range GET (partial read)

Supports Range: bytes=start-end.

curl -H "Range: bytes=0-99" -i "http://localhost:8080/api/v1/objects/foo"

Response:

    206 Partial Content

    Content-Range: bytes 0-99/<total>

    Accept-Ranges: bytes

HEAD /api/v1/objects/:key

Returns metadata headers without a body:

    ETag

    Content-Length

    Accept-Ranges

curl -I "http://localhost:8080/api/v1/objects/foo"

DELETE /api/v1/objects/:key

Deletes the key (tombstone). CAS refcount is decremented.

curl -X DELETE -i "http://localhost:8080/api/v1/objects/foo"

Response:

    204 No Content on success

    404 Not Found if key does not exist

ETag, If-Match, If-None-Match

ETag is derived from SHA-256:

    ETag: "sha256:<hex>"

If-None-Match

Use this to avoid downloading unchanged content:

curl -i "http://localhost:8080/api/v1/objects/foo" \
  -H 'If-None-Match: "sha256:..."'

Response:

    304 Not Modified if ETag matches

If-Match

Use this to enforce preconditions:

curl -i "http://localhost:8080/api/v1/objects/foo" \
  -H 'If-Match: "sha256:..."'

Response:

    412 Precondition Failed if ETag does not match

Health and Probes

    /api/v1/livez: liveness probe (always 200 if process is running)

    /api/v1/readyz: readiness probe

        200 when ready

        503 during shutdown/draining

These are intended for Kubernetes.
Metrics

GET /metrics returns Prometheus-compatible text.

Example:

curl http://localhost:8080/metrics

Counters include:

    request totals (PUT/GET/DELETE)

    bytes in/out

    range GET count

    conditional cache responses (304 / 412)

Graceful Shutdown

On SIGTERM / Ctrl+C:

    /readyz begins returning 503

    write requests (PUT/DELETE) are rejected (503)

    existing requests are allowed to finish

    background worker performs a final checkpoint

    the process exits

This behavior supports rolling restarts in orchestration systems.
GC and Compaction
When GC triggers

GC/compaction is threshold-based:

    Runs only when dead ratio or dead bytes exceed configured thresholds.

Default thresholds:

    global dead ratio >= 0.30 OR dead bytes >= 1 GiB

    per segment:

        rewrite if dead ratio >= 0.35

        drop if dead ratio >= 0.95

What GC does

    Mark: compute set of live hashes referenced by latest KeyIndex

    Plan: compute per-segment dead ratio and build actions

    Execute:

        Drop fully-dead sealed segments

        Rewrite partially-dead sealed segments into new sealed segments

        Update manifest in a crash-safe order

Safety rules

    Active segment is never compacted.

    Segment replacement is done via manifest records:

        NewSegment -> write file -> SealSegment -> DropSegment old

Operational Playbook
Disk usage grows continuously

    Ensure background GC is enabled

    Verify thresholds are reachable

    Check that sealed segments exist (active must rotate to create sealed segments)

Crash recovery

    Restart the daemon

    Recovery is manifest + WAL based

    If a segment is missing but was dropped in manifest, it is tolerated

Rolling restart (Kubernetes)

    Readiness goes to 503 on SIGTERM

    Load balancer stops routing new traffic

    Existing requests drain

    Pod terminates cleanly

Limitations

    Streaming PUT without buffering is not yet implemented.

    Full HTTP/2 draining semantics may require additional hyper-level connection control.

    KeyIndex "latest view" iteration must include SSTables for large datasets.

Roadmap

    Reader-based PUT API (true streaming ingest)

    Multipart upload (S3-like)

    Full LSM compaction for KeyIndex/CAS SSTables

    Segment rotation policy (size/object limits)

    Admin endpoints (manual checkpoint, manual compaction)

    Authentication and quotas


---