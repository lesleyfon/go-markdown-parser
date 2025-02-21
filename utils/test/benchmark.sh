#!/bin/bash

# Basic benchmark with readable units (-v for verbose)
go test ./ -bench=. -benchmem -v

# Time in milliseconds and memory in MB
go test ./ -bench=. -benchmem -benchtime=1s -v |
  awk '/Benchmark/ { \
        printf "%-40s    %8.2f ms    %8.2f MB    %8d allocs\n", \
        $1, \
        $3/1000000, \
        $5/1048576, \
        $7 \
    }'

# Run multiple times for more stable results
# go test ./ -bench=. -benchmem -count=5 -v

# Compare benchmarks side by side
# go test ./ -bench=. -benchmem -v |
#   benchstat /dev/stdin
