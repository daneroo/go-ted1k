# Verify data pipeline end-to-end

**Safety review**: NOT VERIFIED — involves connecting to remote databases and potentially writing data

## Context

The core data flow is: MySQL (capture host) → pump → TimescaleDB (d1-px1) → Grafana.
After 600 days, need to confirm each link in this chain is operational.

## Progress

- [ ] **MySQL source**: Is `darwin.imetrical.com:3306/ted` still the capture host? Is it reachable? What's the latest data?
- [ ] **d1-px1 TimescaleDB**: Is the subscribe service running? What's the latest stamp in the watt table?
- [ ] **pump**: Can we pump recent data from MySQL → local TimescaleDB?
- [ ] **pump from d1-px1**: Can we pump from the production TimescaleDB mirror?
- [ ] **Grafana**: Is `grafana-ted.imetrical.com` showing current data?
- [ ] **NATS subscribe**: Is the NATS broker at `nats.ts.imetrical.com:4222` still active? Is subscribe receiving messages?

## Acceptance criteria

- Know the current state of each data source (latest timestamp, row counts)
- pump runs successfully against at least one source
- Gaps in data coverage are identified and documented
