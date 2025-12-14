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
