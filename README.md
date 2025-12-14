# go-server

simple go server

Group 18: topic-"Network Infrastructure Hardening: Evaluating Against DDoS Attack on Web Servers

To start with, the enviriroment require the golang, docker, and vscode as more convenient

```
sudo apt update
sudo apt install snapd -y
sudo apt install git -y
sudo snap install go --classic
sudo snap install code --classic
sudo snap install docker
```

The web server use postgres with docker container, to start with:

```bash
sudo docker compose -f 'docker-compose.yml' up -d --build
```

if need reset database completely:

```bash
sudo docker compose -f 'docker-compose.yml' down
sudo docker volume rm go-server_postgres_data
sudo docker compose -f 'docker-compose.yml' up -d --build
```

## Development

### Normal Mode (Full authentication, HTTP + HTTPS):

```bash
go build -o ./bin/main ./cmd/ && sudo ./bin/main
```

### Role D Testing Mode (Authentication disabled for folder downloads):

```bash
go build -o ./bin/main ./cmd/ && sudo ROLE_D_MODE=true ./bin/main
```

### Alternative (without sudo, using setcap):

```bash
go build -o ./bin/main ./cmd/ && sudo setcap 'cap_net_bind_service=+ep' ./bin/main && ./bin/main
```

## SSL/TLS Certificate Setup

For use browser, You may need import the ca_cert.pem in browser setting.

Incase you want create you own cert for testing:

```bash
openssl ecparam -genkey -name secp384r1 -out go_key.pem
openssl req -new -key go_key.pem -out go_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"
openssl x509 -req -days 365 -in go_csr.pem -signkey go_key.pem -out go_cert.pem

openssl ecparam -genkey -name secp384r1 -out ca_key.pem

openssl req -x509 -new -key ca_key.pem -out ca_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-CA" -days 365

openssl x509 -req -in go_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 01 -out go_cert.pem -days 365
```

---

## Default flag in ./cmd/main.go


That exists two flag:

drop_flag in line 51 (Default is True): if true, drop the entire database and preform database initialization.

(Each time restart the web server, all the store file, record, and user data is removed. )

If false, the database maintains the all the store file, record, and user data if any.

random_request_addressin line 55 (Fefault is True): If true, allows simulated different ip address on the rate-limiting machinal aspect.

(The effect scope only within the /json_handler/ login_handler , register_handler and methed_helper)

If false, allow trueip address on the rate-limiting machinal aspect, as visit localhost with given 127.0.0.1.

---



## Login page error: fail to generate token

**Stop the web server**

Go to the :

📦cmd
 ┣ 📂api
......
 ┗ 📜main.go <--

**Manually set the drop_flag to true**

Restart the web server, due to the database do not having the key that insert into database.

Such that drop the database and restart the initialization database process.

After that

**Manually set the **drop_flag** to false**

---

## Role A - Web Server Hardening Baseline

### web server setting

Go to the :

📦cmd
 ┣ 📂api
......
 ┗ 📜main.go <--

Manually set the random_request_address to true, to enable the simulate feature

### Prerequisites

Install OWASP zaproxy for the fuzz:

```
sudo apt update
sudo apt install snapd-y
sudo snap install zaproxy --classic
```

### Prepare Test Data

Before running DDoS tests, create the account with user name:

user name:

`comp4334`

user password:

`comp4334password`

user recover:

`comp4334recover`

Testing with the login page to verity can login.

### Simulation Case

go to the:

📦cmd
 ┣ 📂api
 ┃ ┣ 📂config
 ┃ ┃ ┗ 📜config.go
 *┃ ┗ 📂router*
 ┃ ┃ ┣ 📂html_handler
.....
 ┃ ┃ ┣ 📂json_handler
 ┃ ┃ ┃ ┣ 📂web_server
 ┃ ┃ ┃ ┃ ┣ 📂server
 ┃ ┃ ┃ ┃ ┃ ┗ 📜server_handler.go
 *┃ ┃ ┃ ┃ ┗ 📜server_router.go
*....**
 ┃ ┃ ┃ ┣ 📜method_helper.go <----
....

** ┣ 📜.DS_Store
 ┗ 📜main.go

uncomment the statement under the case comment

For instance, the case 1 is used:

```
func generateRandomIP() string {
	// 1. Legitimate user with exactly 10 devices (192.168.0.2 - 192.168.0.11)
	return fmt.Sprintf("192.168.0.%d", 2 + rand.Intn(10))

	// 2. Attacker with exactly 50 devices (192.168.0.2 - 192.168.0.51)
	// return fmt.Sprintf("192.168.0.%d", 2 + rand.Intn(50))
	......
}
```

