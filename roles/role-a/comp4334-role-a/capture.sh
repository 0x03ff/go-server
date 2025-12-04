#!/bin/bash
# Capture CPU and memory profiles during attack
# Usage: ./capture.sh [duration_seconds] [interval_seconds]

DURATION=${1:-1500}  # 25 minutes in seconds
INTERVAL=${2:-30}   # 30 seconds
OUTPUT_DIR="profiles_$(date +%Y%m%d_%H%M)"
mkdir -p "$OUTPUT_DIR"

echo "Starting profile capture for $((DURATION/60)) minutes (every $((INTERVAL)) seconds)"
echo "Output directory: $OUTPUT_DIR"

START_TIME=$(date +%s)
END_TIME=$((START_TIME + DURATION))

while [ $(date +%s) -lt $END_TIME ]; do
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    
    # Capture CPU profile (10-second sample) using wget
    wget -q "http://localhost:8086/debug/pprof/profile?seconds=10" -O "$OUTPUT_DIR/cpu_$TIMESTAMP.pprof"
    
    # Capture heap profile using wget
    wget -q "http://localhost:8086/debug/pprof/heap" -O "$OUTPUT_DIR/heap_$TIMESTAMP.pprof"
    
    echo "Captured at $(date): cpu_$TIMESTAMP.pprof, heap_$TIMESTAMP.pprof"
    
    # Wait until next interval
    sleep $((INTERVAL - ( $(date +%s) - START_TIME ) % INTERVAL ))
done

echo "Capture complete. Profiles saved in $OUTPUT_DIR"

