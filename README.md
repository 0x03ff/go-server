# go-server

simple go server

Group 18: topic-"Network Infrastructure Hardening: Evaluating Against DDoS Attack on Web Servers

The web server use postgres with docker container, to start with:
sudo docker compose -f 'docker-compose.yml' up -d --build

for development golang

go build  -o ./bin/main ./cmd/main.go && sudo setcap 'cap_net_bind_service=+ep' ./bin/main &&  ./bin/main

For use browser, You may need import the ca_cert.pem in browser setting.

Incase you want create you own cert for testing:

openssl ecparam -genkey -name secp384r1 -out go_key.pem
openssl req -new -key go_key.pem -out go_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"
openssl x509 -req -days 365 -in go_csr.pem -signkey go_key.pem -out go_cert.pem

openssl ecparam -genkey -name secp384r1 -out ca_key.pem

openssl req -x509 -new -key ca_key.pem -out ca_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-CA" -days 365

openssl x509 -req -in go_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 01 -out go_cert.pem -days 365

In case you want RSA:

# Generate RSA 2048 / 4096cert

openssl genrsa -out rsa2048_key.pem 2048

openssl req -new -key rsa2048_key.pem -out rsa2048_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"

openssl x509 -req -days 365 -in rsa2048_csr.pem -signkey rsa2048_key.pem -out rsa2048_cert.pem

openssl x509 -req -in rsa2048_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 02 -out rsa2048_cert.pem -days 365

openssl genrsa -out rsa4096_key.pem 4096

openssl req -new -key rsa4096_key.pem -out rsa4096_csr.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=localhost"

openssl x509 -req -days 365 -in rsa4096_csr.pem -signkey rsa4096_key.pem -out rsa4096_cert.pem

openssl x509 -req -in rsa4096_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 03 -out rsa4096_cert.pem -days 365

if you want the ca_cert also with rsa:

# Create self-signed CA certificate

openssl genrsa -out ca_rsa2048_key.pem 2048

openssl req -x509 -new -key ca_rsa2048_key.pem -out ca_rsa2048_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-RSA2048-CA" -days 365

# Create self-signed CA certificate

openssl genrsa -out ca_rsa4096_key.pem 4096

openssl req -x509 -new -key ca_rsa4096_key.pem -out ca_rsa4096_cert.pem -subj "/C=HK/ST=HK/L=HK/O=comp4334/CN=comp4334-RSA4096-CA" -days 365
