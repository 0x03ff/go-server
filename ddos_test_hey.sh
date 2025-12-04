#!/bin/bash

# DDoS Performance Testing Script using hey (supports HTTPS)
# Tests 5 configurations: HTTP+non-enc, HTTPS+non-enc, HTTPS+AES, HTTPS+RSA-2048, HTTPS+RSA-4096

USER_ID="a8f91811-99e8-43d0-8671-11c011a7af37"
RESULTS_DIR="./ddos_results"
HEY="$HOME/go/bin/hey"
CONCURRENCY=200
DURATION="30s"
DURATION_SEC=30

# Check if hey is installed
if [ ! -f "$HEY" ]; then
    echo "Error: hey not found at $HEY"
    echo "Please run: go install github.com/rakyll/hey@latest"
    exit 1
fi

mkdir -p "$RESULTS_DIR"

# Get number of CPU cores
NUM_CPUS=$(nproc)

# Get server PID (the actual ./bin/main process, not sudo)
SERVER_PID=$(pgrep -f "^\./bin/main$" || pgrep -x "main")
if [ -z "$SERVER_PID" ]; then
    echo "Error: Server not running (bin/main process not found)"
    exit 1
fi

echo "========================================"
echo "DDoS Performance Testing Started (using hey)"
echo "Time: $(date)"
echo "Server PID: $SERVER_PID"
echo "========================================"
echo ""

# Function to monitor CPU during test (system-wide)
monitor_cpu() {
    local test_name=$1
    local duration=$2
    # Monitor overall CPU usage using mpstat
    ( mpstat 1 $duration | grep -E "^[0-9]" | awk '{print 100 - $NF}' ) > "${RESULTS_DIR}/${test_name}_cpu.txt" 2>&1 &
    echo $!
}

# Test 1: HTTP + Non-encrypted (Baseline)
echo "[1/5] Testing HTTP + Non-encrypted (Baseline)..."
CPU_MON_PID=$(monitor_cpu "1_http_nonenc" $DURATION_SEC)
$HEY -c $CONCURRENCY -z $DURATION "http://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/1_http_nonenc.txt" 2>&1
wait $CPU_MON_PID 2>/dev/null
echo "✓ Completed"
echo ""
sleep 5

# Test 2: HTTPS + Non-encrypted
echo "[2/5] Testing HTTPS + Non-encrypted..."
CPU_MON_PID=$(monitor_cpu "2_https_nonenc" $DURATION_SEC)
$HEY -c $CONCURRENCY -z $DURATION "https://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/2_https_nonenc.txt" 2>&1
wait $CPU_MON_PID 2>/dev/null
echo "✓ Completed"
echo ""
sleep 5

# Test 3: HTTPS + AES (with decryption)
echo "[3/5] Testing HTTPS + AES..."
CPU_MON_PID=$(monitor_cpu "3_https_aes" $DURATION_SEC)
$HEY -c $CONCURRENCY -z $DURATION "https://localhost/api/download_folder/${USER_ID}/2?decrypt=true" > "${RESULTS_DIR}/3_https_aes.txt" 2>&1
wait $CPU_MON_PID 2>/dev/null
echo "✓ Completed"
echo ""
sleep 5

# Test 4: HTTPS + RSA-2048 (with decryption)
echo "[4/5] Testing HTTPS + RSA-2048..."
CPU_MON_PID=$(monitor_cpu "4_https_rsa2048" $DURATION_SEC)
$HEY -c $CONCURRENCY -z $DURATION "https://localhost/api/download_folder/${USER_ID}/3?decrypt=true" > "${RESULTS_DIR}/4_https_rsa2048.txt" 2>&1
wait $CPU_MON_PID 2>/dev/null
echo "✓ Completed"
echo ""
sleep 5

