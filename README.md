# MySQL test with Go

## Todo

[See Evernote](https://www.evernote.com/shard/s60/nl/1773032759/ae1b9921-7e85-4b75-a21b-86be7d524295/)

## Vendoring

[Usage](https://github.com/govend/govend):

	govend -v  # download all the dependencies in the vendor.yml file
	govend -v -u # scan your project, update all dependencies, and update the vendor.yml revision versions

To install [`govend`](https://github.com/gophersaurus/govend) itself:

	go get -u github.com/govend/govend

and made sure our `GOPATH` was set and `$GOPATH/bin` is on our `$PATH`.

See this general vendoring entry: [Go/Wiki for reference](https://github.com/golang/go/wiki/PackageManagementTools).
Prior to `go1.6`, we also had to set `GO15VENDOREXPERIMENT=1`.

## InfluxDB

	select mean(value)*24/1000 from watt where time > '2008-01-01' and time < '2016-01-01' group by time(7d)

### Downsampled time series

Truncate for D,M,Y: http://play.golang.org/p/PUNNHq9sh6

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

    3600*24: 294s,290s  (Read-only) Now 412s,405s, with IgnoreAll
	10000: --s  (Batch Writes)



## MySQL inserts are ridiculously slow
This is what we di to spead things up:
(all time measured on dirac/docker)

- Initial naive approach: 550 ins/s
- Prepared statements: 650 ins/s
- Transactions: 2500-3000 ins/s (depending on size of batch...)

