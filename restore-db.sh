# This script creates or recreates a running container and restores a MySQL database snapshot into it.
# This script creates a volatile container, which will disappear when the container is destroyed.
#  Mapping an osx Volume, seems to create permissions problems
#  Use a data-volume

DBIMAGE=mysql	
CONTAINER_NAME=teddb
CONTAINER_NAME_DATA=teddata
SNAPSHOTFILE=~/Downloads/ted/ted.watt.20150928.1003.sql.bz2
DB_NAME=ted
DB_USER=ted
DB_PASS=secret

# allow passwordless root
# MYSQL_ROOT_PASSWORD=
# MYSQL_ALLOW_EMPTY_PASSWORD=yes

# No persistence for now (permission problems)
# TODO: How about we create a path outside osx's default share: /data on host
# DB_PATH_ON_DOCKER_HOST=`pwd`/data
# DB_PATH_ON_DOCKER_HOST=/data/mysql

restoreDB() {

	echo "- Restoring database from snapshot: ${SNAPSHOTFILE}"
	echo "  Using DB_NAME: ${DB_NAME}, DB_USER=${DB_USER}"
	echo "  Data will be stored in a docker container named: ${CONTAINER_NAME}"
	echo "  Data Volume will persisted on host in data volume: ${CONTAINER_NAME_DATA}"
	echo

	# Exit if docker is not setup in this shell
	checkDocker;

	# createDataVolume - move to function
	# Should I force delete: No : instruct on how to remove (and therefore recreate data-volume)
	echo "  Creating data volume: ${CONTAINER_NAME_DATA}"
	docker create -v /var/lib/mysql --name ${CONTAINER_NAME_DATA} ${DBIMAGE} /bin/true


	# remove any existing container with name ${CONTAINER_NAME}
	echo "- Remove the database container if it exists. (${CONTAINER_NAME})"
	docker kill ${CONTAINER_NAME} 2>&1 >/dev/null
	docker rm ${CONTAINER_NAME} 2>&1 >/dev/null

	echo "- Start the database container ${CONTAINER_NAME}"
	# 
	cid=$(docker run -d --name ${CONTAINER_NAME} -p 3306:3306 --volumes-from ${CONTAINER_NAME_DATA} -e MYSQL_USER=${DB_USER} -e MYSQL_PASSWORD=${DB_PASS} -e MYSQL_DATABASE=${DB_NAME} -e MYSQL_ALLOW_EMPTY_PASSWORD=yes ${DBIMAGE})

	waitForMysql $cid;

	# command string for executing mysql command inside running container (with exec)
	# Notice no `-t' to pipe into docker
	mysqlExecCmd="docker exec -i ${cid} mysql -p${DB_PASS} -u${DB_USER} ${DB_NAME}"
	# echo CMD: $mysqlExecCmd

	# pump the snapshot into mysql (through docker exec)
	echo "- Restoring database..."
	# gzip --decompress --stdout  ${SNAPSHOTFILE} | $mysqlExecCmd
	bzcat  ${SNAPSHOTFILE} | $mysqlExecCmd
	echo

	echo "- Expect something recent in Gps table"
	echo "select min(stamp),max(stamp),count(*) from watt" | $mysqlExecCmd
	echo

	echo "- Done Restoring Database"
	echo	
}

# Exit if docker info return an error
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

# Wait for asuccesful connection (using credentials , DB_USER, DB_PASS)
waitForMysql() {
	local containerId=$1
	timeout=30
	# wait for mysql server to start (max $timeout seconds)
	echo -n "- Waiting for database server to accept connections (max $timeout seconds)"
	while ! docker exec -i ${containerId} mysqladmin -p${DB_PASS} -u${DB_USER} status >/dev/null 2>&1; do
		timeout=$(($timeout - 1))
		if [ $timeout -eq 0 ]; then	
			echo	
			echo "  ERROR: Could not connect to database. Aborting..."
			echo "  This could mean that the credentials do not match a previously initialiazed database"

			# clean up
			cleanUp $containerId;

			exit 1
		fi
		echo -n "."


		sleep 1
	done
	echo
	echo "  Connected"
}

# Stop and remove the container (parameter)
cleanUp() {
	local containerId=$1
	echo "- Cleaning up ($containerId)"
	echo "  Stopping the database container..."
	# docker stop ${containerId} >/dev/null

	echo "  Removing the database container..."
	# docker rm ${containerId} >/dev/null
}

restoreDB;