To modify to case 2:

```
func generateRandomIP() string {
	// 1. Legitimate user with exactly 10 devices (192.168.0.2 - 192.168.0.11)
	// return fmt.Sprintf("192.168.0.%d", 2 + rand.Intn(10))

	// 2. Attacker with exactly 50 devices (192.168.0.2 - 192.168.0.51)
	return fmt.Sprintf("192.168.0.%d", 2 + rand.Intn(50))
	......
}
```

Such that simulate different ip address on the rate-limit aspect.

### Capture setup

The script is locate on role/role-a:

```
📦role-a
 ┣ 📂capture
 ┃ ┣ 📜capture.sh <----
....
```

1. Change directory to the folder capture
2. ./capture.sh [total time second] [capture preiod second], where default setting is: 1800 (30 min) and 30 second
3. start capture the web server pprof data (After start the web server)

   If that exist the error, consider:

   ```
   chmod +x ./capture.sh
   ```

To analyze the result:

Download the requirement:

```
sudo apt update
sudo apt install graphviz -y 

go install github.com/google/pprof@latest 
```

1. Change directory to the folder capture
2. Change directory to the create folder name call: profiles_2025.....
3. Start analysis the cpu (CPU usage)/ heap (memort usage) with web

   ```
   go tool pprof -http=:6060 cpu_[the name of the file]*.pprof 
   ```

### ZAP environment setup

1. Open the zaproxy
2. Import the session, File -> open the session , choose the session

   ```
   📦roles
    ┣ 📂role-a
   ...
    ┃ ┗ 📂zap-session
    ┃ ┃ ┣ 📂Session.session.tmp
    ┃ ┃ ┣ 📜.DS_Store
    ┃ ┃ ┣ 📜Session.session <----
    ┃ ┃ ┣ 📜Session.session.data
    ┃ ┃ ┣ 📜Session.session.lck
    ┃ ┃ ┣ 📜Session.session.log
    ┃ ┃ ┣ 📜Session.session.properties
    ┃ ┃ ┗ 📜Session.session.script
   ....
   ```
3. After import the session, choose SITES:HTTPS://127.0.0.1 -> api -> POST:Login(){"username":.....}
4. Choose attack -> fuzz
5. Highlight the json payload
6. Click Add.... -> Add...
7. Choose import type: file, import the fuzz payload file (payload_generater.py to generate the payload)

   ```
   📦roles
    ┣ 📂role-a
   ...
    ┃ ┣ 📂files
    ┃ ┃ ┣ 📜200000-fuzz_payloads.txt <---
    ┃ ┃ ┗ 📜payload_generater.py
   ......
   ```
8. Close the Windows
9. Click choose Start Fuzzer(After active the capture.sh and web server)

   To analysis the result:

   Watch the ZAP fuzz result:

   or

   watch the web server log

   ```
   📦logs <--
    ┣ 📜security_events_0001.csv
    ┗ 📜system.log
   ```

That have instruction with image on report.

---

## Role B: firewall

## Prerequisites

This assumes the server is already set up. The attacker and server should be different devices (VM-to-VM or host-to-VM), though attacking and hosting on the same VM is acceptable for testing.

## GUFW install

```bash
sudo add-apt-repository universe   # required for GUFW to function
sudo apt update -y
sudo apt install gufw -y
```

## GUFW set up

- Set GUFW logging level to medium:
  ```bash
  sudo ufw logging medium
  ```
- Import rules: Open GUFW UI → File → Import profile → select the profile to import.
- Export rules: Open GUFW UI → File → Export this profile → choose destination.
- Verify UFW status after applying the profile:
  ```bash
  sudo ufw status verbose
  ```

## Test instructions

### Prerequisites

Install hping3 (attack tool):

```bash
sudo apt install hping3
```

### Before test

- (Optional) Clear UFW log before testing.

```bash
cat /dev/null | sudo tee /var/log/ufw.log
```

### Test

1) Enable/disable UFW: open GUFW GUI and toggle `Status`.
2) Live monitor CPU usage & logs from server:

```bash
sudo top
```

3) Run attack from attacker:
   Note: replace `<HOST IP>` with your server VM's IP address. The `--flood` flag sends high-rate traffic—use only on isolated lab networks you control.

```bash
sudo hping3 <HOST IP> -S --flood -p 443
```

