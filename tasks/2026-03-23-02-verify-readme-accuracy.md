# Verify README accuracy

**Safety review**: NOT VERIFIED — some steps involve connecting to remote hosts

## Context

README contains operational details, host references, URLs, data tables, and commands that are ~600 days stale. Need to walk through each section and confirm or correct.

## Progress

- [ ] **Top links**: grafana-ted URL, d1-px1 direct link — are these still live?
- [ ] **Backups section**: data tables for mysql/timescale/jsonl — rerun queries, update counts
- [ ] **Operations section**: host names (galois, d1-px1), SSH paths, docker commands — still accurate?
- [ ] **Development section**: do the listed commands still work?
- [ ] **Setup tips**: NATS, IPFS, Grafana, Postgres sections — still relevant?
- [ ] **JSON section**: easyjson instructions — still valid with current toolchain?
- [ ] **InfluxDB section**: is InfluxDB still in use or should this be removed?

## Out of scope

- Actually fixing infrastructure issues (separate tasks)
- Rewriting the README structure

## Acceptance criteria

- Every command in the README has been tested or marked as unverified
- Stale data tables are updated or flagged
- Dead links are noted
- Sections for unused tech (InfluxDB? IPFS?) are flagged for removal
