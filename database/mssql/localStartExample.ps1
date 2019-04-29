Param (
    [String]$MsSqlPassword = "SecureNow!",
    [String]$MsSqlPort = "1401",
    [String]$MsSqlSkipSetupIfExist = "N",
    [String]$MsSqlSkipSetup = "N",
    [String]$PerceptiaDockerNet = "perceptia-net"

)
# File: localStartExample.ps1

Set-Variable -Name MSSQL_SCHEMA_VERSION -Value "0.5.0"
Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name MSSQL_SERVICE_NAME -Value "mssql"
Set-Variable -Name MSSQL_IMAGE_NAME -Value "${MSSQL_SERVICE_NAME}"
Set-Variable -Name MSSQL_IMAGE_TAG -Value "${MSSQL_SCHEMA_VERSION}-${LATEST_COMMIT}"
Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${MSSQL_IMAGE_NAME}:${MSSQL_IMAGE_TAG}"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "mssql_vol"

docker build --tag "${MSSQL_IMAGE_AND_TAG}" --no-cache .

docker rm --force ${MSSQL_SERVICE_NAME}
#docker volume rm ${MSSQL_VOLUME_NAME}

docker network create -d bridge $PerceptiaDockerNet

docker run `
--detach `
--env 'ACCEPT_EULA=Y' `
--env "SA_PASSWORD=$MsSqlPassword" `
--env "SKIP_SETUP_IF_EXISTS=$MsSqlSkipSetupIfExist" `
--env "SKIP_SETUP=$MsSqlSkipSetup" `
--mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
--name=${MSSQL_SERVICE_NAME} `
--network $PerceptiaDockerNet `
${MSSQL_IMAGE_AND_TAG}

Write-Host "MsSql Server is listening inside docker network: ${PerceptiaDockerNet} at: ${MSSQL_SERVICE_NAME}:1433"