4) Watch logs or CPU usage or try accessing website in a private browser.
5) Export logs for reference:

   Note: Replace `<FILTER WORD>` with the log keyword you want to match (see Logging section for some of the filter words); and replace `/path/to/exported/log` with your desired output path/filename.

```bash
sudo cat /var/log/ufw.log > /path/to/exported/log
sudo grep <FILTER WORD> /var/log/ufw.log > /path/to/exported/log
```

6) Repeat with UFW enabled.

## Logging

Filter keywords (case sensitive):

- `UFW LIMIT BLOCK` / `UFW BLOCK`: blocked by UFW rules
- `UFW LIMIT ACCEPT` / `UFW ACCEPT`: not blocked by UFW rules
- `PROTO=TCP`: TCP packets
- `SYN`: SYN packets
- `SRC=<ip addr>`: packet source IP
- `DPT=443`: destination port 443
- `IN=` / `OUT=`: inbound/outbound interface labels

Example combined filter:

```bash
sudo grep "UFW BLOCK" /path/to/exported/log | grep "DPT=443"
```

---

---

# Role C: TCP Stack Settings

## Objective

Improve the server's resilience against DDoS attacks (such as TLS Handshake Flood) and optimize its performance under stress.

---

## Environment

Ideally, attacker and server should be different devices (VM-to-VM or host-to-VM), but same VM is acceptable.

---

## Install Required Tools

### Linux Utilities

```bash
sudo apt update
sudo apt upgrade -y
sudo apt install net-tools iproute2 htop curl wget -y
```

### Attack Simulation Tools

- **hping3** (for SYN flood):

```bash
sudo apt install hping3 -y
```

- **hey** (for HTTP/TLS stress test):

```bash
sudo apt install golang-go -y
go install github.com/rakyll/hey@latest
# Ensure hey is in ~/go/bin/hey
```

### Monitoring Tools

```bash
sudo apt install iftop iotop -y
```

---

## TCP Tuning Steps

### Step 1: Start VM

Login via terminal or SSH.

### Step 2: Enable SYN Cookies

```bash
sudo sysctl -w net.ipv4.tcp_syncookies=1
```

### Step 3: Increase Backlog Queue Size

```bash
sudo sysctl -w net.ipv4.tcp_max_syn_backlog=4096
```

### Step 4: Increase Connection Limits

```bash
sudo sysctl -w net.core.somaxconn=1024
```

### Step 5: Reduce FIN Timeout

```bash
sudo sysctl -w net.ipv4.tcp_fin_timeout=15
```

### Step 6: Adjust Buffer Sizes

```bash
sudo sysctl -w net.core.rmem_max=16777216
sudo sysctl -w net.core.wmem_max=16777216
```

---

## Make Changes Persistent

Edit `/etc/sysctl.conf` and add:

```
net.ipv4.tcp_syncookies=1
net.ipv4.tcp_max_syn_backlog=4096
net.core.somaxconn=1024
net.ipv4.tcp_fin_timeout=15
net.core.rmem_max=16777216
net.core.wmem_max=16777216
```

Apply changes:

```bash
sudo sysctl -p
```

---

## Testing

### Attack Simulation

- **SYN Flood**:

```bash
sudo hping3 -S -p 443 --flood <your IP>
```

- **HTTP/TLS Stress Test**:

```bash
hey -n 10000 -c 200 https://<your IP>
```

### Monitoring

```bash
netstat -s    # Connection state
ss -s         # Connection state
top / htop    # CPU usage
```

---

## Expected Results

- `hping3` and `hey` do not completely take down the server.
- `netstat` or `ss` shows fewer stuck connections.
- `htop` shows CPU usage is high but not maxed out during stress tests.

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

Edit \`ddos_test_hey.sh\` and update the \`USER_ID\` variable with your actual user ID:

```bash
USER_ID="your-user-id-here"
```

Optionally adjust test parameters:

- \`CONCURRENCY\`: Number of concurrent workers (default: 200)
- \`DURATION\`: Test duration per scenario (default: 30s)

### Run DDoS Test

Start the server in Role D testing mode:

```bash
go build -o ./bin/main ./cmd/ && sudo ROLE_D_MODE=true ./bin/main
```

In another terminal, run the test script:

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

Results are saved in \`ddos_results/\` folder:

- CPU monitoring uses \`mpstat\` to capture system-wide CPU utilization
- The script automatically detects CPU core count using \`nproc\`
- Server must be running before executing tests
- Tests run sequentially with 5-second cooldown between scenarios
