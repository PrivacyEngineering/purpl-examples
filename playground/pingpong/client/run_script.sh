#!/bin/bash

for ((i=1; i<=10; i++)); do
    output=$(go run clients.go)
    duration=$(echo "$output" | grep -oE '[0-9]')
    echo "$duration" >> durations.txt
done
