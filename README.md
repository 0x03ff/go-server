# go-server

simple go server

Group 18: topic-"Network Infrastructure Hardening: Evaluating Against DDoS Attack on Web Servers

The web server use postgres with docker container, to start with:
sudo docker compose -f 'docker-compose.yml' up -d --build

if need reset database completely:

sudo docker compose -f 'docker-compose.yml' down
sudo docker volume rm go-server_postgres_data
sudo docker compose -f 'docker-compose.yml' up -d --build

for development golang

go build  -o ./bin/main ./cmd/ && sudo ./bin/main

For use browser, You may need import the ca_cert.pem in browser setting.

Incase you want create you own cert for testing:

openssl ecparam -genkey -name secp384r1 -out go_key.pem
openssl req -new -key go_key.pem -out go_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"
openssl x509 -req -days 365 -in go_csr.pem -signkey go_key.pem -out go_cert.pem

openssl ecparam -genkey -name secp384r1 -out ca_key.pem

openssl req -x509 -new -key ca_key.pem -out ca_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-CA" -days 365

openssl x509 -req -in go_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 01 -out go_cert.pem -days 365

## Role D - DDoS Performance Testing

### Prerequisites

Install hey tool for HTTP load testing:

```bash
go install github.com/rakyll/hey@latest
```

The hey tool will be installed in `~/go/bin/hey`

Install sysstat for CPU monitoring if it is not installed:

```bash
sudo apt install sysstat
```

### Prepare Test Data

Before running DDoS tests, upload test folders with different encryption types through the web interface:
- Folder 1: non-encrypted
- Folder 2: AES
- Folder 3: RSA-2048
- Folder 4: RSA-4096

Each folder should contain ~10MB files for meaningful performance comparison.

### Configure Test Script

Check your user ID and folder IDs from database:

```bash
sudo docker exec postgres_db psql -U comp4334 -d go_server -c "SELECT id, username FROM users;"
sudo docker exec postgres_db psql -U comp4334 -d go_server -c "SELECT id, title, encrypt FROM folders ORDER BY id;"
```

Edit `ddos_test_hey.sh` and update the `USER_ID` variable with your actual user ID:

```bash
USER_ID="your-user-id-here"
```

Optionally adjust test parameters:
- `CONCURRENCY`: Number of concurrent workers (default: 1700)
- `DURATION`: Test duration per scenario (default: 5s)

### Run DDoS Test

Run the test script:

```bash
chmod +x ddos_test_hey.sh
./ddos_test_hey.sh
```

To run fresh tests (clean previous results):

```bash
rm -rf ddos_results && ./ddos_test_hey.sh
```

### Test Scenarios

The script will automatically run 5 test scenarios:
1. HTTP + Non-encrypted (Baseline)
2. HTTPS + Non-encrypted
3. HTTPS + AES
4. HTTPS + RSA-2048
5. HTTPS + RSA-4096

Each test measures:
- Throughput (requests/sec)
- Average response time
- Fastest/Slowest response
- Total test duration
- System-wide CPU usage (%)

### View Results

Results are saved in `ddos_results/` folder:

- CPU monitoring uses `mpstat` to capture system-wide CPU utilization
- The script automatically detects CPU core count using `nproc`
- Server must be running before executing tests
- Tests run sequentially with 5-second cooldown between scenarios


