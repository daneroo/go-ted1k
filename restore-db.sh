# This script creates or recreates a running container and restores a MySQL database snapshot into it.
# This script creates a volatile container, which will disappear when the container is destroyed.
#  Mapping an osx Volume, seems to create permissions problems
#  Use a data-volume

# Read the environment variables from environment specific .env file
. ./MYSQL.env

main() {

	echo "- Restoring database from snapshot: ${SNAPSHOTFILE}"
	echo "  To docker container named: ${CONTAINER_NAME}"
	echo "  Using database: ${MYSQL_DATABASE}, MYSQL_USER=${MYSQL_USER}"
	echo "  Data Volume will persisted inside that docker container"
	echo

	# Exit if docker is not setup in this shell
	checkDocker;

	waitForMysql ${CONTAINER_NAME};

	restoreSnapshot;
	
	echo "- Done Restoring Database"
	echo	
}

restoreSnapshot() {
	# command string for executing mysql command inside running container (with exec)
	# Notice no `-t' to pipe into docker
	mysqlExecCmd="docker exec -i ${CONTAINER_NAME} mysql -p${MYSQL_PASSWORD} -u${MYSQL_USER} ${MYSQL_DATABASE}"
	# echo CMD: $mysqlExecCmd

	# pump the snapshot into mysql (through docker exec)
	echo "- Restoring database..."
	# gzip --decompress --stdout  ${SNAPSHOTFILE} | $mysqlExecCmd
	bzcat  ${SNAPSHOTFILE} | $mysqlExecCmd
	echo

	echo "- Expect something recent in watt table"
	echo "select min(stamp),max(stamp),count(*) from watt" | $mysqlExecCmd
	echo
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