# Test 5: HTTPS + RSA-4096 (with decryption)
echo "[5/5] Testing HTTPS + RSA-4096..."
CPU_MON_PID=$(monitor_cpu "5_https_rsa4096" $DURATION_SEC)
$HEY -c $CONCURRENCY -z $DURATION "https://localhost/api/download_folder/${USER_ID}/4?decrypt=true" > "${RESULTS_DIR}/5_https_rsa4096.txt" 2>&1
wait $CPU_MON_PID 2>/dev/null
echo "✓ Completed"
echo ""

# Generate summary report
echo "========================================"
echo "Generating Performance Summary..."
echo "========================================"
echo ""

{
    echo "======================================================================"
    echo "          DDoS Performance Test Results - HTTP Flood Attack"
    echo "======================================================================"
    echo ""
    echo "Test Date: $(date)"
    echo "Test Configuration:"
    echo "  - Concurrent Workers: $CONCURRENCY"
    echo "  - Test Duration: $DURATION per scenario"
    echo "  - File Size: ~10MB (encrypted folders)"
    echo ""
    echo "======================================================================"
    echo ""
    
    # Process results in order
    for test_num in 1 2 3 4 5; do
        # Find the test file
        test_file=$(ls "${RESULTS_DIR}/${test_num}_"*.txt 2>/dev/null | grep -v "_cpu.txt" | head -1)
        [ -f "$test_file" ] || continue
        
        config=$(basename "$test_file" .txt)
        cpu_file="${RESULTS_DIR}/${config}_cpu.txt"
        
        # Extract test name
        case $config in
            1_http_nonenc) test_name="HTTP + Non-encrypted (Baseline)" ;;
            2_https_nonenc) test_name="HTTPS + Non-encrypted" ;;
            3_https_aes) test_name="HTTPS + AES" ;;
            4_https_rsa2048) test_name="HTTPS + RSA-2048" ;;
            5_https_rsa4096) test_name="HTTPS + RSA-4096" ;;
            *) test_name="Unknown Test" ;;
        esac
        
        echo "Test ${test_num}: ${test_name}"
        echo "----------------------------------------------------------------------"
        
        # Extract metrics
        req_sec=$(grep "Requests/sec:" "$test_file" | awk '{print $2}')
        avg_time=$(grep "Average:" "$test_file" | awk '{print $2, $3}')
        total_time=$(grep "Total:" "$test_file" | head -1 | awk '{print $2, $3}')
        slowest=$(grep "Slowest:" "$test_file" | awk '{print $2, $3}')
        fastest=$(grep "Fastest:" "$test_file" | awk '{print $2, $3}')
        
        # Calculate CPU usage (system-wide, already in percentage)
        if [ -f "$cpu_file" ]; then
            avg_cpu=$(awk '{sum+=$1; count++} END {if(count>0) printf "%.2f", sum/count}' "$cpu_file")
        else
            avg_cpu="N/A"
        fi
        
        # Display results
        printf "  %-25s %s\n" "Throughput:" "${req_sec:-N/A} req/sec"
        printf "  %-25s %s\n" "Average Response Time:" "${avg_time:-N/A}"
        printf "  %-25s %s\n" "Fastest Response:" "${fastest:-N/A}"
        printf "  %-25s %s\n" "Slowest Response:" "${slowest:-N/A}"
        printf "  %-25s %s\n" "Total Test Duration:" "${total_time:-N/A}"
        printf "  %-25s %s%%\n" "Average CPU Usage:" "${avg_cpu}"
        echo ""
    done
    
    echo "======================================================================"
    echo "Detailed results: ${RESULTS_DIR}/*.txt"
    echo "CPU monitoring data: ${RESULTS_DIR}/*_cpu.txt"
    echo "======================================================================"
    
} > "${RESULTS_DIR}/summary_hey.txt"

cat "${RESULTS_DIR}/summary_hey.txt"

echo ""
echo "========================================"
echo "All Tests Completed!"
echo "Results saved in: $RESULTS_DIR"
echo "========================================"
