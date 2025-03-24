#!/bin/bash

# set -euxo 

# ./run-parallel.sh echo 'hello'
# ./run-parallel.sh -n 2 echo 'hello'

# Default number of instances
instances=10

# Parse command line arguments
while getopts "n:h" opt; do
    case $opt in
        n) instances=$OPTARG ;;
        h) echo "Usage: $0 [-n instances] binary [args...]"
           echo "  -n: number of instances to run (default: 10)"
           exit 0
           ;;
        \?) echo "Invalid option: -$OPTARG" >&2
            exit 1
            ;;
    esac
done

# Shift away the parsed options
shift $((OPTIND-1))

# Check if binary is provided
if [ $# -eq 0 ]; then
    echo "Error: No binary specified"
    echo "Usage: $0 [-n instances] binary [args...]"
    exit 1
fi

# Store binary and its arguments
binary="$@"

# Array to store PIDs
declare -a pids

echo "Starting $instances instances of: $binary"

# Start processes
for ((i=1; i<=$instances; i++)); do
    $binary &
    pids+=($!)
    echo "Started instance $i with PID ${pids[-1]}"
done

echo -e "\nAll instances started. PIDs: ${pids[@]}"

# Wait for all processes to complete
for pid in "${pids[@]}"; do
    wait $pid
    status=$?
    echo "Process $pid completed with status $status"
done

echo "All processes completed"
