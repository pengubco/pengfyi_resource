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

PROGRAM=../build/bin/mpg_example

if [ ! -e "$PROGRAM" ]; then
    echo "$PROGRAM does not exist. build it first"
		exit 1
fi

run_program() {
    GOMAXPROCS=$1 $PROGRAM -quiet -numberOfGoroutines $2 &
		PROGRAM_PID=$!
		sleep 3
		threadCnt=$(gops $(pidof mpg_example) | grep threads |  awk '{print $2}')
		kill $PROGRAM_PID
		echo $1,$(($2+1)),$threadCnt
		sleep 3
}

echo "GOMAXPROCS,user-level goroutines,os threads"

# for gomaxpros in {1..9..1} {10..19..5} {20..100..10}; do
for gomaxprocs in {1..100..1}; do
	run_program $gomaxprocs 10000
done

