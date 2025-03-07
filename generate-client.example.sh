#!/bin/bash

# === CONFIGURATION ===
CLIENT_CN="myclient.com"
CLIENT_IP="192.168.1.100"

# Step 5: Create a client SAN config file
cat > client_san.cnf <<EOF
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
CN = $CLIENT_CN

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $CLIENT_CN
IP.1  = $CLIENT_IP
EOF

# Step 6: Generate client key and CSR
openssl genrsa -out client.key 4096
openssl req -new -key client.key -out client.csr -config client_san.cnf

# Step 7: Sign the client certificate with the CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out client.crt -days 3650 -sha256 -extfile client_san.cnf -extensions req_ext

echo "Verifying client certificate SANs:"
openssl x509 -in client.crt -noout -text | grep -A 1 "Subject Alternative Name"

echo "✅ Certificate generation completed successfully!"
