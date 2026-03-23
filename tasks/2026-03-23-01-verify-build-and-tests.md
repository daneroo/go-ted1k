# Verify build and tests

**Safety review**: SAFE — all local, read-only operations

## Context

Repo untouched for ~600 days. Go 1.20 in go.mod and Dockerfile, dependencies are vintage 2020-2021.
Need to confirm everything still compiles and tests pass before changing anything.

## Progress

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `go vet ./...` clean
- [ ] Docker build succeeds (`docker compose build`)
- [ ] Note any deprecation warnings or issues
- [ ] Issues documented as new tasks

## Out of scope

- Upgrading Go version or dependencies (separate task)
- Fixing any issues found (separate tasks)

## Acceptance criteria

- All binaries build: subscribe, pump, dumpjsonl, iterator, ipfs, mysqlrestore
- All tests pass
- Docker image builds
- Any issues documented as new tasks
