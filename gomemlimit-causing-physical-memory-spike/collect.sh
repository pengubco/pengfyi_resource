#!/bin/bash

# Given a PID and a duration, collect VmSize and VmRSS of the process and print to stdout every second.
# The output is in CSV format. 
if [ $# -ne 2 ]; then
    echo "Usage: $0 <PID> <duration_in_seconds>"
    exit 1
fi

PID=$1
DURATION=$2

# Check if process exists
if ! kill -0 "$PID" 2>/dev/null; then
    echo "Process $PID does not exist"
    exit 1
fi

echo "seconds,virtual_memory_kb,physical_memory_kb" 

counter=0

while [ $counter -lt "$DURATION" ]; do
    virtual=$(grep VmSize /proc/"$PID"/status | awk '{print $2}')
    physical=$(grep VmRSS /proc/"$PID"/status | awk '{print $2}')
    
    if [ -z "$virtual" ] || [ -z "$physical" ]; then
        echo "Process $PID has terminated"
        exit 1
    fi
    
    echo "$counter,$virtual,$physical" 
    
    counter=$((counter + 1))
    sleep 1
done
