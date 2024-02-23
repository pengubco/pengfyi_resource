#!/bin/bash

# Check if the user provided a command to run
if [ -z "$1" ]; then
	echo "Usage: $0 <command_to_run>"
	exit 1
fi

command_to_run="$1"

for ((i = 1; i <= 5; i++)); do
	$command_to_run
done
