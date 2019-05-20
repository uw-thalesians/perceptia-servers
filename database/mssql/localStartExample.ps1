Param (
    [String]$MsSqlSaPassword = "SecureNow!",
    [String]$MsSqlGatewaySpUsername = "gateway_sp",
    [String]$MsSqlGatewaySpPassword = "ThisIsReal!",
    [String]$MsSqlPortPublish = "1401",
    [String]$MsSqlSkipSetupIfExist = "Y",
    [String]$MsSqlSkipSetup = "N",
    [switch]$BuildMsSql,
    [string]$MsSqlVersion = "0.8.1",
    [string]$MsSqlBuild = "232",
    [string]$MsSqlBranch = "develop",
    [String]$PerceptiaDockerNet = "perceptia-net",
    [switch]$MsSqlRemoveDbVolume,
    [switch]$CleanUp

)
# File: localStartExample.ps1

Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name MSSQL_SERVICE_NAME -Value "mssql"
Set-Variable -Name MSSQL_IMAGE_NAME -Value "${MSSQL_SERVICE_NAME}"
Set-Variable -Name MSSQL_IMAGE_TAG -Value "${LATEST_COMMIT}"
Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${MSSQL_IMAGE_NAME}:${MSSQL_IMAGE_TAG}"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "mssql-svc_mssql_vol"

if (!$CleanUp) {
    if ($BuildMsSql) {
        docker build --tag "${MSSQL_IMAGE_AND_TAG}" --no-cache .
    } else {
        Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${MSSQL_IMAGE_NAME}:${MsSqlVersion}-build-${MsSqlBuild}-branch-${MsSqlBranch}"
    }
    

    docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=mssql-svc")
    docker rm --force $MSSQL_SERVICE_NAME
    if ($MsSqlRemoveDbVolume) {
        Write-Host "MsSqlRemoveDbVolume option true, removing previous database"
        docker volume rm $MSSQL_VOLUME_NAME
    }

    docker network create -d bridge $PerceptiaDockerNet

    docker run `
    --detach `
    --env 'ACCEPT_EULA=Y' `
    --env "GATEWAY_SP_USERNAME=${MsSqlGatewaySpUsername}" `
    --env "GATEWAY_SP_PASSWORD=${MsSqlGatewaySpPassword}" `
    --env "MSSQL_ENVIRONMENT=development" `
    --env "SA_PASSWORD=$MsSqlSaPassword" `
    --env "SKIP_SETUP_IF_EXISTS=$MsSqlSkipSetupIfExist" `
    --env "SKIP_SETUP=$MsSqlSkipSetup" `
    --label "label.perceptia.info/name=mssql" `
    --label "label.perceptia.info/instance=mssql-1" `
    --label "label.perceptia.info/managed-by=docker" `
    --label "label.perceptia.info/component=database" `
    --label "label.perceptia.info/part-of=mssql-svc" `
    --mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
    --name=${MSSQL_SERVICE_NAME} `
    --network $PerceptiaDockerNet `
    --publish "${MsSqlPortPublish}:1433" `
    ${MSSQL_IMAGE_AND_TAG}

    Write-Host "MsSql Server is listening inside docker network: ${PerceptiaDockerNet} at: ${MSSQL_SERVICE_NAME}:1433"
    Write-Host "MsSql Server is listening on the host at: localhost:${MsSqlPortPublish}"
} else {
    docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=mssql-svc")
    if ($MsSqlRemoveDbVolume) {
        Write-Host "MsSqlRemoveDbVolume option true, removing previous database"
        docker volume rm $MSSQL_VOLUME_NAME
    }
}


