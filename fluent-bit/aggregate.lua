-- Aggregate logs by traceId into a single document
-- Buffers log entries and combines messages into multi-line text
-- Flushes when "request completed" message is detected

local buffer = {}

function combine_logs(tag, timestamp, record)
    -- Extract traceId from record
    local traceId = record["traceId"]
    if not traceId or traceId == "" then
        -- If no traceId, pass through as-is
        return 1, timestamp, record
    end

    local message = record["message"] or ""
    local is_complete = string.find(message, "request completed") ~= nil

    -- Initialize buffer entry if it doesn't exist
    if not buffer[traceId] then
        buffer[traceId] = {
            traceId = traceId,
            messages = {},
            method = record["method"],
            path = record["path"],
            status = record["status"],
            latencyMs = record["latencyMs"]
        }
    end

    -- Append message to buffer
    table.insert(buffer[traceId].messages, message)

    -- Update other fields with latest values
    if record["method"] then
        buffer[traceId].method = record["method"]
    end
    if record["path"] then
        buffer[traceId].path = record["path"]
    end
    if record["status"] then
        buffer[traceId].status = record["status"]
    end
    if record["latencyMs"] then
        buffer[traceId].latencyMs = record["latencyMs"]
    end

    -- If this is the completion message, flush the combined entry
    if is_complete then
        local combined = {
            traceId = buffer[traceId].traceId,
            message = table.concat(buffer[traceId].messages, "\n"),
            method = buffer[traceId].method,
            path = buffer[traceId].path,
            status = buffer[traceId].status,
            latencyMs = buffer[traceId].latencyMs
        }

        -- Clean up buffer
        buffer[traceId] = nil

        -- Return the combined record
        return 1, timestamp, combined
    end

    -- Otherwise, suppress this record (don't pass it through yet)
    return -1, timestamp, record
end
