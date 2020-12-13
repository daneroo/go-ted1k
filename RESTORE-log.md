# Restoring snapshots

Move a condensed version of this into main README.md

## Systematic restore of all snapshots

- Excluding any footprints
- `ted.200YYMMDD.HHMM.sql.bz2` contains `ted_native`, `ted_service` and `watt` tables
  - we can ignore ted_service
    - seems in error, probably a time shift, we can ignore
    - has a problem with 2008-11-14T23:18:13Z - 2008-11-28T04:59:59Z
    - was in use only for [2008-11-14 23:18:13 , 2008-12-17 19:19:20] 
  - we can ignore `watt_day|hour|minute|tensec`

The tables `ted_native` and `watt`, should be equivalent starting at 2008-12-17T19:37:16Z, when the TedNative capture started. We only have dumps of both tables until: `ted.20150928.1006.sql.bz2`, and the last stamp is `2015-09-28T14:06:52Z`, at that point there are only 25 entries present in ted_native and missing in watt:

```txt
2020-12-11T05:41:34.294Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
```

So we combine both in a jsonl dump: `./data/jsonl-ted-rollup.20150928.1006/`, and confirm by restoring each snapshot (ted.200YYMMDD.HHMM.sql.bz2) and confirm that all data is in that combined json.

Phase-1: We then restored each `ted.200YYMMDD.HHMM.sql.bz2` to verify that all samples of `ted_native` and `watt` tables were included in the rollup. The output is in `RESTORE-phase-1.md`, (which took	24 hours), confirmint that there were only `MissingInB`, that is no Conflict or MissingInA entries.

Final rollup we will accumulate in postgres.
- seed with ./data/jsonl-ted-rollup.20150928.100 (result of phase-1 above)
- restore each `ted.watt.200YYMMDD.HHMM.sql.bz2`

### /archive/mirror/ted/ted.20090214.1756.sql.bz2
```
- Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1756.sql.bz2
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 22:55:47	15495136
- Done Restoring Database
real	9m6.224s

JSON Size: 660M
Summary of mysql restore: 

2020-12-10T22:01:13.576Z - Verified jsonl <-> mysql(watt):
2020-12-10T22:01:13.576Z - [2008-07-30T00:04:40Z, 2009-02-14T22:55:47Z](15495136) Equal
2020-12-10T22:01:13.576Z - -=- jsonl <-> mysql(ted_native)
2020-12-10T22:01:56.657Z - jsonl <-> mysql(ted_native) took 43.081s, rate ~ 114.5k/s count: 4934047
2020-12-10T22:01:56.657Z - Verified jsonl <-> mysql(ted_native):
2020-12-10T22:01:56.657Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-10T22:01:56.657Z - [2008-12-17T19:37:16Z, 2009-02-14T22:55:47Z](4934047) Equal
```

#### analysis

```
mysql> select min(stamp),max(stamp),count(*) from ted_native;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-12-17 19:37:16 | 2009-02-14 22:55:47 |  4934047 |
+---------------------+---------------------+----------+
1 row in set (0.00 sec)

mysql> select min(stamp),max(stamp),count(*) from ted_service;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-11-14 23:18:13 | 2008-12-17 19:19:20 |  2710592 |
+---------------------+---------------------+----------+
1 row in set (0.00 sec)


mysql> select min(stamp),max(stamp),count(*) from watt;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-07-30 00:04:40 | 2009-02-14 22:55:47 | 15495136 |
+---------------------+---------------------+----------+

4934047+2710592 = 7644639

mysql> select min(stamp),max(stamp),count(*) from watt where stamp>'2008-12-17 19:19:20';
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-12-17 19:37:16 | 2009-02-14 22:55:47 |  4934047 |
+---------------------+---------------------+----------+
4934047+2710592 = 7644639 + 4934047
```

### /archive/mirror/ted/ted.20150928.1006.sql.bz2

- Restoring database from snapshot: ./data/archive/mirror/ted/ted.20150928.1006.sql.bz2
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2015-09-28 14:06:52	212737945
- Done Restoring Database
real	162m43.936s

### setup
```bash
mkdir -p ./data/archive/mirror/ted/
rsync -av --progress dirac:/Volumes/Space/archive/mirror/ted/ted.20*.sql.bz2 ./data/archive/mirror/ted/
```

