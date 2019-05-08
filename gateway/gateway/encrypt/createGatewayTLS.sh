#!/usr/bin/env bash
echo "Creating request and generating x509 cert and key"

openssl req -x509 -out gateway_tlscert.pem -keyout gateway_tlskey.pem \
    -newkey rsa:2048 -nodes -sha256 \
    -extensions EXT -config openssl.conf
