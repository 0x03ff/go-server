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
