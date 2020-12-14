# This script creates or recreates a running container and restores a MySQL database snapshot into it.
# This script creates a volatile container, which will disappear when the container is destroyed.
#  Mapping an osx Volume, seems to create permissions problems
#  Use a data-volume

# Read the environment variables from environment specific .env file
. ./MYSQL.env

CONTAINER_NAME=go-ted1k_mysql_1
# for tiumeing format
export TIMEFORMAT="%Rs"

# turn this into a loop:
main() {

	# echo "- Restoring database from snapshot: ${SNAPSHOTFILE}"
	# echo "  To docker container named: ${CONTAINER_NAME}"
	# echo "  Using database: ${MYSQL_DATABASE}, MYSQL_USER=${MYSQL_USER}"
	# echo "  Data Volume will persisted inside that docker container"
	# echo

	# Exit if docker is not setup in this shell
	checkDocker;
	waitForMysql ${CONTAINER_NAME};

	# SNAPSHOTFILE=./data/archive/mirror/ted/ted.20090214.1756.sql.bz2 # first ted.20*.sql.bz2
	# SNAPSHOTFILE=./data/archive/mirror/ted/ted.20150928.1006.sql.bz2 # last ted.20*.sql.bz2
	# Phase-1
	# for SNAPSHOTFILE in ./data/archive/mirror/ted/ted.20??????.????.sql.bz2; do
	# pre-Phase-2
	# - ted.watt.2016-02-14-1555.sql.bz2 - last full fromt 2008-07-30
	# - ted.watt-just2016.2016-02-14-1624.sql.bz2 first partial from 2016-01-01
	# - ted.watt.20201120.2332Z.sql.bz2 most recent partial from 2016-01-01
	# Phase-2 verification
	for SNAPSHOTFILE in ./data/archive/mirror/ted/ted.watt*.sql.bz2; do
		echo
		echo "-=-= Restoring database from snapshot: ${SNAPSHOTFILE}"
		restoreSnapshot;
		echo "- Done Restoring Database"
	done
	
	echo	
}

restoreSnapshot() {
	# command string for executing mysql command inside running container (with exec)
	# Notice no `-t' to pipe into docker
	# mysqlExecCmd="docker exec -i ${CONTAINER_NAME} mysql -p${MYSQL_PASSWORD} -u${MYSQL_USER} ${MYSQL_DATABASE}"
	mysqlExecCmd="docker exec -i ${CONTAINER_NAME} mysql ${MYSQL_DATABASE}"
	# echo CMD: $mysqlExecCmd

	echo "- Drop tables watt and ted_native before restore, if present"
	echo 'drop table if exists ted_native;' | $mysqlExecCmd
	echo 'drop table if exists watt;'| $mysqlExecCmd

	echo "- Show remaining tables, before restore"
	echo 'show tables;' | $mysqlExecCmd
	# pump the snapshot into mysql (through docker exec)
	echo "- Restoring database..."
	# gzip --decompress --stdout  ${SNAPSHOTFILE} | $mysqlExecCmd
	time bzcat  ${SNAPSHOTFILE} | $mysqlExecCmd
	echo

	echo "- Expect something recent in watt table"
	echo "select min(stamp),max(stamp),count(*) from watt" | $mysqlExecCmd
	echo
	echo "- run mysql restore"
	go run cmd/mysqlrestore/mysqlrestore.go
}

# Wait for asuccesful connection (using credentials , MYSQL_USER, MYSQL_PASSWORD)
waitForMysql() {
	local timeout=30
	# wait for mysql server to start (max $timeout seconds)
	echo -n "- Waiting for database server (${CONTAINER_NAME}) to accept connections (max $timeout seconds)"
	while ! docker exec -i ${CONTAINER_NAME} mysqladmin -p${MYSQL_PASSWORD} -u${MYSQL_USER} status >/dev/null 2>&1; do
		timeout=$(($timeout - 1))
		if [ $timeout -eq 0 ]; then	
			echo	
			echo "  ERROR: Could not connect to database. Aborting..."
			echo "  This could mean that the credentials do not match a previously initialiazed database"

			exit 1
		fi
		echo -n "."
		sleep 1
	done
	echo
	echo "  Connected"
}

# Exit if docker info returns an error
checkDocker() {
	echo "- Verifying docker environment"
	# Check that docker is setup (docker info haz zero return code)
	docker info >/dev/null
	local dockerOK=$?
	if [ ${dockerOK} -gt 0 ]; then
		echo "  Docker does not seem to be setup properly"
		exit -1	
	else
		echo "  Docker seems to be setup properly"
		echo
	fi
}

main;