#!/bin/bash

# Step 1: Generate CA key and self-signed certificate
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/CN=MyCA"

# Step 2: Create a server SAN config file (server_san.cnf)
cat > server_san.cnf <<EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
commonName = myserver.com

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = myserver.com
DNS.2 = localhost
IP.1  = 127.0.0.1
EOF

# Step 3: Generate server key and CSR
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config server_san.cnf

# Step 4: Sign the server certificate with the CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days 3650 -sha256 -extfile server_san.cnf -extensions req_ext

# Step 5: Create a client SAN config file (client_san.cnf)
cat > client_san.cnf <<EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
commonName = myclient.com

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = myclient.com
IP.1  = 192.168.1.100
EOF

# Step 6: Generate client key and CSR
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -config client_san.cnf

# Step 7: Sign the client certificate with the CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out client.crt -days 3650 -sha256 -extfile client_san.cnf -extensions req_ext

# Step 8: Verify that the certificates include SANs
echo "Verifying server certificate SANs:"
openssl x509 -in server.crt -noout -text | grep -A 1 "Subject Alternative Name"

echo "Verifying client certificate SANs:"
openssl x509 -in client.crt -noout -text | grep -A 1 "Subject Alternative Name"
