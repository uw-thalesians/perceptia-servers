kubectl create secret tls api-tls --cert=$PERCEPTIA_SERVERS_SECRET\cert\fullchain.pem --key=$PERCEPTIA_SERVERS_SECRET\cert\privkey.pem --namespace production
kubectl create secret tls api-tls --cert=$PERCEPTIA_SERVERS_SECRET\certdev\fullchain.pem --key=$PERCEPTIA_SERVERS_SECRET\certdev\privkey.pem --namespace development
