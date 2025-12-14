# Role B: firewall
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
