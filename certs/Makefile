.PHONY: build run clean

CA_BITS          = 4096
CA_DAYS          = 3650
CA_COUNTRY       = FI
CA_STATE         = Tampere
CA_LOCALITY      = Tampere
CA_ORG           = ExampleOrg
CA_ORG_UNIT      = Developers
CA_COMMON_NAME   = localhost

# Variables for the config file
SERVER_DAYS          = 3650
SERVER_BITS          = 2048
SERVER_MD            = sha256
SERVER_COUNTRY       = FI
SERVER_STATE         = Tampere
SERVER_LOCALITY      = Tampere
SERVER_ORG           = ExampleOrg
SERVER_ORG_UNIT      = Developers
SERVER_COMMON_NAME   = localhost
SERVER_DNS_NAME      = localhost
SERVER_IP_ADDRESS    = 127.0.0.1

all: build

build: server.crt

ca.key:
	openssl genrsa -out $@ $(CA_BITS)

ca.crt: ca.key
	openssl req -x509 -new -nodes -key ca.key -sha256 -days $(CA_DAYS) -out $@ -subj "/C=$(CA_COUNTRY)/ST=$(CA_STATE)/L=$(CA_LOCALITY)/O=$(CA_ORG)/OU=$(CA_ORG_UNIT)/CN=$(CA_COMMON_NAME)"

server.cnf:
	@echo "[ req ]" > $@
	@echo "default_bits       = $(SERVER_BITS)" >> $@
	@echo "prompt             = no" >> $@
	@echo "default_md         = $(SERVER_MD)" >> $@
	@echo "distinguished_name = req_distinguished_name" >> $@
	@echo "req_extensions     = req_ext" >> $@
	@echo "" >> $@
	@echo "[ req_distinguished_name ]" >> $@
	@echo "C  = $(SERVER_COUNTRY)" >> $@
	@echo "ST = $(SERVER_STATE)" >> $@
	@echo "L  = $(SERVER_LOCALITY)" >> $@
	@echo "O  = $(SERVER_ORG)" >> $@
	@echo "OU = $(SERVER_ORG_UNIT)" >> $@
	@echo "CN = $(SERVER_COMMON_NAME)" >> $@
	@echo "" >> $@
	@echo "[ req_ext ]" >> $@
	@echo "subjectAltName = @alt_names" >> $@
	@echo "" >> $@
	@echo "[ alt_names ]" >> $@
	@echo "DNS.1 = $(SERVER_DNS_NAME)" >> $@
	@echo "IP.1 = $(SERVER_IP_ADDRESS)" >> $@

server.key:
	openssl genpkey -algorithm RSA -out $@ -pkeyopt rsa_keygen_bits:$(SERVER_BITS)

server.csr: server.cnf server.key
	openssl req -new -key server.key -out server.csr -config server.cnf -extensions req_ext

server.crt: server.key server.cnf server.csr ca.crt ca.key server.cnf
	openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $@ -days $(SERVER_DAYS) -sha256 -extfile server.cnf -extensions req_ext

clean:
	rm -f server.key server.crt server.csr server.cnf ca.key ca.crt ca.csr ca.srl
