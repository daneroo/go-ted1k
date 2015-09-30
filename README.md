# MySQL test with Go

## Todo

- Cleanup docker script(s) : add start-db.sh
- Batch writes to ted/watt2


## Timing of MySQL reads
For timing of MySQL selects with maxCount results

From goedel to cantor

	3600: 989s
	3600*24: 357s
	3600*24*10: 324s

From Godel to local docker:

	3600*24: 294s,290s


## MySQL inserts are ridiculously slow
This is what we di to spead things up:
(all time measured on dirac/docker)

- Initial naive approach: 550 ins/s
- Prepared statements: 650 ins/s
- Transactions: 2500-3000 ins/s (depending on size of batch...)

