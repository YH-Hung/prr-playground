# Fluent Bit → Elasticsearch demo (Go server & client)

This example shows how to log per-request trace IDs, collect logs with Fluent Bit, and ship them to Elasticsearch 8 using Docker Compose.

## What's inside
- Go HTTP server (`server/`) writing JSON logs with `traceId` to `/var/log/app/app.log`.
- Go client (`client/`) sending requests with a fresh `X-Trace-Id` header.
- Fluent Bit tailing the server log and forwarding to Elasticsearch + stdout.
- Elasticsearch 8 single node with security disabled for simplicity.

## Run it (with explanation)
1) Build and start the stack (detached):
```sh
docker compose up --build -d
```
- `--build` forces rebuilding the Go server/client images so code + go.mod changes are picked up.
- `-d` runs containers in the background so you can use the same terminal for the next commands.

2) Generate traffic using the client container:
```sh
docker compose run --rm client \
  -target http://server:8080/hello \
  -count 20 \
  -concurrency 3 \
  -interval 300ms
```
- `docker compose run --rm client` starts a one-off client container then removes it when done, keeping your compose stack clean.
- `-target http://server:8080/hello` tells the client where to send requests; `server` resolves via the compose network to the Go server.
- `-count 20` sends 20 total requests so you get enough samples in Elasticsearch.
- `-concurrency 3` runs three workers in parallel to interleave trace IDs and show grouping across overlapping requests.
- `-interval 300ms` waits 300ms between requests per worker to avoid overwhelming the simple server while still generating multiple overlapping traces.

3) Watch Fluent Bit output (parsed logs):
```sh
docker compose logs -f fluent-bit
```
- `logs -f` tails the Fluent Bit container logs so you can see parsed JSON records and delivery status to Elasticsearch in real time.

4) Query Elasticsearch for recent logs (get 5 newest):
```sh
curl -s "http://localhost:9200/requests-*/_search" \
  -H 'Content-Type: application/json' \
  -d '{"size":5,"sort":[{"@timestamp":{"order":"desc"}}]}'
```
- `requests-*` matches the daily index name produced by Fluent Bit (`requests-%Y-%m-%d`).
- `size:5` limits to 5 docs to keep the output readable.
- `sort @timestamp desc` shows the most recent ingested entries first.

Filter by a specific trace (replace `<trace>` with a traceId from the logs):
```sh
curl -s "http://localhost:9200/requests-*/_search" \
  -H 'Content-Type: application/json' \
  -d '{"query":{"term":{"traceId.keyword":"<trace>"}}}'
```
- Each `traceId` should return exactly **one document** with all log lines combined.

5) Tear down and clean volumes:
```sh
docker compose down -v
```
- `down` stops and removes containers.
- `-v` also removes the named volume holding the server log, ensuring a clean slate for the next run.

## How Log Aggregation Works

This implementation combines multiple log lines from the same request (same `traceId`) into a **single Elasticsearch document** with a multi-line message field. Here's how it works:

### Pipeline Flow

```
Server Log File → Fluent Bit Tail Input → JSON Parser → Lua Aggregator → Elasticsearch
```

### Step-by-Step Process

1. **Input** (`[INPUT]` section in `fluent-bit.conf`):
   - Fluent Bit tails `/var/log/app/app.log` line by line
   - Each line is tagged as `app` and passed to the filter chain

2. **JSON Parsing** (`[FILTER]` parser):
   - Parses each log line as JSON
   - Extracts fields like `traceId`, `message`, `method`, `path`, `status`, `latencyMs`
   - The parsed record continues to the next filter

3. **Lua Aggregation** (`[FILTER]` lua + `aggregate.lua`):
   - **Buffering**: Maintains an in-memory buffer keyed by `traceId`
   - **For each incoming log entry**:
     - If `traceId` exists in buffer: appends the `message` to the buffer entry's message array
     - If new `traceId`: creates a new buffer entry
     - Updates other fields (method, path, status, latencyMs) with latest values
   - **Suppression**: Returns `-1` to drop intermediate entries (they're buffered, not sent yet)
   - **Flush trigger**: When a log entry contains `"request completed"`:
     - Combines all buffered messages with `\n` separator: `message1\nmessage2\n...`
     - Creates a single combined record with merged fields
     - Returns `1` to pass the combined record to output
     - Removes the entry from buffer

4. **Output** (`[OUTPUT]` es):
   - Only the combined records (one per `traceId`) are sent to Elasticsearch
   - Each document contains all log lines from that request in a single `message` field

### Example Transformation

**Before aggregation** (multiple log lines):
```json
{"traceId": "abc", "message": "handler finished", "method": "GET", "path": "/hello", "status": 200, "latencyMs": 0}
{"traceId": "abc", "message": "request completed", "method": "GET", "path": "/hello", "status": 200, "latencyMs": 52}
```

**After aggregation** (single Elasticsearch document):
```json
{
  "traceId": "abc",
  "message": "handler finished\nrequest completed",
  "method": "GET",
  "path": "/hello",
  "status": 200,
  "latencyMs": 52
}
```

### Key Implementation Details

- **Buffer storage**: The Lua script uses a global `buffer` table to store entries by `traceId`
- **Return codes**: 
  - `-1`: Drop/suppress the record (intermediate entries)
  - `1`: Pass the record through (combined entry or entries without traceId)
- **Field merging**: Latest values for `status` and `latencyMs` are kept (from the "request completed" entry)
- **Message combination**: All messages are joined with `\n` to create a multi-line text field
- **Completion detection**: Uses `string.find()` to detect the "request completed" message pattern

### Configuration Files

- **`fluent-bit/fluent-bit.conf`**: Defines the filter pipeline with Lua filter
- **`fluent-bit/aggregate.lua`**: Contains the aggregation logic
- Both files are mounted into the Fluent Bit container via `docker-compose.yml`

## Notes
- Server log file is volume-mounted so Fluent Bit and Elasticsearch see the same data.
- Elasticsearch security is disabled here for brevity; enable auth in real deployments.
- The Lua buffer is in-memory only; if Fluent Bit restarts, buffered entries may be lost (acceptable for this use case).

