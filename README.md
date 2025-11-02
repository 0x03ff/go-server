# go-server

simple go server

Group 18: topic-"Network Infrastructure Hardening: Evaluating Against DDoS Attack on Web Servers

The web server use postgres with docker container, to start with:
sudo docker compose -f 'docker-compose.yml' up -d --build

for development golang

go build  -o ./bin/main ./cmd/ && ./bin/main

For use browser, You may need import the ca_cert.pem in browser setting.

Incase you want create you own cert for testing:

openssl ecparam -genkey -name secp384r1 -out go_key.pem
openssl req -new -key go_key.pem -out go_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"
openssl x509 -req -days 365 -in go_csr.pem -signkey go_key.pem -out go_cert.pem

openssl ecparam -genkey -name secp384r1 -out ca_key.pem

openssl req -x509 -new -key ca_key.pem -out ca_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-CA" -days 365

openssl x509 -req -in go_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 01 -out go_cert.pem -days 365