### list of restores
```bash
2020-12-11T19:43:35.600Z start
51	/Volumes/Space/archive/mirror/ted/ted.20090214.1756.sql.bz2
51	/Volumes/Space/archive/mirror/ted/ted.20090214.1800.sql.bz2
64	/Volumes/Space/archive/mirror/ted/ted.20090326.1052.sql.bz2
65	/Volumes/Space/archive/mirror/ted/ted.20090328.1335.sql.bz2
100	/Volumes/Space/archive/mirror/ted/ted.20090528.0815.sql.bz2
100	/Volumes/Space/archive/mirror/ted/ted.20090609.0858.sql.bz2
132	/Volumes/Space/archive/mirror/ted/ted.20090918.0240.sql.bz2
2020-12-11T23:06:45.155Z start
148	/Volumes/Space/archive/mirror/ted/ted.20091022.0301.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.20091102.0131.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.20091113.2035.sql.bz2
2020-12-12T01:26:06.140Z start
308	/Volumes/Space/archive/mirror/ted/ted.20110406.0317.sql.bz2
324	/Volumes/Space/archive/mirror/ted/ted.20110607.0118.sql.bz2
372	/Volumes/Space/archive/mirror/ted/ted.20111017.2034.sql.bz2
452	/Volumes/Space/archive/mirror/ted/ted.20120608.0122.sql.bz2
532	/Volumes/Space/archive/mirror/ted/ted.20130221.2122.sql.bz2
660	/Volumes/Space/archive/mirror/ted/ted.20140219.2021.sql.bz2
708	/Volumes/Space/archive/mirror/ted/ted.20140806.0019.sql.bz2
836	/Volumes/Space/archive/mirror/ted/ted.20150928.1006.sql.bz2
2020-12-12T20:23:26.248Z done

7   /Volumes/Space/archive/mirror/ted/ted.watt-just2016.2016-02-14-1624.sql.gz

63	/Volumes/Space/archive/mirror/ted/ted.watt.20090918.0300.sql.bz2
68	/Volumes/Space/archive/mirror/ted/ted.watt.20091022.0258.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20091102.0134.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20110406.0316.sql.bz2
164	/Volumes/Space/archive/mirror/ted/ted.watt.20110607.0115.sql.bz2
180	/Volumes/Space/archive/mirror/ted/ted.watt.20111017.2045.sql.bz2
212	/Volumes/Space/archive/mirror/ted/ted.watt.20120608.0119.sql.bz2
260	/Volumes/Space/archive/mirror/ted/ted.watt.20130221.2119.sql.bz2
308	/Volumes/Space/archive/mirror/ted/ted.watt.20140219.2038.sql.bz2
340	/Volumes/Space/archive/mirror/ted/ted.watt.20140806.0016.sql.bz2
340	/Volumes/Space/archive/mirror/ted/ted.watt.20141005.2218.sql.bz2
388	/Volumes/Space/archive/mirror/ted/ted.watt.20150928.1003.sql.bz2
420	/Volumes/Space/archive/mirror/ted/ted.watt.2016-02-14-1555.sql.bz2

8   /Volumes/Space/archive/mirror/ted/ted.watt.20160430.0232Z.sql.bz2
16	/Volumes/Space/archive/mirror/ted/ted.watt.20160616.0229Z.sql.bz2
21	/Volumes/Space/archive/mirror/ted/ted.watt.20160719.1848Z.sql.bz2
31	/Volumes/Space/archive/mirror/ted/ted.watt.20160918.0059Z.sql.bz2
43	/Volumes/Space/archive/mirror/ted/ted.watt.20161202.0733Z.sql.bz2
49	/Volumes/Space/archive/mirror/ted/ted.watt.20170106.0629Z.sql.bz2
63	/Volumes/Space/archive/mirror/ted/ted.watt.20170326.1528Z.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20170607.0541Z.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20170727.1724Z.sql.bz2
116	/Volumes/Space/archive/mirror/ted/ted.watt.20180217.2219Z.sql.bz2
132	/Volumes/Space/archive/mirror/ted/ted.watt.20180326.0312Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180612.0035Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180720.2138Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180831.2033Z.sql.bz2
164	/Volumes/Space/archive/mirror/ted/ted.watt.20181024.1913Z.sql.bz2
195	/Volumes/Space/archive/mirror/ted/ted.watt.20190414.0128Z.sql.bz2
201	/Volumes/Space/archive/mirror/ted/ted.watt.20190617.0443Z.sql.bz2
209	/Volumes/Space/archive/mirror/ted/ted.watt.20190818.0554Z.sql.bz2
233	/Volumes/Space/archive/mirror/ted/ted.watt.20191129.0710Z.sql.bz2
251	/Volumes/Space/archive/mirror/ted/ted.watt.20200413.1503Z.sql.bz2
269	/Volumes/Space/archive/mirror/ted/ted.watt.20200807.2218Z.sql.bz2
297	/Volumes/Space/archive/mirror/ted/ted.watt.20201120.2332Z.sql.bz2
```

