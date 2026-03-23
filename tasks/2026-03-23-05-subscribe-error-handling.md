# Improve subscribe error handling

**Safety review**: NOT VERIFIED — modifies production service behavior

## Context

Subscribe currently has basic retry with exponential backoff (1s-10s) on DB insertion failures, but the overall error handling strategy needs review. From the original TODO: "subscribe should exit on error, sleep before restart."

Docker restart policy (`unless-stopped`) provides external restart, but the process should exit cleanly on unrecoverable errors rather than spinning.

## Progress

- [ ] Review subscribe's error handling for NATS disconnection, DB connection loss, and insertion failures
- [ ] Ensure clean exit on unrecoverable errors (let Docker restart handle recovery)
- [ ] Verify graceful shutdown path (signal handling, queue drain with 20s timeout)

## Acceptance criteria

- Subscribe exits with non-zero code on unrecoverable errors
- Transient errors (brief DB hiccup) are retried
- Permanent errors (bad config, auth failure) cause exit
- Graceful shutdown still drains the queue
