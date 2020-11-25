# Restoring snapshots

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
  To docker container named: go-ted1k_teddb_1
  Using database: ted, MYSQL_USER=ted
  Data Volume will persisted inside that docker container

- Verifying docker environment
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_teddb_1) to accept connections (max 30 seconds)
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

## ted.watt-just2016.2016-02-14-1624.sql.gz

2016-01-01 00:00:00 - 2016-02-14 21:24:21
