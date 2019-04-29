Param (
    [String]$MsSqlDatabase = "Perceptia",
    [String]$MsSqlHost = "mssql",
    [String]$MsSqlPassword = "SecureNow!",
    [String]$MsSqlPort = "1433",
    [String]$MsSqlScheme = "sqlserver",
    [String]$MsSqlUsername = "sa",
    [String]$PerceptiaDockerNet = "perceptia-net",
    [String]$GatewayPort = "4443"

)
Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name GATEWAY_IMAGE_TAG -Value "${LATEST_COMMIT}"
Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${GATEWAY_IMAGE_NAME}:${GATEWAY_IMAGE_TAG}"
Set-Variable -Name GATEWAY_CONTAINER_NAME -Value "gateway"

docker build --tag "${GATEWAY_IMAGE_AND_TAG}" --no-cache .

Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/gateway_tlscert.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/gateway_tlskey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)\gateway\encrypt\"

Set-Variable -Name MSSQL_DATABASE -Value $MsSqlDatabase
Set-Variable -Name MSSQL_HOST -Value $MsSqlHost
Set-Variable -Name MSSQL_PASSWORD -Value $MsSqlPassword
Set-Variable -Name MSSQL_PORT -Value $MsSqlPort
Set-Variable -Name MSSQL_SCHEME -Value $MsSqlScheme
Set-Variable -Name MSSQL_USERNAME -Value $MsSqlUsername

docker rm --force ${GATEWAY_CONTAINER_NAME}

docker network create -d bridge $PerceptiaDockerNet

docker run `
--detach `
--env GATEWAY_TLSCERTPATH="$GATEWAY_TLSCERTPATH" `
--env GATEWAY_TLSKEYPATH="$GATEWAY_TLSKEYPATH" `
--env MSSQL_DATABASE="$MSSQL_DATABASE" `
--env MSSQL_HOST="$MSSQL_HOST" `
--env MSSQL_PASSWORD="$MSSQL_PASSWORD" `
--env MSSQL_PORT="$MSSQL_PORT" `
--env MSSQL_SCHEME="$MSSQL_SCHEME" `
--env MSSQL_USERNAME="$MSSQL_USERNAME" `
--name ${GATEWAY_CONTAINER_NAME} `
--network $PerceptiaDockerNet `
--publish "${GatewayPort}:443" `
--restart on-failure `
--mount type=bind,source="$GATEWAY_TLSMOUNTSOURCE",target="/encrypt/",readonly `
${GATEWAY_IMAGE_AND_TAG}

Write-Host "Gateway is listening at https://localhost:${GatewayPort}"