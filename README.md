# P2P tunnel

A P2P tunnel for secure bidirectional communication between two or more P2P nodes

```
openssl genpkey -algorithm RSA -out server.key
openssl req -new -key server.key -out server.csr
openssl x509 -req -in server.csr -signkey server.key -out server.crt
```