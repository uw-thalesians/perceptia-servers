# File: localStartExample.ps1

Set-Variable -Name MSSQL_SCHEMA_VERSION -Value "0.3.2"
Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name DOCKERHUB_NAME -Value "uwthalesians"
Set-Variable -Name MSSQL_SERVICE_NAME -Value "mssql"
Set-Variable -Name MSSQL_IMAGE_NAME -Value "${DOCKERHUB_NAME}/${MSSQL_SERVICE_NAME}"
Set-Variable -Name MSSQL_IMAGE_TAG -Value "${MSSQL_SCHEMA_VERSION}-${LATEST_COMMIT}"
Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${MSSQL_IMAGE_NAME}:${MSSQL_IMAGE_TAG}"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "mssql_vol"

docker build --tag "${MSSQL_IMAGE_AND_TAG}" .

docker rm --force ${MSSQL_SERVICE_NAME}
#docker volume rm ${MSSQL_VOLUME_NAME}

docker run `
--detach `
--env 'ACCEPT_EULA=Y' `
--env "SA_PASSWORD=$Env:SA_PASSWORD" `
--env "SKIP_SETUP_IF_EXISTS=Y" `
--mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
--name=mssql `
--publish 1401:1433 `
${MSSQL_IMAGE_AND_TAG}
