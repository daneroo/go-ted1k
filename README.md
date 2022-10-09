# Ted1k with Go - data pump

- 2022-10-09 Moved this to gateway 
  - `sudo snap install go --classic` on gateway ubuntu host
    - until I containerize pump and subscribe.
  - also copied over `data/grafana/grafana.db` (`grafana-2022-10-09.db`)

## TODO

- Bring back Evernote TODO to here...
- [Separate e2e tests](https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go/25970712)
- subscribe: reconnect on conn error(s)
  - <https://docs.docker.com/compose/compose-file/compose-file-v2/#healthcheck>
- [docker-compose composition](https://docs.docker.com/compose/extends/)
- Verify(eph,eph) - has a timing bug for output Verified before .. took..
- Consider byDay as final format
  - Document file layout and format (IPFS/jsonl) including file names and aggregation directories
  - Write to ipfs byDay:~650k/s vs byMonth:1.0M/s (same to a lesser extent with json) 930k/s vs 950k/s
- channels of slices `chan []types.Entry`
  - Extract slice manipulation
- Verify and **merge** - for faster inserts
- ipfs
  - Powergate / Filecoin - <https://blog.textile.io/hosted-powergate/>
- flux (at least write) <https://github.com/influxdata/influxdb-client-go#writes>
- off-by-one error in gaps? add tests
  - progress.Gaps: 2020-05-15T23:59:59Z 2020-05-17T00:00:00Z : 24h0m1s
- e2e testing (with ephemeral source) in docker-compose
- progress on nats - with client
- Gather performance/integrity and history in markdown (PERFORMANCE.md)
- [See Evernote](https://www.evernote.com/shard/s60/nl/1773032759/ae1b9921-7e85-4b75-a21b-86be7d524295/)

## Operations

```bash
# Start backing services
docker-compose up -d

# restore a mysqldump snapshot to mysql (see MYSQL.env)
./restore-db.sh

# unit tests
go test -v ./...

# run subscribe ( in screen )
time go run cmd/subscribe/subscribe.go
# run the pump
time go run cmd/pump/pump.go
```

## History

Given all the history of implementations and hardware adventures and failures,

- Scraping of the footprints software (Windows) (Aria Windows Media Server - 2008)
- Remote mounting the footprints data to linux, and doing ETL from SQLite -> MYSQL
- Implementing the python native scraper on linux
- Disk failures on goedel and later cantor
- Porting the scraper to Raspberry Pi (Which had a catastrophic Flash Card Failure) - that lasted a month 8-(
- Moving back to euler (née cantor) linux server (We lost data from 2016-02-14 to 2016-03-12)
- Re-implementing in Go (2018)
- Implementing the pump in Go for persistence neutral backups

In December 2020 we consolidated the data pump (this repo) to transport and synchronize snapshots, agnostic to storage and transport. Accommodating MySQL,Postgres, `.jsonl` files, IPFS

The aggregation into a final rollup of all legacy snapshots was performed based on backups up to 2020-11-20.
The detailed aggregation logs and process are found in the [RESTORE-log.md](RESTORE-log.md) file

The final rollup archive: `ted.20201120.2332Z.rollup-clean.jsonl.tar.bz2`, with IPFS/CID:`QmSLJPEZocdPZ99pazEkiJTaf3B1zeBmAQWEr7n9fSNgEu`

## Setup tips

### nats

- starts a nats image with websocket support
- pass in a websocket enabled configuration
- 4222 is for clients.
- 6222 is a routing port for clustering.
- 8222 is an HTTP management port for information reporting.
- 9222 is for websocket clients

```bash
docker run -it --rm -p 4222:4222 -p 8222:8222 -p 9222:9222 -v $(pwd)/data/nats/nats-server.conf:/nats/conf/nats-server.conf --name nats synadia/nats-server:nightly
```

### IPFS

```bash
# unpin all recursive pins - and run the gc
ipfs pin ls --type recursive | cut -d' ' -f1 | xargs -n1 ipfs pin rm
ipfs repo gc
```

### Postgres/TimescaleDB

```sql
CREATE DATABASE ted;

CREATE TABLE IF NOT EXISTS watt (
  -- timestamp [ (p) ] with time zone
  stamp TIMESTAMP WITHOUT TIME ZONE NOT NULL PRIMARY KEY,
  watt integer NOT NULL DEFAULT '0'
);
-- For timescaledb (must be done on an empty table)
SELECT create_hypertable('watt', 'stamp')
```

## JSON

### Decoding / UnMarshaling

For decoding which, is a bottleneck, we looked at many streaming modules (ffjson/fastjson, etc), not many of which can properly handle our json per line format well, so we stuck with the `encoding/json` implementation

```
for (*json.Decoder).More() {
  err := dec.Decode(&entry)
}
```

We did get a slight improvement (`540k/s -> 610k/s`) from `easyjson` by generating a `json.Unmarshaler` interface:

```
go get -u github.com/mailru/easyjson/...
${GOPATH-~/go}/bin/easyjson types/types.go
```

For json encoding, easyjson actually slightly worsened performance (`500k/s -> 495k/s`), but we are not using it because`fmt.Fprintf()` is faster than `json.Encoder.Encode()` with or without `easyjson`, yielding `500k/s -> 875k/s`

## InfluxDB

```bash
docker exec -it go-ted1k_tedflux_1 bash
influx -database ted -execute 'select count(value) from watt'
select mean(value)*24/1000 from watt where time > '2008-01-01' and time < '2016-01-01' group by time(7d)
```

### Downsampled time series

Truncate for D,M,Y: <http://play.golang.org/p/PUNNHq9sh6>

Continuous Queries are not appropriate for historical data loading.
I should implement my own select .. into (in go), using table names as in mysql

```InfluxQL
select mean(value)*24/1000 into kwh_1d from watt where time > '2015-09-01' group by time(1d)
```

## Performance - Pump

This was performed a an ubuntu:20.04 VM (Proxmox), on a mac mini 2012/8G/2TB-SSD, databases/ipfs running in docker in the same VM. The `ephemeral` data set is a synthetic 31M data points representing~1year of second data; 100MB/month, 1.35GB total

```bash
$ time go run cmd/pump/pump.go
2020-12-10T17:52:39.982Z - Starting TED1K pump
2020-12-10T17:52:39.991Z - Connected to MySQL
2020-12-10T17:52:40.122Z - Connected to Postgres

2020-12-10T17:52:40.163Z - -=- ephemeral -> ephemeral
2020-12-10T17:52:41.288Z - ephemeral -> ephemeral took 1.125s, rate ~ 27.9M/s count: 31415926
2020-12-10T17:52:41.288Z - -=- ephemeral <-> ephemeral
2020-12-10T17:52:44.788Z - ephemeral <-> ephemeral took 3.499s, rate ~ 9.0M/s count: 31415926
2020-12-10T17:52:44.788Z - Verified ephemeral <-> ephemeral:
2020-12-10T17:52:44.788Z - [2020-01-01T00:00:00Z, 2020-12-29T14:38:45Z](31415926) Equal

2020-12-10T17:52:44.788Z - -=- ephemeral -> jsonl
2020-12-10T17:53:18.584Z - ephemeral -> jsonl took 33.796s, rate ~ 929.6k/s count: 31415926
2020-12-10T17:53:18.966Z - -=- jsonl -> ephemeral
2020-12-10T17:54:11.536Z - jsonl -> ephemeral took 52.57s, rate ~ 597.6k/s count: 31415926
2020-12-10T17:54:11.536Z - -=- ephemeral <-> jsonl
2020-12-10T17:55:05.434Z - ephemeral <-> jsonl took 53.889s, rate ~ 583.0k/s count: 31415926
2020-12-10T17:55:05.434Z - Verified ephemeral <-> jsonl:
2020-12-10T17:55:05.434Z - [2020-01-01T00:00:00Z, 2020-12-29T14:38:45Z](31415926) Equal

2020-12-10T17:55:05.443Z - -=- ephemeral -> ipfs
2020-12-10T17:55:34.981Z - ephemeral -> ipfs took 29.537s, rate ~ 1.1M/s count: 31415926
2020-12-10T17:55:35.071Z - -=- ipfs -> ephemeral
2020-12-10T17:56:31.769Z - ipfs -> ephemeral took 56.698s, rate ~ 554.1k/s count: 31415926
2020-12-10T17:56:31.769Z - -=- ephemeral <-> ipfs
2020-12-10T17:57:30.256Z - ephemeral <-> ipfs took 58.486s, rate ~ 537.2k/s count: 31415926
2020-12-10T17:57:30.256Z - Verified ephemeral <-> ipfs:
2020-12-10T17:57:30.256Z - [2020-01-01T00:00:00Z, 2020-12-29T14:38:45Z](31415926) Equal

2020-12-10T17:57:30.256Z - -=- ephemeral -> postgres
2020-12-10T17:59:28.548Z - ephemeral -> postgres took 1m58.291s, rate ~ 265.6k/s count: 31415926
2020-12-10T17:59:28.567Z - -=- postgres -> ephemeral
2020-12-10T17:59:44.963Z - postgres -> ephemeral took 16.395s, rate ~ 1.9M/s count: 31415926
2020-12-10T17:59:45.005Z - -=- ephemeral <-> postgres
2020-12-10T18:00:01.588Z - ephemeral <-> postgres took 16.582s, rate ~ 1.9M/s count: 31415926
2020-12-10T18:00:01.588Z - Verified ephemeral <-> postgres:
2020-12-10T18:00:01.588Z - [2020-01-01T00:00:00Z, 2020-12-29T14:38:45Z](31415926) Equal

2020-12-10T18:00:01.588Z - -=- ephemeral -> mysql
2020-12-10T18:08:21.848Z - ephemeral -> mysql took 8m20.259s, rate ~ 62.8k/s count: 31415926
2020-12-10T18:08:21.999Z - -=- mysql -> ephemeral
2020-12-10T18:09:05.233Z - mysql -> ephemeral took 43.234s, rate ~ 726.6k/s count: 31415926
2020-12-10T18:09:05.233Z - -=- ephemeral <-> mysql
2020-12-10T18:09:49.996Z - ephemeral <-> mysql took 44.763s, rate ~ 701.8k/s count: 31415926
2020-12-10T18:09:49.996Z - Verified ephemeral <-> mysql:
2020-12-10T18:09:49.996Z - [2020-01-01T00:00:00Z, 2020-12-29T14:38:45Z](31415926) Equal

real	17m11.313s
```

### Should check monthly sums after snapshots

2015-09-28:

|    format |  avg | total |
|----------:|-----:|------:|
|    .jsonl | 100M | 9026M |
| .jsonl.gz |   7M |  629M |
| jsonl.bz2 |   5M |  384M |

```bash
time md5sum $(find data/jsonl/month -type f -name \*.jsonl) | tee sums.txt
md5sum -c sums.txt
md5sum $(find data/jsonl/month -type f -name \*.jsonl)|cut -d \  -f 1|sort|md5sum
```

### ted.20150928.1006.sql.bz2 (watt, and native...?)

```bash
mysql> select min(stamp),max(stamp),count(*) from watt;
+---------------------+---------------------+-----------+
| min(stamp)          | max(stamp)          | count(*)  |
+---------------------+---------------------+-----------+
| 2008-07-30 00:04:40 | 2015-09-28 14:06:52 | 212737945 |
+---------------------+---------------------+-----------+
1 row in set (0.00 sec)
```

## Gaps (< 1h)

I also use:

```sql
select min(stamp),max(stamp),count(*) from watt group by left(stamp,10);
```

### Gaps ted.20150928.1006.sql.bz2 (watt, and native...?)

```bash
mysql> select min(stamp),max(stamp),count(*) from watt;
+---------------------+---------------------+-----------+
| min(stamp)          | max(stamp)          | count(*)  |
+---------------------+---------------------+-----------+
| 2008-07-30 00:04:40 | 2015-09-28 14:06:52 | 212737945 |
+---------------------+---------------------+-----------+
1 row in set (0.00 sec)

mysql> select min(stamp),max(stamp),count(*) from ted_native;
+---------------------+---------------------+-----------+
| min(stamp)          | max(stamp)          | count(*)  |
+---------------------+---------------------+-----------+
| 2008-12-17 19:37:16 | 2015-09-28 14:06:52 | 202176877 |
+---------------------+---------------------+-----------+

mysql> select min(stamp),max(stamp),count(*) from ted_service;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-11-14 23:18:13 | 2008-12-17 19:19:20 |  2710592 |
+---------------------+---------------------+----------+

progress.Gaps: 2008-09-17T16:24:01Z 2008-09-18T03:28:17Z : 11h4m16s
progress.Gaps: 2008-10-12T23:19:28Z 2008-10-15T01:56:16Z : 50h36m48s
progress.Gaps: 2008-11-28T05:45:30Z 2008-11-28T17:10:53Z : 11h25m23s
progress.Gaps: 2008-12-02T09:09:55Z 2008-12-02T16:24:27Z : 7h14m32s
progress.Gaps: 2009-01-23T14:33:30Z 2009-01-25T00:41:57Z : 34h8m27s
progress.Gaps: 2009-04-05T22:12:33Z 2009-04-06T04:44:27Z : 6h31m54s
progress.Gaps: 2009-10-22T07:01:15Z 2009-10-22T20:13:02Z : 13h11m47s
progress.Gaps: 2009-10-24T01:38:40Z 2009-10-24T15:09:55Z : 13h31m15s
progress.Gaps: 2009-12-11T21:50:15Z 2009-12-13T17:38:53Z : 43h48m38s
progress.Gaps: 2010-07-23T01:52:32Z 2010-09-09T06:53:08Z : 1157h0m36s
progress.Gaps: 2010-12-07T13:34:55Z 2010-12-09T01:55:23Z : 36h20m28s
progress.Gaps: 2011-05-19T18:56:59Z 2011-05-19T22:51:10Z : 3h54m11s
progress.Gaps: 2011-09-04T07:21:50Z 2011-09-24T16:23:19Z : 489h1m29s
progress.Gaps: 2012-07-14T21:22:01Z 2012-07-15T05:27:01Z : 8h5m0s
progress.Gaps: 2012-07-17T12:33:41Z 2012-07-17T22:53:13Z : 10h19m32s
progress.Gaps: 2012-11-17T13:30:59Z 2012-11-17T20:22:14Z : 6h51m15s
progress.Gaps: 2013-02-27T22:49:25Z 2013-02-28T03:10:04Z : 4h20m39s
progress.Gaps: 2013-06-13T12:45:12Z 2013-06-14T07:54:07Z : 19h8m55s
progress.Gaps: 2013-06-26T09:56:07Z 2013-06-26T13:22:35Z : 3h26m28s
progress.Gaps: 2013-09-19T13:56:21Z 2013-09-20T03:05:28Z : 13h9m7s
progress.Gaps: 2013-10-01T12:44:24Z 2013-10-02T02:19:42Z : 13h35m18s
progress.Gaps: 2013-10-18T13:53:18Z 2013-10-18T15:53:27Z : 2h0m9s
progress.Gaps: 2014-01-18T13:42:25Z 2014-01-18T15:41:53Z : 1h59m28s
progress.Gaps: 2014-02-19T13:58:18Z 2014-02-20T01:44:19Z : 11h46m1s
progress.Gaps: 2014-04-22T11:51:38Z 2014-04-22T23:00:51Z : 11h9m13s
progress.Gaps: 2014-05-16T11:21:02Z 2014-05-17T03:46:19Z : 16h25m17s
progress.Gaps: 2014-07-17T12:34:53Z 2014-07-17T16:44:49Z : 4h9m56s
progress.Gaps: 2014-08-16T06:35:34Z 2014-08-16T07:57:15Z : 1h21m41s
progress.Gaps: 2014-08-16T07:57:21Z 2014-08-16T19:14:39Z : 11h17m18s
progress.Gaps: 2014-08-17T05:06:51Z 2014-08-17T06:13:56Z : 1h7m5s
progress.Gaps: 2014-08-17T08:53:56Z 2014-08-18T02:18:42Z : 17h24m46s
progress.Gaps: 2014-08-19T19:46:47Z 2014-08-19T21:07:44Z : 1h20m57s
progress.Gaps: 2014-08-23T02:17:46Z 2014-08-27T07:41:48Z : 101h24m2s
progress.Gaps: 2014-09-01T11:56:18Z 2014-09-01T23:12:34Z : 11h16m16s
progress.Gaps: 2014-09-05T19:50:02Z 2014-09-05T22:42:41Z : 2h52m39s
progress.Gaps: 2014-09-29T17:40:37Z 2014-09-29T21:37:11Z : 3h56m34s
progress.Gaps: 2014-10-05T14:02:43Z 2014-10-05T17:05:05Z : 3h2m22s
progress.Gaps: 2014-11-22T17:59:14Z 2014-12-03T10:00:17Z : 256h1m3s
progress.Gaps: 2014-12-21T05:36:29Z 2014-12-21T09:09:03Z : 3h32m34s
progress.Gaps: 2014-12-25T19:47:01Z 2015-01-04T08:41:46Z : 228h54m45s
progress.Gaps: 2015-03-11T17:36:52Z 2015-03-11T20:09:10Z : 2h32m18s
progress.Gaps: 2015-03-15T13:56:24Z 2015-03-15T16:09:04Z : 2h12m40s
progress.Gaps: 2015-05-13T14:55:34Z 2015-05-13T16:09:08Z : 1h13m34s
progress.Gaps: 2015-06-01T18:33:34Z 2015-06-01T20:09:05Z : 1h35m31s
progress.Gaps: 2015-06-23T18:58:24Z 2015-06-23T20:09:05Z : 1h10m41s
progress.Gaps: 2015-06-27T10:30:07Z 2015-06-27T12:09:04Z : 1h38m57s
progress.Gaps: 2015-06-29T10:54:22Z 2015-06-29T12:09:04Z : 1h14m42s
progress.Gaps: 2015-07-19T06:05:33Z 2015-07-19T08:09:11Z : 2h3m38s
progress.Gaps: 2015-07-22T21:56:04Z 2015-07-23T16:09:07Z : 18h13m3s
progress.Gaps: 2015-08-07T02:50:53Z 2015-08-07T04:09:04Z : 1h18m11s
progress.Gaps: 2015-08-27T14:56:34Z 2015-08-27T16:09:03Z : 1h12m29s
progress.Gaps: 2015-08-28T08:03:20Z 2015-08-29T00:29:19Z : 16h25m59s
progress.Gaps: 2015-09-27T22:07:10Z 2015-09-28T02:56:46Z : 4h49m36s

Progress.Gaps: 53 gaps totaling 2703h29m23s (9732563 entries)
Progress.Gaps: 212737945 total entries in [2008-07-30T00:04:40Z, 2015-09-28T14:06:52Z] 62798h2m12s
```

### Gaps ted.watt.2016-02-14-1555.sql.bz2

```bash
progress.Gaps: 2008-09-17T16:24:01Z 2008-09-18T03:28:17Z : 11h4m16s
progress.Gaps: 2008-10-12T23:19:28Z 2008-10-15T01:56:16Z : 50h36m48s
progress.Gaps: 2008-11-28T05:45:30Z 2008-11-28T17:10:53Z : 11h25m23s
progress.Gaps: 2008-12-02T09:09:55Z 2008-12-02T16:24:27Z : 7h14m32s
progress.Gaps: 2009-01-23T14:33:30Z 2009-01-25T00:41:57Z : 34h8m27s
progress.Gaps: 2009-04-05T22:12:33Z 2009-04-06T04:44:27Z : 6h31m54s
progress.Gaps: 2009-10-22T07:01:15Z 2009-10-22T20:13:02Z : 13h11m47s
progress.Gaps: 2009-10-24T01:38:40Z 2009-10-24T15:09:55Z : 13h31m15s
progress.Gaps: 2009-12-11T21:50:15Z 2009-12-13T17:38:53Z : 43h48m38s
progress.Gaps: 2010-07-23T01:52:32Z 2010-09-09T06:53:08Z : 1157h0m36s
progress.Gaps: 2010-12-07T13:34:55Z 2010-12-09T01:55:23Z : 36h20m28s
progress.Gaps: 2011-05-19T18:56:59Z 2011-05-19T22:51:10Z : 3h54m11s
progress.Gaps: 2011-09-04T07:21:50Z 2011-09-24T16:23:19Z : 489h1m29s
progress.Gaps: 2012-07-14T21:22:01Z 2012-07-15T05:27:01Z : 8h5m0s
progress.Gaps: 2012-07-17T12:33:41Z 2012-07-17T22:53:13Z : 10h19m32s
progress.Gaps: 2012-11-17T13:30:59Z 2012-11-17T20:22:14Z : 6h51m15s
progress.Gaps: 2013-02-27T22:49:25Z 2013-02-28T03:10:04Z : 4h20m39s
progress.Gaps: 2013-06-13T12:45:12Z 2013-06-14T07:54:07Z : 19h8m55s
progress.Gaps: 2013-06-26T09:56:07Z 2013-06-26T13:22:35Z : 3h26m28s
progress.Gaps: 2013-09-19T13:56:21Z 2013-09-20T03:05:28Z : 13h9m7s
progress.Gaps: 2013-10-01T12:44:24Z 2013-10-02T02:19:42Z : 13h35m18s
progress.Gaps: 2013-10-18T13:53:18Z 2013-10-18T15:53:27Z : 2h0m9s
progress.Gaps: 2014-01-18T13:42:25Z 2014-01-18T15:41:53Z : 1h59m28s
progress.Gaps: 2014-02-19T13:58:18Z 2014-02-20T01:44:19Z : 11h46m1s
progress.Gaps: 2014-04-22T11:51:38Z 2014-04-22T23:00:51Z : 11h9m13s
progress.Gaps: 2014-05-16T11:21:02Z 2014-05-17T03:46:19Z : 16h25m17s
progress.Gaps: 2014-07-17T12:34:53Z 2014-07-17T16:44:49Z : 4h9m56s
progress.Gaps: 2014-08-16T06:35:34Z 2014-08-16T07:57:15Z : 1h21m41s
progress.Gaps: 2014-08-16T07:57:21Z 2014-08-16T19:14:39Z : 11h17m18s
progress.Gaps: 2014-08-17T05:06:51Z 2014-08-17T06:13:56Z : 1h7m5s
progress.Gaps: 2014-08-17T08:53:56Z 2014-08-18T02:18:42Z : 17h24m46s
progress.Gaps: 2014-08-19T19:46:47Z 2014-08-19T21:07:44Z : 1h20m57s
progress.Gaps: 2014-08-23T02:17:46Z 2014-08-27T07:41:48Z : 101h24m2s
progress.Gaps: 2014-09-01T11:56:18Z 2014-09-01T23:12:34Z : 11h16m16s
progress.Gaps: 2014-09-05T19:50:02Z 2014-09-05T22:42:41Z : 2h52m39s
progress.Gaps: 2014-09-29T17:40:37Z 2014-09-29T21:37:11Z : 3h56m34s
progress.Gaps: 2014-10-05T14:02:43Z 2014-10-05T17:05:05Z : 3h2m22s
progress.Gaps: 2014-11-22T17:59:14Z 2014-12-03T10:00:17Z : 256h1m3s
progress.Gaps: 2014-12-21T05:36:29Z 2014-12-21T09:09:03Z : 3h32m34s
progress.Gaps: 2014-12-25T19:47:01Z 2015-01-04T08:41:46Z : 228h54m45s
progress.Gaps: 2015-03-11T17:36:52Z 2015-03-11T20:09:10Z : 2h32m18s
progress.Gaps: 2015-03-15T13:56:24Z 2015-03-15T16:09:04Z : 2h12m40s
progress.Gaps: 2015-05-13T14:55:34Z 2015-05-13T16:09:08Z : 1h13m34s
progress.Gaps: 2015-06-01T18:33:34Z 2015-06-01T20:09:05Z : 1h35m31s
progress.Gaps: 2015-06-23T18:58:24Z 2015-06-23T20:09:05Z : 1h10m41s
progress.Gaps: 2015-06-27T10:30:07Z 2015-06-27T12:09:04Z : 1h38m57s
progress.Gaps: 2015-06-29T10:54:22Z 2015-06-29T12:09:04Z : 1h14m42s
progress.Gaps: 2015-07-19T06:05:33Z 2015-07-19T08:09:11Z : 2h3m38s
progress.Gaps: 2015-07-22T21:56:04Z 2015-07-23T16:09:07Z : 18h13m3s
progress.Gaps: 2015-08-07T02:50:53Z 2015-08-07T04:09:04Z : 1h18m11s
progress.Gaps: 2015-08-27T14:56:34Z 2015-08-27T16:09:03Z : 1h12m29s
progress.Gaps: 2015-08-28T08:03:20Z 2015-08-29T00:29:19Z : 16h25m59s
progress.Gaps: 2015-09-27T22:07:10Z 2015-09-28T02:56:46Z : 4h49m36s
progress.Gaps: 2015-11-30T12:54:48Z 2015-11-30T14:09:07Z : 1h14m19s
progress.Gaps: 2015-12-15T18:34:04Z 2015-12-15T21:09:03Z : 2h34m59s
progress.Gaps: 2016-01-12T19:03:01Z 2016-01-12T21:09:04Z : 2h6m3s
progress.Gaps: 2016-01-24T19:52:42Z 2016-01-29T21:09:03Z : 121h16m21s
progress.Gaps: 2016-02-01T19:47:19Z 2016-02-03T03:03:04Z : 31h15m45s
progress.Gaps: 2016-02-03T08:52:05Z 2016-02-11T06:29:39Z : 189h37m34s

Progress.Gaps: 59 gaps totaling 3051h34m24s (10985664 entries)
Progress.Gaps: 223101124 total entries in [2008-07-30T00:04:40Z, 2016-02-11T12:22:45Z] 66060h18m5s
```

### ted.watt-just2016.2016-02-14-1624.sql.bz2

progress.Gaps: 2016-01-12T19:03:01Z 2016-01-12T21:09:04Z : 2h6m3s
progress.Gaps: 2016-01-24T19:52:42Z 2016-01-29T21:09:03Z : 121h16m21s
progress.Gaps: 2016-02-01T19:47:19Z 2016-02-03T03:03:04Z : 31h15m45s
progress.Gaps: 2016-02-03T08:52:05Z 2016-02-11T06:29:39Z : 189h37m34s
progress.Gaps: 2016-02-11T12:22:45Z 2016-02-14T21:09:15Z : 80h46m30s

Progress.Gaps: 5 gaps totaling 425h2m13s (1530133 entries)
Progress.Gaps: 2318726 total entries in [2016-01-01T00:00:00Z, 2016-02-14T21:24:21Z] 1077h24m21s

### ted.watt.20180326.0312Z.sql.bz2

```bash
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2016-03-12 06:35:35 | 2018-03-26 03:12:25 | 61205052 |
+---------------------+---------------------+----------+
progress.Gaps: 2016-06-20T23:52:47Z 2016-06-21T01:20:36Z : 1h27m49s
progress.Gaps: 2016-06-22T03:08:51Z 2016-06-23T04:01:24Z : 24h52m33s
progress.Gaps: 2016-11-06T14:02:10Z 2016-11-06T15:16:45Z : 1h14m35s
progress.Gaps: 2017-01-04T11:23:42Z 2017-01-04T14:59:51Z : 3h36m9s
progress.Gaps: 2017-01-04T22:20:28Z 2017-01-05T05:25:43Z : 7h5m15s
progress.Gaps: 2017-02-07T18:06:02Z 2017-02-08T01:21:15Z : 7h15m13s
progress.Gaps: 2017-03-08T22:20:24Z 2017-03-09T05:23:14Z : 7h2m50s
progress.Gaps: 2017-05-30T16:01:34Z 2017-05-31T03:59:29Z : 11h57m55s
progress.Gaps: 2017-07-19T11:28:18Z 2017-07-26T03:36:58Z : 160h8m40s
progress.Gaps: 2017-10-18T14:20:52Z 2017-10-25T05:03:04Z : 158h42m12s
progress.Gaps: 2017-11-03T04:27:15Z 2017-11-07T20:41:03Z : 112h13m48s
progress.Gaps: 2017-12-28T17:54:31Z 2017-12-28T21:58:12Z : 4h3m41s
progress.Gaps: 2018-03-08T09:10:28Z 2018-03-10T08:59:28Z : 47h49m0s

Progress.Gaps: 13 gaps totaling 547h29m40s (1970980 entries)
Progress.Gaps: 61205052 total entries in [2016-03-12T06:35:35Z, 2018-03-26T03:12:25Z] 17852h36m50s
```
