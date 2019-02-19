#!/usr/bin/env bash
set -euo pipefail

echo ""
echo "++++++++ GENERATE CA KEY AND CERTIFICATES +++++++"
openssl genrsa -out ca.key 4096 > /dev/null
openssl req -new -x509 -days 365 -key ca.key -out ca.crt --subj "/C=DE/ST=NRW/O=goss/OU=IT/CN=localhost:8081" > /dev/null

echo ""
echo "++++++++ GENERATE SERVER KEY AND CSR++++++++++++"
openssl genrsa -out server.key 1024 > /dev/null
openssl req -new -key server.key -out server.csr --subj "/C=DE/ST=NRW/O=goss/OU=IT/CN=localhost:8081" > /dev/null

echo ""
echo "+++++++ SIGN SERVER CSR +++++++++++++++++"
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt > /dev/null


echo ""
echo "++++++++ GENERATE CLIENT KEY AND CSR++++++++++++"
openssl genrsa -out client.key 1024 > /dev/null
openssl req -new -key client.key -out client.csr -subj "/C=DE/ST=NRW/O=goss/OU=IT/CN=localhost:8081" > /dev/null

echo ""
echo "++++++++ SIGN CLIENT CERTIFICATE REQUEST ++++++++++++"
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt > /dev/null