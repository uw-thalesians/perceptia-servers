kubectl create secret tls api-tls --cert=$PERCEPTIA_SERVERS_SECRET\ssl\archive\fullchain1.pem --key=$PERCEPTIA_SERVERS_SECRET\ssl\archive\privkey1.pem --namespace production
kubectl create secret tls api-tls --cert=$PERCEPTIA_SERVERS_SECRET\ssldev\archive\fullchain1.pem --key=$PERCEPTIA_SERVERS_SECRET\ssldev\archive\privkey1.pem --namespace development
