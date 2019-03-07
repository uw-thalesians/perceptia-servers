
Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name GATEWAY_IMAGE_TAG -Value "${LATEST_COMMIT}"
Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${GATEWAY_IMAGE_NAME}:${GATEWAY_IMAGE_TAG}"
Set-Variable -Name GATEWAY_CONTAINER_NAME -Value "gateway"

docker build --tag "${GATEWAY_IMAGE_AND_TAG}" .

Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/gateway_tlscert.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/gateway_tlskey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)\gateway\encrypt\"

docker rm --force ${GATEWAY_CONTAINER_NAME}

docker run `
--detach `
--env GATEWAY_TLSCERTPATH="${GATEWAY_TLSCERTPATH}" `
--env GATEWAY_TLSKEYPATH="${GATEWAY_TLSKEYPATH}" `
--name ${GATEWAY_CONTAINER_NAME} `
--publish "4443:443" `
--restart on-failure `
--mount type=bind,source="$GATEWAY_TLSMOUNTSOURCE",target="/encrypt/",readonly `
${GATEWAY_IMAGE_AND_TAG}