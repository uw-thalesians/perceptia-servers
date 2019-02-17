#Local Development Encryption Files
##Gateway server TLS Key and Cert
The script createGatewayTLS.sh will generate a private key and a certificate for "localhost" signed with that key.
The script uses the openssl application to generate the key and cert. The local .gitignore file in this directory ensures that the generated .pem files are not added to the Git repository. 