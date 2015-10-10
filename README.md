# MySQL test with Go

## Todo

- Batch writes to ted/watt2 - Reader or LOAD INFILE
- Refactor Progress (TimeTrack...) Interceptor channel ?

## Vendoring
See this [Go/Wiki for reference](https://github.com/golang/go/wiki/PackageManagementTools)

We want to use `GO15VENDOREXPERIMENT=1` and place our external dependencies in a `vendor folder`.
We are using [`govend`](https://github.com/gophersaurus/govend) as listed.

To install `govend`, we did a standard `go get -u github.com/gophersaurus/govend`, and makde sure our `GOPATH` was set and `$GOPATH/bin` is on our `$PATH`. also `GO15VENDOREXPERIMENT=1` needs to be set.

## InfluxDB

	select mean(value)*24/1000 from watt where time > '2008-01-01' and time < '2016-01-01' group by time(7d)

Continuous Queries are not appropriate for historical data loading.
I should implement my own select .. into (in go), using tablenames as in mysql

	select mean(value)*24/1000 into kwh_1d from watt where time > '2015-09-01' group by time(1d)

## Docker
We have abandoned data volumes for now.
`docker-compose` command brings up MySQL and InfluxDB instances, and the `restore` script restores a MySQL snapshot/

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

    3600*24: 294s,290s  (Read-only)
	10000: --s  (Batch Writes)



## MySQL inserts are ridiculously slow
This is what we di to spead things up:
(all time measured on dirac/docker)

- Initial naive approach: 550 ins/s
- Prepared statements: 650 ins/s
- Transactions: 2500-3000 ins/s (depending on size of batch...)

