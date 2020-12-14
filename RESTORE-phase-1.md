# Restore phase 1

```go
verify("jsonl <-> mysql(watt)", jsonl.NewReader(), mysql.NewReader(db, "watt"))
verify("jsonl <-> mysql(ted_native)", jsonl.NewReader(), mysql.NewReader(db, "ted_native"))
```

```bash
daniel@dockerdev:~/go-ted1k$ time ./restore-db.sh 
- Verifying docker environment
WARNING: No swap limit support
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_mysql_1) to accept connections (max 30 seconds)
  Connected

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1756.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	9m7.466s
user	0m16.462s
sys	0m4.937s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-02-14 22:55:47	4934047
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 22:55:47	15495136

- run mysqlrestore
2020-12-11T19:53:35.600Z - Starting TED1K mysql restore
2020-12-11T19:53:35.602Z - Connected to MySQL
2020-12-11T19:53:35.602Z - -=- jsonl <-> mysql(watt)
2020-12-11T19:54:58.587Z - jsonl <-> mysql(watt) took 1m22.985s, rate ~ 186.7k/s count: 15495136
2020-12-11T20:00:48.723Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:00:48.723Z - [2008-07-30T00:04:40Z, 2009-02-14T22:55:47Z](15495136) Equal
2020-12-11T20:00:48.723Z - [2009-02-14T22:55:48Z, 2015-09-28T14:06:52Z](197242834) MissingInB
2020-12-11T20:00:48.723Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:01:33.334Z - jsonl <-> mysql(ted_native) took 44.611s, rate ~ 110.6k/s count: 4934047
2020-12-11T20:07:27.297Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:07:27.297Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:07:27.297Z - [2008-12-17T19:37:16Z, 2009-02-14T22:55:47Z](4934047) Equal
2020-12-11T20:07:27.297Z - [2009-02-14T22:55:48Z, 2015-09-28T14:06:52Z](197242834) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1800.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	9m3.722s
user	0m16.089s
sys	0m4.846s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-02-14 23:00:05	4934304
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 23:00:05	15495393

- run mysqlrestore
2020-12-11T20:16:37.288Z - Starting TED1K mysql restore
2020-12-11T20:16:37.309Z - Connected to MySQL
2020-12-11T20:16:37.309Z - -=- jsonl <-> mysql(watt)
2020-12-11T20:17:58.138Z - jsonl <-> mysql(watt) took 1m20.829s, rate ~ 191.7k/s count: 15495393
2020-12-11T20:23:47.680Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:23:47.680Z - [2008-07-30T00:04:40Z, 2009-02-14T23:00:05Z](15495393) Equal
2020-12-11T20:23:47.680Z - [2009-02-14T23:00:06Z, 2015-09-28T14:06:52Z](197242577) MissingInB
2020-12-11T20:23:47.680Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:24:32.667Z - jsonl <-> mysql(ted_native) took 44.987s, rate ~ 109.7k/s count: 4934304
2020-12-11T20:30:24.277Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:30:24.277Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:30:24.277Z - [2008-12-17T19:37:16Z, 2009-02-14T23:00:05Z](4934304) Equal
2020-12-11T20:30:24.277Z - [2009-02-14T23:00:06Z, 2015-09-28T14:06:52Z](197242577) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090326.1052.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	11m54.704s
user	0m21.218s
sys	0m6.202s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-03-26 14:52:59	8314404
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-03-26 14:52:59	18875493

- run mysqlrestore
2020-12-11T20:42:24.360Z - Starting TED1K mysql restore
2020-12-11T20:42:24.383Z - Connected to MySQL
2020-12-11T20:42:24.383Z - -=- jsonl <-> mysql(watt)
2020-12-11T20:44:05.108Z - jsonl <-> mysql(watt) took 1m40.725s, rate ~ 187.4k/s count: 18875493
2020-12-11T20:49:53.188Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:49:53.188Z - [2008-07-30T00:04:40Z, 2009-03-26T14:52:59Z](18875493) Equal
2020-12-11T20:49:53.188Z - [2009-03-26T14:53:00Z, 2015-09-28T14:06:52Z](193862477) MissingInB
2020-12-11T20:49:53.188Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:50:56.180Z - jsonl <-> mysql(ted_native) took 1m2.992s, rate ~ 132.0k/s count: 8314404
2020-12-11T20:56:43.698Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:56:43.698Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:56:43.698Z - [2008-12-17T19:37:16Z, 2009-03-26T14:52:59Z](8314404) Equal
2020-12-11T20:56:43.698Z - [2009-03-26T14:53:00Z, 2015-09-28T14:06:52Z](193862477) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090328.1335.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	12m11.046s
user	0m21.024s
sys	0m6.546s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-03-28 17:35:27	8494287
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-03-28 17:35:27	19055376

- run mysqlrestore
2020-12-11T21:08:59.818Z - Starting TED1K mysql restore
2020-12-11T21:08:59.844Z - Connected to MySQL
2020-12-11T21:08:59.844Z - -=- jsonl <-> mysql(watt)
2020-12-11T21:10:42.171Z - jsonl <-> mysql(watt) took 1m42.326s, rate ~ 186.2k/s count: 19055376
2020-12-11T21:16:31.380Z - Verified jsonl <-> mysql(watt):
2020-12-11T21:16:31.380Z - [2008-07-30T00:04:40Z, 2009-03-28T17:35:27Z](19055376) Equal
2020-12-11T21:16:31.380Z - [2009-03-28T17:35:28Z, 2015-09-28T14:06:52Z](193682594) MissingInB
2020-12-11T21:16:31.380Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T21:17:35.693Z - jsonl <-> mysql(ted_native) took 1m4.313s, rate ~ 132.1k/s count: 8494287
2020-12-11T21:23:24.028Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T21:23:24.028Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T21:23:24.028Z - [2008-12-17T19:37:16Z, 2009-03-28T17:35:27Z](8494287) Equal
2020-12-11T21:23:24.028Z - [2009-03-28T17:35:28Z, 2015-09-28T14:06:52Z](193682594) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090528.0815.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	16m11.319s
user	0m28.646s
sys	0m8.181s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-05-28 12:15:46	13668586
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-05-28 12:15:46	24229675

- run mysqlrestore
2020-12-11T21:39:40.675Z - Starting TED1K mysql restore
2020-12-11T21:39:40.735Z - Connected to MySQL
2020-12-11T21:39:40.735Z - -=- jsonl <-> mysql(watt)
2020-12-11T21:41:51.377Z - jsonl <-> mysql(watt) took 2m10.641s, rate ~ 185.5k/s count: 24229675
2020-12-11T21:47:27.620Z - Verified jsonl <-> mysql(watt):
2020-12-11T21:47:27.620Z - [2008-07-30T00:04:40Z, 2009-05-28T12:15:46Z](24229675) Equal
2020-12-11T21:47:27.620Z - [2009-05-28T12:15:47Z, 2015-09-28T14:06:52Z](188508295) MissingInB
2020-12-11T21:47:27.620Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T21:49:00.478Z - jsonl <-> mysql(ted_native) took 1m32.858s, rate ~ 147.2k/s count: 13668586
2020-12-11T21:54:36.027Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T21:54:36.027Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T21:54:36.028Z - [2008-12-17T19:37:16Z, 2009-05-28T12:15:46Z](13668586) Equal
2020-12-11T21:54:36.028Z - [2009-05-28T12:15:47Z, 2015-09-28T14:06:52Z](188508295) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090609.0858.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	16m52.486s
user	0m30.672s
sys	0m8.975s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-06-09 12:58:55	14697489
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-06-09 12:58:55	25258578

- run mysqlrestore
2020-12-11T22:11:33.552Z - Starting TED1K mysql restore
2020-12-11T22:11:33.578Z - Connected to MySQL
2020-12-11T22:11:33.578Z - -=- jsonl <-> mysql(watt)
2020-12-11T22:13:48.434Z - jsonl <-> mysql(watt) took 2m14.855s, rate ~ 187.3k/s count: 25258578
2020-12-11T22:19:24.506Z - Verified jsonl <-> mysql(watt):
2020-12-11T22:19:24.506Z - [2008-07-30T00:04:40Z, 2009-06-09T12:58:55Z](25258578) Equal
2020-12-11T22:19:24.506Z - [2009-06-09T12:58:56Z, 2015-09-28T14:06:52Z](187479392) MissingInB
2020-12-11T22:19:24.506Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T22:21:01.608Z - jsonl <-> mysql(ted_native) took 1m37.102s, rate ~ 151.4k/s count: 14697489
2020-12-11T22:26:36.025Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T22:26:36.026Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T22:26:36.026Z - [2008-12-17T19:37:16Z, 2009-06-09T12:58:55Z](14697489) Equal
2020-12-11T22:26:36.026Z - [2009-06-09T12:58:56Z, 2015-09-28T14:06:52Z](187479392) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090918.0240.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	23m59.822s
user	0m40.754s
sys	0m12.476s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-09-18 06:40:17	23312671
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-09-18 06:40:17	33873760

- run mysqlrestore
2020-12-11T22:50:40.882Z - Starting TED1K mysql restore
2020-12-11T22:50:40.912Z - Connected to MySQL
2020-12-11T22:50:40.912Z - -=- jsonl <-> mysql(watt)
2020-12-11T22:53:31.467Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 184.9k/s, took 2m50.554s, rate ~ 184.9k/s count: 31536000
2020-12-11T22:53:44.116Z - jsonl <-> mysql(watt) took 3m3.203s, rate ~ 184.9k/s count: 33873760
2020-12-11T22:59:04.371Z - Verified jsonl <-> mysql(watt):
2020-12-11T22:59:04.371Z - [2008-07-30T00:04:40Z, 2009-09-18T06:40:17Z](33873760) Equal
2020-12-11T22:59:04.371Z - [2009-09-18T06:40:18Z, 2015-09-28T14:06:52Z](178864210) MissingInB
2020-12-11T22:59:04.371Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T23:01:26.657Z - jsonl <-> mysql(ted_native) took 2m22.286s, rate ~ 163.8k/s count: 23312671
2020-12-11T23:06:45.155Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T23:06:45.155Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T23:06:45.155Z - [2008-12-17T19:37:16Z, 2009-09-18T06:40:17Z](23312671) Equal
2020-12-11T23:06:45.155Z - [2009-09-18T06:40:18Z, 2015-09-28T14:06:52Z](178864210) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091022.0301.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	27m1.348s
user	0m46.351s
sys	0m14.627s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-10-22 07:01:15	26221838
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-10-22 07:01:15	36782927

- run mysqlrestore
2020-12-11T23:33:51.864Z - Starting TED1K mysql restore
2020-12-11T23:33:51.874Z - Connected to MySQL
2020-12-11T23:33:51.874Z - -=- jsonl <-> mysql(watt)
2020-12-11T23:37:16.448Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 154.2k/s, took 3m24.573s, rate ~ 154.2k/s count: 31536000
2020-12-11T23:37:50.443Z - jsonl <-> mysql(watt) took 3m58.568s, rate ~ 154.2k/s count: 36782927
^[f2020-12-11T23:43:06.840Z - Verified jsonl <-> mysql(watt):
2020-12-11T23:43:06.840Z - [2008-07-30T00:04:40Z, 2009-10-22T07:01:15Z](36782927) Equal
2020-12-11T23:43:06.840Z - [2009-10-22T20:13:02Z, 2015-09-28T14:06:52Z](175955043) MissingInB
2020-12-11T23:43:06.840Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T23:46:10.781Z - jsonl <-> mysql(ted_native) took 3m3.94s, rate ~ 142.6k/s count: 26221838
2020-12-11T23:51:28.072Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T23:51:28.072Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T23:51:28.072Z - [2008-12-17T19:37:16Z, 2009-10-22T07:01:15Z](26221838) Equal
2020-12-11T23:51:28.072Z - [2009-10-22T20:13:02Z, 2015-09-28T14:06:52Z](175955043) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091102.0131.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	29m0.210s
user	0m46.468s
sys	0m15.921s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-11-02 06:31:14	27064883
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-11-02 06:31:14	37625972

- run mysqlrestore
2020-12-12T00:20:33.511Z - Starting TED1K mysql restore
2020-12-12T00:20:33.513Z - Connected to MySQL
2020-12-12T00:20:33.513Z - -=- jsonl <-> mysql(watt)
2020-12-12T00:23:56.071Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.7k/s, took 3m22.557s, rate ~ 155.7k/s count: 31536000
2020-12-12T00:24:34.853Z - jsonl <-> mysql(watt) took 4m1.34s, rate ~ 155.9k/s count: 37625972
2020-12-12T00:29:50.750Z - Verified jsonl <-> mysql(watt):
2020-12-12T00:29:50.750Z - [2008-07-30T00:04:40Z, 2009-11-02T06:31:14Z](37625972) Equal
2020-12-12T00:29:50.750Z - [2009-11-02T06:31:15Z, 2015-09-28T14:06:52Z](175111998) MissingInB
2020-12-12T00:29:50.750Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T00:33:01.838Z - jsonl <-> mysql(ted_native) took 3m11.088s, rate ~ 141.6k/s count: 27064883
2020-12-12T00:38:17.164Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T00:38:17.164Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T00:38:17.164Z - [2008-12-17T19:37:16Z, 2009-11-02T06:31:14Z](27064883) Equal
2020-12-12T00:38:17.164Z - [2009-11-02T06:31:15Z, 2015-09-28T14:06:52Z](175111998) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091113.2035.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	29m50.511s
user	0m47.759s
sys	0m16.531s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-11-14 01:35:35	28066796
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-11-14 01:35:35	38627885

- run mysqlrestore
2020-12-12T01:08:13.039Z - Starting TED1K mysql restore
2020-12-12T01:08:13.041Z - Connected to MySQL
2020-12-12T01:08:13.041Z - -=- jsonl <-> mysql(watt)
2020-12-12T01:11:35.395Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.355s, rate ~ 155.8k/s count: 31536000
2020-12-12T01:12:20.750Z - jsonl <-> mysql(watt) took 4m7.709s, rate ~ 155.9k/s count: 38627885
2020-12-12T01:17:35.557Z - Verified jsonl <-> mysql(watt):
2020-12-12T01:17:35.557Z - [2008-07-30T00:04:40Z, 2009-11-14T01:35:35Z](38627885) Equal
2020-12-12T01:17:35.557Z - [2009-11-14T01:48:36Z, 2015-09-28T14:06:52Z](174110085) MissingInB
2020-12-12T01:17:35.557Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T01:20:51.538Z - jsonl <-> mysql(ted_native) took 3m15.98s, rate ~ 143.2k/s count: 28066796
2020-12-12T01:26:06.140Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T01:26:06.140Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T01:26:06.140Z - [2008-12-17T19:37:16Z, 2009-11-14T01:35:35Z](28066796) Equal
2020-12-12T01:26:06.140Z - [2009-11-14T01:48:36Z, 2015-09-28T14:06:52Z](174110085) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20110406.0317.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	63m33.269s
user	1m41.827s
sys	0m35.158s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-04-06 07:15:35	67085417
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-04-06 07:15:35	77646482

- run mysqlrestore
2020-12-12T02:29:44.824Z - Starting TED1K mysql restore
2020-12-12T02:29:44.827Z - Connected to MySQL
2020-12-12T02:29:44.827Z - -=- jsonl <-> mysql(watt)
2020-12-12T02:33:07.631Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.5k/s, took 3m22.804s, rate ~ 155.5k/s count: 31536000
2020-12-12T02:36:28.536Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.0k/s, took 6m43.709s, rate ~ 156.2k/s count: 63072000
2020-12-12T02:38:01.010Z - jsonl <-> mysql(watt) took 8m16.183s, rate ~ 156.5k/s count: 77646482
2020-12-12T02:42:02.652Z - Verified jsonl <-> mysql(watt):
2020-12-12T02:42:02.652Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T02:42:02.652Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T02:42:02.652Z - [2009-11-17T03:00:30Z, 2011-04-06T07:15:35Z](38760401) Equal
2020-12-12T02:42:02.652Z - [2011-04-06T07:24:52Z, 2015-09-28T14:06:52Z](135091463) MissingInB
2020-12-12T02:42:02.652Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T02:45:42.851Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.2k/s, took 3m40.199s, rate ~ 143.2k/s count: 31536000
2020-12-12T02:49:02.894Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 7m0.241s, rate ~ 150.1k/s count: 63072000
2020-12-12T02:49:28.212Z - jsonl <-> mysql(ted_native) took 7m25.56s, rate ~ 150.6k/s count: 67085417
2020-12-12T02:53:30.145Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T02:53:30.145Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T02:53:30.145Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T02:53:30.145Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T02:53:30.145Z - [2009-11-17T03:02:49Z, 2011-04-06T07:15:35Z](38760400) Equal
2020-12-12T02:53:30.145Z - [2011-04-06T07:24:52Z, 2015-09-28T14:06:52Z](135091463) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20110607.0118.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	67m16.029s
user	1m47.745s
sys	0m36.741s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-06-07 05:18:07	72370533
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-06-07 05:18:07	82931598

- run mysqlrestore
2020-12-12T04:00:51.521Z - Starting TED1K mysql restore
2020-12-12T04:00:51.523Z - Connected to MySQL
2020-12-12T04:00:51.523Z - -=- jsonl <-> mysql(watt)
2020-12-12T04:04:13.991Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.467s, rate ~ 155.8k/s count: 31536000
2020-12-12T04:07:35.114Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 156.8k/s, took 6m43.59s, rate ~ 156.3k/s count: 63072000
2020-12-12T04:09:41.351Z - jsonl <-> mysql(watt) took 8m49.827s, rate ~ 156.5k/s count: 82931598
2020-12-12T04:13:35.374Z - Verified jsonl <-> mysql(watt):
2020-12-12T04:13:35.374Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T04:13:35.374Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T04:13:35.374Z - [2009-11-17T03:00:30Z, 2011-06-07T05:18:07Z](44045517) Equal
2020-12-12T04:13:35.374Z - [2011-06-07T05:18:08Z, 2015-09-28T14:06:52Z](129806347) MissingInB
2020-12-12T04:13:35.374Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T04:17:15.020Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.6k/s, took 3m39.645s, rate ~ 143.6k/s count: 31536000
2020-12-12T04:20:36.332Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 156.7k/s, took 7m0.957s, rate ~ 149.8k/s count: 63072000
2020-12-12T04:21:35.184Z - jsonl <-> mysql(ted_native) took 7m59.809s, rate ~ 150.8k/s count: 72370533
2020-12-12T04:25:28.341Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T04:25:28.341Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T04:25:28.341Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T04:25:28.341Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T04:25:28.341Z - [2009-11-17T03:02:49Z, 2011-06-07T05:18:07Z](44045516) Equal
2020-12-12T04:25:28.341Z - [2011-06-07T05:18:08Z, 2015-09-28T14:06:52Z](129806347) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20111017.2034.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	75m3.564s
user	2m0.566s
sys	0m41.336s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-10-18 00:33:13	81986218
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-10-18 00:33:13	92547284

- run mysqlrestore
2020-12-12T05:40:37.456Z - Starting TED1K mysql restore
2020-12-12T05:40:37.458Z - Connected to MySQL
2020-12-12T05:40:37.458Z - -=- jsonl <-> mysql(watt)
2020-12-12T05:44:01.130Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 154.8k/s, took 3m23.671s, rate ~ 154.8k/s count: 31536000
2020-12-12T05:47:24.944Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 154.7k/s, took 6m47.485s, rate ~ 154.8k/s count: 63072000
2020-12-12T05:50:35.005Z - jsonl <-> mysql(watt) took 9m57.547s, rate ~ 154.9k/s count: 92547284
2020-12-12T05:54:10.505Z - Verified jsonl <-> mysql(watt):
2020-12-12T05:54:10.505Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T05:54:10.505Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T05:54:10.505Z - [2009-11-17T03:00:30Z, 2011-10-18T00:33:13Z](53661203) Equal
2020-12-12T05:54:10.505Z - [2011-10-18T00:51:14Z, 2015-09-28T14:06:52Z](120190661) MissingInB
2020-12-12T05:54:10.505Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T05:57:50.793Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.2k/s, took 3m40.287s, rate ~ 143.2k/s count: 31536000
2020-12-12T06:01:11.576Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.1k/s, took 7m1.071s, rate ~ 149.8k/s count: 63072000
2020-12-12T06:03:10.927Z - jsonl <-> mysql(ted_native) took 9m0.422s, rate ~ 151.7k/s count: 81986218
2020-12-12T06:06:45.631Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T06:06:45.631Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T06:06:45.631Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T06:06:45.631Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T06:06:45.631Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T06:06:45.631Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T06:06:45.631Z - [2011-10-16T21:54:25Z, 2011-10-18T00:33:13Z](94912) Equal
2020-12-12T06:06:45.632Z - [2011-10-18T00:51:14Z, 2015-09-28T14:06:52Z](120190661) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20120608.0122.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	93m42.371s
user	2m28.268s
sys	0m51.801s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2012-06-08 05:22:10	101938261
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2012-06-08 05:22:10	112499327

- run mysqlrestore
2020-12-12T07:40:33.556Z - Starting TED1K mysql restore
2020-12-12T07:40:33.558Z - Connected to MySQL
2020-12-12T07:40:33.558Z - -=- jsonl <-> mysql(watt)
2020-12-12T07:43:52.555Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 158.5k/s, took 3m18.997s, rate ~ 158.5k/s count: 31536000
2020-12-12T07:47:09.874Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 159.8k/s, took 6m36.316s, rate ~ 159.1k/s count: 63072000
2020-12-12T07:50:26.057Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 160.7k/s, took 9m52.499s, rate ~ 159.7k/s count: 94608000
2020-12-12T07:52:16.787Z - jsonl <-> mysql(watt) took 11m43.229s, rate ~ 160.0k/s count: 112499327
2020-12-12T07:55:17.159Z - Verified jsonl <-> mysql(watt):
2020-12-12T07:55:17.159Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T07:55:17.159Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T07:55:17.159Z - [2009-11-17T03:00:30Z, 2012-06-08T05:22:10Z](73613246) Equal
2020-12-12T07:55:17.159Z - [2012-06-08T05:22:12Z, 2015-09-28T14:06:52Z](100238618) MissingInB
2020-12-12T07:55:17.159Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T07:58:55.970Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 144.1k/s, took 3m38.811s, rate ~ 144.1k/s count: 31536000
2020-12-12T08:02:16.134Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 6m58.975s, rate ~ 150.5k/s count: 63072000
2020-12-12T08:05:35.446Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 158.2k/s, took 10m18.287s, rate ~ 153.0k/s count: 94608000
2020-12-12T08:06:21.455Z - jsonl <-> mysql(ted_native) took 11m4.296s, rate ~ 153.5k/s count: 101938261
2020-12-12T08:09:20.793Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T08:09:20.793Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T08:09:20.793Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T08:09:20.793Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T08:09:20.793Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T08:09:20.793Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T08:09:20.793Z - [2011-10-16T21:54:25Z, 2012-06-08T05:22:10Z](20046955) Equal
2020-12-12T08:09:20.793Z - [2012-06-08T05:22:12Z, 2015-09-28T14:06:52Z](100238618) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20130221.2122.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	112m54.141s
user	2m57.062s
sys	1m1.232s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2013-02-22 02:22:26	123979218
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2013-02-22 02:22:26	134540284

- run mysqlrestore
2020-12-12T10:02:20.797Z - Starting TED1K mysql restore
2020-12-12T10:02:20.799Z - Connected to MySQL
2020-12-12T10:02:20.799Z - -=- jsonl <-> mysql(watt)
2020-12-12T10:05:42.034Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 156.7k/s, took 3m21.235s, rate ~ 156.7k/s count: 31536000
2020-12-12T10:09:02.179Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.6k/s, took 6m41.38s, rate ~ 157.1k/s count: 63072000
2020-12-12T10:12:21.672Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 158.1k/s, took 10m0.873s, rate ~ 157.5k/s count: 94608000
2020-12-12T10:15:40.255Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 158.8k/s, took 13m19.456s, rate ~ 157.8k/s count: 126144000
2020-12-12T10:16:33.353Z - jsonl <-> mysql(watt) took 14m12.555s, rate ~ 157.8k/s count: 134540284
2020-12-12T10:18:53.624Z - Verified jsonl <-> mysql(watt):
2020-12-12T10:18:53.624Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T10:18:53.624Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T10:18:53.624Z - [2009-11-17T03:00:30Z, 2013-02-22T02:22:26Z](95654203) Equal
2020-12-12T10:18:53.624Z - [2013-02-22T02:22:27Z, 2015-09-28T14:06:52Z](78197661) MissingInB
2020-12-12T10:18:53.624Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T10:22:32.904Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.8k/s, took 3m39.28s, rate ~ 143.8k/s count: 31536000
2020-12-12T10:25:52.975Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 6m59.351s, rate ~ 150.4k/s count: 63072000
2020-12-12T10:29:12.390Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 158.1k/s, took 10m18.766s, rate ~ 152.9k/s count: 94608000
2020-12-12T10:32:18.052Z - jsonl <-> mysql(ted_native) took 13m24.428s, rate ~ 154.1k/s count: 123979218
2020-12-12T10:34:38.073Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T10:34:38.073Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T10:34:38.073Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T10:34:38.073Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T10:34:38.073Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T10:34:38.073Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T10:34:38.073Z - [2011-10-16T21:54:25Z, 2013-02-22T02:22:26Z](42087912) Equal
2020-12-12T10:34:38.073Z - [2013-02-22T02:22:27Z, 2015-09-28T14:06:52Z](78197661) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20140219.2021.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	137m59.671s
user	3m40.321s
sys	1m16.615s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2014-02-19 13:58:18	154748431
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-02-19 13:58:18	165309499

- run mysqlrestore
2020-12-12T12:52:43.302Z - Starting TED1K mysql restore
2020-12-12T12:52:43.305Z - Connected to MySQL
2020-12-12T12:52:43.305Z - -=- jsonl <-> mysql(watt)
2020-12-12T12:56:05.707Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.402s, rate ~ 155.8k/s count: 31536000
2020-12-12T12:59:26.945Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 156.7k/s, took 6m43.64s, rate ~ 156.3k/s count: 63072000
2020-12-12T13:02:48.223Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 156.7k/s, took 10m4.918s, rate ~ 156.4k/s count: 94608000
2020-12-12T13:06:08.606Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 157.4k/s, took 13m25.3s, rate ~ 156.6k/s count: 126144000
2020-12-12T13:09:27.310Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 158.7k/s, took 16m44.004s, rate ~ 157.1k/s count: 157680000
2020-12-12T13:10:15.117Z - jsonl <-> mysql(watt) took 17m31.811s, rate ~ 157.2k/s count: 165309499
2020-12-12T13:11:40.424Z - Verified jsonl <-> mysql(watt):
2020-12-12T13:11:40.424Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T13:11:40.424Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T13:11:40.424Z - [2009-11-17T03:00:30Z, 2014-02-19T13:58:18Z](126423418) Equal
2020-12-12T13:11:40.424Z - [2014-02-20T01:44:19Z, 2015-09-28T14:06:52Z](47428446) MissingInB
2020-12-12T13:11:40.424Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T13:15:19.441Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 144.0k/s, took 3m39.016s, rate ~ 144.0k/s count: 31536000
2020-12-12T13:18:39.670Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.5k/s, took 6m59.245s, rate ~ 150.4k/s count: 63072000
2020-12-12T13:21:59.463Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 157.8k/s, took 10m19.038s, rate ~ 152.8k/s count: 94608000
2020-12-12T13:25:20.132Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 157.2k/s, took 13m39.708s, rate ~ 153.9k/s count: 126144000
2020-12-12T13:28:21.946Z - jsonl <-> mysql(ted_native) took 16m41.521s, rate ~ 154.5k/s count: 154748431
2020-12-12T13:29:46.838Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T13:29:46.838Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T13:29:46.839Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T13:29:46.839Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T13:29:46.839Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T13:29:46.839Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T13:29:46.839Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T13:29:46.839Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T13:29:46.839Z - [2013-10-18T16:06:13Z, 2014-02-19T13:58:18Z](10584273) Equal
2020-12-12T13:29:46.839Z - [2014-02-20T01:44:19Z, 2015-09-28T14:06:52Z](47428446) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20140806.0019.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	150m48.392s
user	4m3.086s
sys	1m25.752s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2014-08-06 04:19:27	168924060
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-08-06 04:19:27	179485128

- run mysqlrestore
2020-12-12T16:00:40.893Z - Starting TED1K mysql restore
2020-12-12T16:00:40.895Z - Connected to MySQL
2020-12-12T16:00:40.896Z - -=- jsonl <-> mysql(watt)
2020-12-12T16:04:04.369Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.0k/s, took 3m23.473s, rate ~ 155.0k/s count: 31536000
2020-12-12T16:07:26.899Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 155.7k/s, took 6m46.003s, rate ~ 155.3k/s count: 63072000
2020-12-12T16:10:48.511Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 156.4k/s, took 10m7.615s, rate ~ 155.7k/s count: 94608000
2020-12-12T16:14:09.700Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 156.7k/s, took 13m28.804s, rate ~ 156.0k/s count: 126144000
2020-12-12T16:17:31.560Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 156.2k/s, took 16m50.664s, rate ~ 156.0k/s count: 157680000
2020-12-12T16:19:50.132Z - jsonl <-> mysql(watt) took 19m9.237s, rate ~ 156.2k/s count: 179485128
2020-12-12T16:20:50.016Z - Verified jsonl <-> mysql(watt):
2020-12-12T16:20:50.016Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T16:20:50.016Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T16:20:50.016Z - [2009-11-17T03:00:30Z, 2014-08-06T04:19:27Z](140599047) Equal
2020-12-12T16:20:50.016Z - [2014-08-06T04:19:28Z, 2015-09-28T14:06:52Z](33252817) MissingInB
2020-12-12T16:20:50.016Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T16:24:25.925Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 146.1k/s, took 3m35.909s, rate ~ 146.1k/s count: 31536000
2020-12-12T16:27:43.964Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 159.2k/s, took 6m53.948s, rate ~ 152.4k/s count: 63072000
2020-12-12T16:31:01.606Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 159.6k/s, took 10m11.589s, rate ~ 154.7k/s count: 94608000
2020-12-12T16:34:17.981Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 160.6k/s, took 13m27.965s, rate ~ 156.1k/s count: 126144000
2020-12-12T16:37:34.224Z - jsonl <-> mysql(ted_native) (2014-03-26) inner rate ~ 160.7k/s, took 16m44.208s, rate ~ 157.0k/s count: 157680000
2020-12-12T16:38:43.575Z - jsonl <-> mysql(ted_native) took 17m53.558s, rate ~ 157.3k/s count: 168924060
2020-12-12T16:39:43.291Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T16:39:43.291Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T16:39:43.291Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T16:39:43.291Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T16:39:43.291Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T16:39:43.291Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T16:39:43.291Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T16:39:43.291Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T16:39:43.291Z - [2013-10-18T16:06:13Z, 2014-08-06T04:19:27Z](24759902) Equal
2020-12-12T16:39:43.291Z - [2014-08-06T04:19:28Z, 2015-09-28T14:06:52Z](33252817) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20150928.1006.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	179m23.209s
user	4m44.234s
sys	1m40.917s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2015-09-28 14:06:52	202176877
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2015-09-28 14:06:52	212737945

- run mysqlrestore
2020-12-12T19:39:12.348Z - Starting TED1K mysql restore
2020-12-12T19:39:12.350Z - Connected to MySQL
2020-12-12T19:39:12.350Z - -=- jsonl <-> mysql(watt)
2020-12-12T19:42:32.749Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 157.4k/s, took 3m20.398s, rate ~ 157.4k/s count: 31536000
2020-12-12T19:45:52.706Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.7k/s, took 6m40.354s, rate ~ 157.5k/s count: 63072000
2020-12-12T19:49:10.759Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 159.2k/s, took 9m58.407s, rate ~ 158.1k/s count: 94608000
2020-12-12T19:52:29.233Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 158.9k/s, took 13m16.881s, rate ~ 158.3k/s count: 126144000
2020-12-12T19:55:46.680Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 159.7k/s, took 16m34.329s, rate ~ 158.6k/s count: 157680000
2020-12-12T19:59:04.338Z - jsonl <-> mysql(watt) (2014-12-15) inner rate ~ 159.5k/s, took 19m51.987s, rate ~ 158.7k/s count: 189216000
2020-12-12T20:01:32.134Z - jsonl <-> mysql(watt) took 22m19.782s, rate ~ 158.8k/s count: 212737945
2020-12-12T20:01:32.134Z - Verified jsonl <-> mysql(watt):
2020-12-12T20:01:32.134Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T20:01:32.134Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T20:01:32.134Z - [2009-11-17T03:00:30Z, 2015-09-28T14:06:52Z](173851864) Equal
2020-12-12T20:01:32.134Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T20:05:13.716Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 142.3k/s, took 3m41.582s, rate ~ 142.3k/s count: 31536000
2020-12-12T20:08:35.481Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 156.3k/s, took 7m3.347s, rate ~ 149.0k/s count: 63072000
2020-12-12T20:11:57.441Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 156.1k/s, took 10m25.307s, rate ~ 151.3k/s count: 94608000
2020-12-12T20:15:19.891Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 155.8k/s, took 13m47.757s, rate ~ 152.4k/s count: 126144000
2020-12-12T20:18:41.993Z - jsonl <-> mysql(ted_native) (2014-03-26) inner rate ~ 156.0k/s, took 17m9.859s, rate ~ 153.1k/s count: 157680000
2020-12-12T20:22:03.403Z - jsonl <-> mysql(ted_native) (2015-04-27) inner rate ~ 156.6k/s, took 20m31.269s, rate ~ 153.7k/s count: 189216000
2020-12-12T20:23:26.247Z - jsonl <-> mysql(ted_native) took 21m54.113s, rate ~ 153.9k/s count: 202176877
2020-12-12T20:23:26.248Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T20:23:26.248Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T20:23:26.248Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T20:23:26.248Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T20:23:26.248Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T20:23:26.248Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T20:23:26.248Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T20:23:26.248Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T20:23:26.248Z - [2013-10-18T16:06:13Z, 2015-09-28T14:06:52Z](58012719) Equal
- Done Restoring Database

åreal	1479m2.641s
user	358m0.091s
sys	47m30.605s
```