## Oldest TedNative sqlite dump

- Data capture started at 2008-07-30T00:04:40Z
- Switchover from TedNative (sqlite3) to python v1 happened at 2008-12-17T19:37:16Z

| stamp               | watt |                                 |
| ------------------- | ---- | ------------------------------- |
| 2008-07-30 00:04:40 | 540  |
| ....                | ...  |
| 2008-12-17 19:19:19 | 570  |
| 2008-12-17 19:19:20 | 560  | <- last sqlite TEDNative insert |
| 2008-12-17 19:37:16 | 610  | <- first mysql python insert    |
| 2008-12-17 19:37:17 | 610  |

From `TED.db.20081217.1419.sqlite` <- `/archive/mirror/ted/footprints/end-of-life/TED.db.20081217.1419.bz2`,
and from `.../im-ted1k/legacy/scalr-utils/scalr.py`,
and finally `docker run --rm -it python:2.7`

```sqlite
sqlite> select min(tick),max(tick),count(*) from rdu_second_data;
0633529550800006250|0633651203640000000|10242588
```

```python
import string
def tedToSecs(tedTimeString):
  millis = string.atol(tedTimeString)/10000
  return millis / 1000 - 62135578800;

import time
def tedToLocal(tedTimeString):
  secs = tedToSecs(tedTimeString)
  return time.strftime("%Y-%m-%d %H:%M:%S %Z",time.localtime(secs))

print tedToSecs('0633529550800006250') # 1217376280
print tedToSecs('0633651203640000000')# 1229541564

print tedToLocal('0633529550800006250') # 2008-07-30 00:04:40 UTC
print tedToLocal('0633651203640000000')# 2008-12-17 19:19:24 UTC
```

## ted.watt.2016-02-14-1555.sql.bz2

- Restoring database from snapshot: /Users/daniel/Downloads/ted/ted.watt.2016-02-14-1555.sql.bz2
  To docker container named: go-ted1k_mysql_1
  Using database: ted, MYSQL_USER=ted
  Data Volume will persisted inside that docker container

- Verifying docker environment
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_mysql_1) to accept connections (max 30 seconds)
  Connected
- Restoring database...

- Expect something recent in watt table
  min(stamp) max(stamp) count(\*)
  2008-07-30 00:04:40 2016-02-11 12:22:45 223101124

- Done Restoring Database

  2976.253s

In `ted.watt.2016-02-14-1555.sql.bz2`, there are 10,561,089 samples before the TEDNative cutoff,
whereas the last sqlite snapshot `/archive/mirror/ted/footprints/end-of-life/TED.db.20081217.1419.bz2` has only 10,242,588 samples, so it seems there are 318,501 more samples in the mysql snapshot over this period. Hmm.

mysql> select min(stamp),max(stamp),count(_) from watt where stamp<'2008-12-17 19:19:30';
+---------------------+---------------------+----------+
| min(stamp) | max(stamp) | count(_) |
+---------------------+---------------------+----------+
| 2008-07-30 00:04:40 | 2008-12-17 19:19:20 | 10561089 |
+---------------------+---------------------+----------+

Once restored into Mysql:8.0 this restore took ~6.3Gb

| size | path               |
| ---- | ------------------ |
| 8811 | /var/lib/mysql/    |
| 6385 | /var/lib/mysql/ted |

## exception .gz: ted.watt-just2016.2016-02-14-1624.sql.gz

2016-01-01 00:00:00 - 2016-02-14 21:24:21
