# Local Development Encryption Files

## Gateway server TLS Key and Cert

The script [createTlsCert.sh](./createTlsCert.sh) will generate a private key and a certificate for "localhost" signed with that key.

The script uses the openssl application to generate the key and cert and is intended to be run in a bash interpreter. . The local .gitignore file in this directory ensures that the generated .pem files are not added to the Git repository.

The openssl.conf file provides additional configuration information to the openssl command used to generate the key and cert.

Note, you may need to enable executing the bash script. To do this, open a bash shell and type `chmod +x ./createTlsCert.sh`