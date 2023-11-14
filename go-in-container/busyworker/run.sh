#!/bin/bash

# make sure it's using bash and the version is >= 5.
validate_shell() {
	echo "bash version: " $BASH_VERSION

	if [ -z "$BASH_VERSION" ]; then
			echo "This script requires Bash to run. Exiting."
			exit 1
	fi

	major_version=${BASH_VERSION:0:1}

	if [ "$major_version" -lt 5 ]; then
			echo "This script requires Bash version greater than 5. Exiting."
			exit 1
	fi
}

validate_shell

trap "exit" INT TERM
trap "kill 0" EXIT

CPUS=0.2
BUSY_WORKERS=100
BUSY_WORK_SIZE=1000000

echo "cpus=${CPUS} numberOfBusyWorkers=${BUSY_WORKERS} busyWorkSize=${BUSY_WORK_SIZE}"

for gomaxprocs in {1..10}; do
	echo "GOMAXPROCS=$gomaxprocs"
	for i in {1..3}; do
		docker run --rm --cpus=${CPUS} -e GOMAXPROCS=$gomaxprocs \
		peng.fyi/gomaxprocs_experiment /bin/busyworker \
		-numberOfBusyWorkers=${BUSY_WORKERS} \
		-busyWorkSize=${BUSY_WORK_SIZE} \
		-quiet=true \
		&
	done
	wait
	sleep 10
done
