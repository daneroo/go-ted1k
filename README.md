# Ted1k with Go

## Todo

[See Evernote](https://www.evernote.com/shard/s60/nl/1773032759/ae1b9921-7e85-4b75-a21b-86be7d524295/)

- merge.Verify(a,b)
- find best writeBatchSize in mysql writer 1k...32k : currently 10k
- mv flux,ignore,jsonl,mysql to store/
- try pg,sqlite (general sql module)


## Vendoring - with vgo
Build and test:
```
vgo build ./scripts/pump.go
vgo test -v  ./...
```

See [this post](https://research.swtch.com/vgo-tour) for summary usage
```
go get -u golang.org/x/vgo
```

### Vendoring with govend (deprecated)
[Usage](https://github.com/govend/govend):

	govend -v  # download all the dependencies in the vendor.yml file
	govend -v -u # scan your project, update all dependencies, and update the vendor.yml revision versions

[New/Update]

	govend github.com/gorilla/mux  # add
	govend -u github.com/gorilla/mux # update

To install [`govend`](https://github.com/gophersaurus/govend) itself:

	go get -u github.com/govend/govend

and made sure our `GOPATH` was set and `$GOPATH/bin` is on our `$PATH`.

See this general vendoring entry: [Go/Wiki for reference](https://github.com/golang/go/wiki/PackageManagementTools).
Prior to `go1.6`, we also had to set `GO15VENDOREXPERIMENT=1`.

## InfluxDB
```
docker exec -it goted1k_tedflux_1 bash
influx -database ted -execute 'select count(value) from watt'
select mean(value)*24/1000 from watt where time > '2008-01-01' and time < '2016-01-01' group by time(7d)
```

### Downsampled time series

Truncate for D,M,Y: http://play.golang.org/p/PUNNHq9sh6

Continuous Queries are not appropriate for historical data loading.
I should implement my own select .. into (in go), using tablenames as in mysql

	select mean(value)*24/1000 into kwh_1d from watt where time > '2015-09-01' group by time(1d)

## Docker
We have abandoned data volumes for now.
`docker-compose` command brings up MySQL, InfluxDB and Grafana instances, and the `restore-db.sh` script restores a MySQL snapshot/

	docker-compose up -d
	./restore-db.sh

## Timing of MySQL reads
For timing of MySQL selects with maxCount results

From goedel to cantor

	3600: 989s
	3600*24: 357s
	3600*24*10: 324s

From Dirac to local docker:

    3600*24: 605s  (Read-only)

From Godel to local docker:

    3600*24: 294s,290s  (Read-only) Now 412s,405s, with IgnoreAll
	10000: --s  (Batch Writes)



## MySQL inserts are ridiculously slow
This is what we di to spead things up:
(all time measured on dirac/docker)

- Initial naive approach: 550 ins/s
- Prepared statements: 650 ins/s
- Transactions: 2500-3000 ins/s (depending on size of batch...)

## Historical aggregation
_We lost data from ( 2016-02-14 21:24:21 , 2016-03-12 06:35:35 ]_

### check monthly sums after snapshots...
2015-09-28:
   .jsonl     avg 100M, total 9026M
	 .jsonl.gz  avg   7M, total  629M
	 .jsonl.bz2 avg   5M, total  384M
```
time md5sum $(find data/jsonl/month -type f -name \*.jsonl) | tee sums.txt
md5sum -c sums.txt
md5sum $(find data/jsonl/month -type f -name \*.jsonl)|cut -d \  -f 1|sort|md5sum
```

### ted.20150928.1006.sql.bz2 (watt, and native...?)
```
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
```
select min(stamp),max(stamp),count(*) from watt group by left(stamp,10);
```

### ted.20150928.1006.sql.bz2 (watt, and native...?)
```
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

### ted.watt.2016-02-14-1555.sql.bz2
```
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
```
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