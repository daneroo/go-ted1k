# Upgrade Go version and dependencies

**Safety review**: NOT VERIFIED — modifies go.mod, Dockerfile, potentially breaks API compatibility

## Context

go.mod specifies Go 1.20 (released Feb 2023). Current stable is Go 1.24.
Dockerfile uses `golang:1.20-alpine` and `alpine:3.17`.
Dependencies are vintage 2020-2021 (pgx v4, nats.go v1.10, etc.).

## Progress

- [ ] Update go.mod to current Go version
- [ ] Update Dockerfile base images
- [ ] Run `go get -u ./...` and `go mod tidy`
- [ ] Fix any breaking API changes
- [ ] Ensure tests still pass after upgrade

## Key dependencies to watch

- `jackc/pgx/v4` → v5 has breaking changes (context-first APIs)
- `nats-io/nats.go` → v1.10 is very old, current is v1.37+
- `mailru/easyjson` — may need regeneration
- `influxdata/influxdb` — may be removable if InfluxDB is no longer used

## Out of scope

- Adding new features
- Changing architecture

## Acceptance criteria

- `go build ./...` and `go test ./...` pass on current Go
- Docker build succeeds with updated base images
- No known security vulnerabilities in dependencies
