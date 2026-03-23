# Verify build and tests

**Safety review**: SAFE — all local, read-only operations

## Context

Repo untouched for ~600 days. Go 1.20 in go.mod and Dockerfile, dependencies are vintage 2020-2021.
Need to confirm everything still compiles and tests pass before changing anything.

## Progress

- [x] `go build ./...` succeeds (Go 1.26.1 local, Go 1.20.14 in Docker)
- [x] `go test ./...` passes — 8 tests across 5 packages (merge, mysql, postgres, timer, util)
- [x] `go vet ./...` clean
- [x] Docker build succeeds (`docker compose build`)
- [x] Note any deprecation warnings or issues
  - `docker-compose.yml`: `version: "3.8"` attribute is obsolete, should be removed
  - Docker image builds with Go 1.20.14 (pinned in Dockerfile), local Go is 1.26.1
- [x] Issues documented as new tasks — none needed, all clean

## Out of scope

- Upgrading Go version or dependencies (separate task)
- Fixing any issues found (separate tasks)

## Acceptance criteria

- All binaries build: subscribe, pump, dumpjsonl, iterator, ipfs, mysqlrestore
- All tests pass
- Docker image builds
- Any issues documented as new tasks
