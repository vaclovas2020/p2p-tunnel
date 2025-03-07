#!/bin/bash

# === CONFIGURATION ===
CA_CN="MyCA"
SERVER_CN="myserver.com"
SERVER_IP="192.168.1.101"

# Step 1: Generate CA key and self-signed certificate
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/C=US/ST=MyState/L=MyCity/O=MyOrganization/CN=$CA_CN"

# Step 2: Create a server SAN config file
cat > server_san.cnf <<EOF
[ req ]
default_bits       = 4096
prompt            = no
default_md        = sha256
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
C  = US
ST = MyState
L  = MyCity
O  = MyOrganization
CN = $SERVER_CN

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $SERVER_CN
DNS.2 = localhost
IP.1  = $SERVER_IP
IP.2  = 127.0.0.1
EOF

# Step 3: Generate server key and CSR
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr -config server_san.cnf

# Step 4: Sign the server certificate with the CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days 3650 -sha256 -extfile server_san.cnf -extensions req_ext

# Step 8: Verify that the certificates include SANs
echo "Verifying server certificate SANs:"
openssl x509 -in server.crt -noout -text | grep -A 1 "Subject Alternative Name"

echo "✅ Certificate generation completed successfully!"
