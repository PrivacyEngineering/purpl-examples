#!/bin/bash

sum=0
count=0

for ((i=1; i<=1000; i++)); do
    output=$(go run clients.go trackingService-minimal)
    duration=$(echo "$output" | awk '{print $1}')
    echo "$duration" >> durations.txt
    sum=$((sum + duration))
    count=$((count + 1))
done

average=$((sum / count))
echo "Average duration: $average"
