#!/bin/bash

# DDoS Performance Testing Script
# Tests 5 configurations: HTTP+non-enc, HTTPS+non-enc, HTTPS+AES, HTTPS+RSA-2048, HTTPS+RSA-4096

USER_ID="39554041-5769-4cb4-89f1-1adb8a70df55"
RESULTS_DIR="./ddos_results"
mkdir -p "$RESULTS_DIR"

echo "========================================"
echo "DDoS Performance Testing Started"
echo "Time: $(date)"
echo "========================================"
echo ""

# Test 1: HTTP + Non-encrypted (Baseline)
echo "[1/5] Testing HTTP + Non-encrypted (Baseline)..."
wrk -t12 -c400 -d30s "http://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/1_http_nonenc.txt"
echo "✓ Completed"
echo ""
sleep 5

# Test 2: HTTPS + Non-encrypted
echo "[2/5] Testing HTTPS + Non-encrypted..."
wrk -t12 -c400 -d30s "https://localhost/api/download_folder/${USER_ID}/1" > "${RESULTS_DIR}/2_https_nonenc.txt"
echo "✓ Completed"
echo ""
sleep 5

# Test 3: HTTPS + AES (with decryption)
echo "[3/5] Testing HTTPS + AES..."
wrk -t12 -c400 -d30s "https://localhost/api/download_folder/${USER_ID}/2?decrypt=true" > "${RESULTS_DIR}/3_https_aes.txt"
echo "✓ Completed"
echo ""
sleep 5

# Test 4: HTTPS + RSA-2048 (with decryption)
echo "[4/5] Testing HTTPS + RSA-2048..."
wrk -t12 -c400 -d30s "https://localhost/api/download_folder/${USER_ID}/3?decrypt=true" > "${RESULTS_DIR}/4_https_rsa2048.txt"
echo "✓ Completed"
echo ""
sleep 5

# Test 5: HTTPS + RSA-4096 (with decryption)
echo "[5/5] Testing HTTPS + RSA-4096..."
wrk -t12 -c400 -d30s "https://localhost/api/download_folder/${USER_ID}/4?decrypt=true" > "${RESULTS_DIR}/5_https_rsa4096.txt"
echo "✓ Completed"
echo ""

# Generate summary report
echo "========================================"
echo "Generating Performance Summary..."
echo "========================================"
echo ""

{
    echo "DDoS Performance Test Results"
    echo "============================="
    echo "Test Date: $(date)"
    echo ""
    echo "Configuration Details:"
    echo "- Threads: 12"
    echo "- Connections: 400"
    echo "- Duration: 30 seconds each"
    echo ""
    echo "----------------------------"
    
    for file in "${RESULTS_DIR}"/*.txt; do
        config=$(basename "$file" .txt)
        echo ""
        echo "[$config]"
        grep "Requests/sec:" "$file" || echo "N/A"
        grep "Transfer/sec:" "$file" || echo "N/A"
        grep "Latency" "$file" | head -1 || echo "N/A"
    done
    
    echo ""
    echo "----------------------------"
    echo "Full results saved in: $RESULTS_DIR"
    
} > "${RESULTS_DIR}/summary.txt"

cat "${RESULTS_DIR}/summary.txt"

echo ""
echo "========================================"
echo "All Tests Completed!"
echo "Results saved in: $RESULTS_DIR"
echo "========================================"
