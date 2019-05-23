# Add production certs

kubectl create secret tls api-tls `
--cert=$Env:SECRET_PERCEPTIA_SERVERS\ssl\archive\api.perceptia.info\fullchain1.pem `
--key=$Env:SECRET_PERCEPTIA_SERVERS\ssl\archive\api.perceptia.info\privkey1.pem `
--namespace production

# Add development certs

kubectl create secret tls api-tls `
--cert=$Env:SECRET_PERCEPTIA_SERVERS\ssldev\archive\api.dev.perceptia.info\fullchain1.pem `
--key=$Env:SECRET_PERCEPTIA_SERVERS\ssldev\archive\api.dev.perceptia.info\privkey1.pem `
--namespace development
