
Set-Variable -Name GATEWAY_VERSION -Value "0.1.1"
Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name DOCKERHUB_NAME -Value "uwthalesians"
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "${DOCKERHUB_NAME}/gateway"
Set-Variable -Name GATEWAY_IMAGE_TAG -Value "${GATEWAY_VERSION}-${LATEST_COMMIT}"
Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${GATEWAY_IMAGE_NAME}:${GATEWAY_IMAGE_TAG}"

docker build --tag "${GATEWAY_IMAGE_AND_TAG}" .

Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/gateway_tlscert.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/gateway_tlskey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)\gateway\encrypt\"

docker rm --force gateway

docker run `
--detach `
--env GATEWAY_TLSCERTPATH="${GATEWAY_TLSCERTPATH}" `
--env GATEWAY_TLSKEYPATH="${GATEWAY_TLSKEYPATH}" `
--name gateway `
--publish "4443:443" `
--restart on-failure `
--mount type=bind,source="$GATEWAY_TLSMOUNTSOURCE",target="/encrypt/",readonly `
${GATEWAY_IMAGE_AND_TAG}