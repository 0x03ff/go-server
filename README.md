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

for role D, need install hey tool for DDoS testing:

go install github.com/rakyll/hey@latest

the hey tool will be installed in ~/go/bin/hey

before run DDoS test, upload test folders with different encryption (non-encrypted, aes, rsa-2048, rsa-4096)

check your user id and folder ids from database:

sudo docker exec postgres_db psql -U comp4334 -d go_server -c "SELECT id FROM users;"
sudo docker exec postgres_db psql -U comp4334 -d go_server -c "SELECT id, title, encrypt FROM folders ORDER BY id;"

edit ddos_test_hey.sh and update USER_ID and folder ids

run DDoS test:

./ddos_test_hey.sh

test will run 5 scenarios (200 workers, 30 seconds each):
1. HTTP + non-encrypted
2. HTTPS + non-encrypted  
3. HTTPS + AES
4. HTTPS + RSA-2048
5. HTTPS + RSA-4096

results saved in ddos_results/ folder


