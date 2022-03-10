run:
	go run . -tls \
	-tls_ca_path "./test/ca.pem" \
	-tls_cert_path "./test/server/server.pem" \
	-tls_key_path "./test/server/server.key" \
	-tls_verify_client
ca:
	# openssl genrsa -out test/ca.key 2048
	# openssl req -new -x509 -days 3650 -key test/ca.key -out test/ca.pem
server:
	openssl genrsa -out test/server/server.key 2048
	openssl req -new -key test/server/server.key -out test/server/server.csr
	openssl x509 -req -sha256 -CA test/ca.pem -CAkey test/ca.key -CAcreateserial -days 3650 -in test/server/server.csr -out test/server/server.pem
client:
	openssl ecparam -genkey -name secp384r1 -out test/client/client.key
	openssl req -new -key test/client/client.key -out test/client/client.csr
	openssl x509 -req -sha256 -CA test/ca.pem -CAkey test/ca.key -CAcreateserial -days 3650 -in test/client/client.csr -out test/client/client.pem
