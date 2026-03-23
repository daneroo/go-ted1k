# Verify README accuracy

**Safety review**: NOT VERIFIED — some steps involve connecting to remote hosts

## Context

README contains operational details, host references, URLs, data tables, and commands that are ~600 days stale. Need to walk through each section and confirm or correct.

The Backups section data tables are a snapshot from 2023-07-01. Rather than manually rerunning queries, we should create a repeatable script that generates a current report.

## Progress

- [x] Write `scripts/safe-data-report.sh` — prototype bash report
  - Connectivity and year-by-year summary for darwin/mysql and d1-px1/timescaledb
  - Local jsonl file inventory
  - Works but slow (~5 min), to be superseded by `cmd/digest`
- [x] Run the script, review `DATA-REPORT.md` — confirmed data across all sources
  - darwin/mysql: 2007-2026, ~311M rows (note: 2007 data was cleaned from rollup but remains in mysql)
  - d1-px1/timescaledb: 2022-2026, ~103M rows
  - jsonl: 148 files, 15G, 2008-07 to 2020-11 (bit-identical to frozen tar archive, verified with `diff -r`)
- [ ] [cmd/digest](2026-03-23-02b-cmd-digest.md) — Go replacement for safe-data-report.sh with crypto digests
- [ ] Replace README Backups section with link to `DATA-REPORT.md` and `cmd/digest`
- [ ] **Top links**: grafana-ted URL, d1-px1 direct link — are these still live?
- [ ] **Operations section**: host names (galois, d1-px1), SSH paths, docker commands — still accurate?
- [ ] **Development section**: do the listed commands still work?
- [ ] **Setup tips**: NATS, IPFS, Grafana, Postgres sections — still relevant?
- [ ] **JSON section**: easyjson instructions — still valid with current toolchain?
- [ ] **InfluxDB section**: is InfluxDB still in use or should this be removed?

## Out of scope

- Actually fixing infrastructure issues (separate tasks)
- Rewriting the README structure

## Acceptance criteria

- `scripts/safe-data-report.sh` exists, is maintainable, and produces `DATA-REPORT.md`
- README Backups section replaced with link to `DATA-REPORT.md` and the script that generates it
- Every command in the README has been tested or marked as unverified
- Dead links noted
- Sections for unused tech (InfluxDB? IPFS?) flagged for removal
