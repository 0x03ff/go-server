#!/bin/bash

# DDoS Performance Testing Script using hey (supports HTTPS)
# Tests 5 configurations: HTTP+non-enc, HTTPS+non-enc, HTTPS+AES, HTTPS+RSA-2048, HTTPS+RSA-4096

USER_ID="265dabde-29f3-487f-9926-b7a593b13467"
RESULTS_DIR="./ddos_results"
HEY="$HOME/go/bin/hey"

# Check if hey is installed
if [ ! -f "$HEY" ]; then
    echo "Error: hey not found at $HEY"
    echo "Please run: go install github.com/rakyll/hey@latest"
    exit 1
fi

mkdir -p "$RESULTS_DIR"

echo "========================================"
echo "DDoS Performance Testing Started (using hey)"
echo "Time: $(date)"
echo "========================================"
echo ""

# Test parameters: -c 200 (200 workers), -z 30s (30 seconds)

# Test 1: HTTP + Non-encrypted (Baseline)
echo "[1/5] Testing HTTP + Non-encrypted (Baseline)..."
$HEY -c 200 -z 30s "http://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/1_http_nonenc.txt" 2>&1
echo "✓ Completed"
echo ""
sleep 5

# Test 2: HTTPS + Non-encrypted
echo "[2/5] Testing HTTPS + Non-encrypted..."
$HEY -c 200 -z 30s "https://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/2_https_nonenc.txt" 2>&1
echo "✓ Completed"
echo ""
sleep 5

# Test 3: HTTPS + AES (with decryption)
echo "[3/5] Testing HTTPS + AES..."
$HEY -c 200 -z 30s "https://localhost/api/download_folder/${USER_ID}/2?decrypt=true" > "${RESULTS_DIR}/3_https_aes.txt" 2>&1
echo "✓ Completed"
echo ""
sleep 5

# Test 4: HTTPS + RSA-2048 (with decryption)
echo "[4/5] Testing HTTPS + RSA-2048..."
$HEY -c 200 -z 30s "https://localhost/api/download_folder/${USER_ID}/3?decrypt=true" > "${RESULTS_DIR}/4_https_rsa2048.txt" 2>&1
echo "✓ Completed"
echo ""
sleep 5

# Test 5: HTTPS + RSA-4096 (with decryption)
echo "[5/5] Testing HTTPS + RSA-4096..."
$HEY -c 200 -z 30s "https://localhost/api/download_folder/${USER_ID}/5?decrypt=true" > "${RESULTS_DIR}/5_https_rsa4096.txt" 2>&1
echo "✓ Completed"
echo ""

# Generate summary report
echo "========================================"
echo "Generating Performance Summary..."
echo "========================================"
echo ""

{
    echo "DDoS Performance Test Results (using hey)"
    echo "========================================="
    echo "Test Date: $(date)"
    echo ""
    echo "Configuration Details:"
    echo "- Workers: 200"
    echo "- Duration: 30 seconds each"
    echo ""
    echo "----------------------------"
    
    for file in "${RESULTS_DIR}"/*.txt; do
        config=$(basename "$file" .txt)
        echo ""
        echo "[$config]"
        grep "Requests/sec:" "$file" || echo "N/A"
        grep "Average:" "$file" || echo "N/A"
        grep "Total:" "$file" | head -1 || echo "N/A"
    done
    
    echo ""
    echo "----------------------------"
    echo "Full results saved in: $RESULTS_DIR"
    
} > "${RESULTS_DIR}/summary_hey.txt"

cat "${RESULTS_DIR}/summary_hey.txt"

echo ""
echo "========================================"
echo "All Tests Completed!"
echo "Results saved in: $RESULTS_DIR"
echo "========================================"
