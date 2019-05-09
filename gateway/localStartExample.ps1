Param (
    [String]$MsSqlDatabase = "Perceptia",
    [String]$MsSqlHost = "mssql",
    [String]$MsSqlPassword = "SecureNow!",
    [String]$MsSqlPort = "1433",
    [String]$MsSqlPortPublish = "1401",
    [String]$MsSqlScheme = "sqlserver",
    [String]$MsSqlUsername = "sa",
    [String]$PerceptiaDockerNet = "perceptia-net",
    [String]$GatewayPort = "4443",
    [String]$RedisPort = "6379",
    [String]$RedisPortPublish = "6379",
    [String]$RedisHost = "redis",
    [String]$AqRestPort = "80",
    [String]$AqRestHost = "aqrest",
    [switch]$SkipRedis = $false,
    [switch]$KeepRedisDb = $false,
    [switch]$KeepMsSqlDb = $false,
    [switch]$SkipMsSql = $false,
    [switch]$BuildGateway = $false
)

if ($SkipRedis) {
    Write-Host "Be sure to start mssql dependency, see README"
}
if ($SkipMsSql) {
    Write-Host "Be sure to start mssql dependency, see README"
}

Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name GATEWAY_IMAGE_TAG -Value "${LATEST_COMMIT}"
Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${GATEWAY_IMAGE_NAME}:${GATEWAY_IMAGE_TAG}"
Set-Variable -Name GATEWAY_CONTAINER_NAME -Value "gateway"


if ($BuildGateway) {
    Write-Host "Building gateway image: $GATEWAY_IMAGE_AND_TAG"
    docker build --tag "${GATEWAY_IMAGE_AND_TAG}" --no-cache .
} else {
    Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "uwthalesians/gateway:0.2.0-build-129-branch-merge"
}

Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/gateway_tlscert.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/gateway_tlskey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)/gateway/encrypt"

Set-Variable -Name GATEWAY_SESSION_KEY -Value "fjsfndreifnfsnm5kngfnklef23kdnfskng"

Set-Variable -Name MSSQL_DATABASE -Value $MsSqlDatabase
Set-Variable -Name MSSQL_HOST -Value $MsSqlHost
Set-Variable -Name MSSQL_PASSWORD -Value $MsSqlPassword
Set-Variable -Name MSSQL_PORT -Value $MsSqlPort
Set-Variable -Name MSSQL_SCHEME -Value $MsSqlScheme
Set-Variable -Name MSSQL_USERNAME -Value $MsSqlUsername

docker rm --force ${GATEWAY_CONTAINER_NAME}

docker network create -d bridge $PerceptiaDockerNet

if (!$SkipRedis) {
    Write-Host "SkipRedis option false, starting redis dependency"
    Set-Variable -Name REDIS_SERVICE_NAME -Value $RedisHost
    Set-Variable -Name REDIS_VOLUME_NAME -Value redis_vol

    docker rm --force $REDIS_SERVICE_NAME

    if (!$KeepRedisDb) {
        Write-Host "KeepRedisDb option false, clearing old sessions"
        docker volume rm $REDIS_VOLUME_NAME
    }

    docker run `
    --detach `
    --name $REDIS_SERVICE_NAME `
    --network $PerceptiaDockerNet `
    --publish "${RedisPortPublish}:6379" `
    --mount type=volume,source=${REDIS_VOLUME_NAME},destination=/data `
    redis:5.0.4-alpine
    Write-Host "Redis Server is listening inside docker network: ${PerceptiaDockerNet} at: ${REDIS_SERVICE_NAME}:1433"
    Write-Host "Redis Server is listening on the host at: localhost:${MsSqlPortPublish}"
}

if (!$SkipMsSql) {
    Write-Host "SkipMsSql option false, starting mssql dependency"
    Set-Variable -Name MSSQL_SERVICE_NAME -Value "mssql"
    Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value uwthalesians/mssql:0.7.1-build-129-branch-merge
    Set-Variable -Name MSSQL_VOLUME_NAME -Value "mssql_vol"

    docker rm --force ${MSSQL_SERVICE_NAME}

    if (!$KeepMsSqlDb) {
        Write-Host "KeepMsSqlDb option false, removing previous database"
        docker volume rm ${MSSQL_VOLUME_NAME}
    }

    docker run `
    --detach `
    --env 'ACCEPT_EULA=Y' `
    --env "SA_PASSWORD=$MsSqlPassword" `
    --env "SKIP_SETUP_IF_EXISTS=Y" `
    --mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
    --name=${MSSQL_SERVICE_NAME} `
    --network $PerceptiaDockerNet `
    --publish "${MsSqlPortPublish}:1433" `
    ${MSSQL_IMAGE_AND_TAG}

    Write-Host "MsSql Server is listening inside docker network: ${PerceptiaDockerNet} at: ${MSSQL_SERVICE_NAME}:1433"
    Write-Host "MsSql Server is listening on the host at: localhost:${MsSqlPortPublish}"
}

docker run `
--detach `
--env AQREST_HOSTNAME="$AqRestHost" `
--env AQREST_PORT="$AqRestPort" `
--env GATEWAY_SESSION_KEY="$GATEWAY_SESSION_KEY" `
--env GATEWAY_TLSCERTPATH="$GATEWAY_TLSCERTPATH" `
--env GATEWAY_TLSKEYPATH="$GATEWAY_TLSKEYPATH" `
--env MSSQL_DATABASE="$MSSQL_DATABASE" `
--env MSSQL_HOST="$MSSQL_HOST" `
--env MSSQL_PASSWORD="$MSSQL_PASSWORD" `
--env MSSQL_PORT="$MSSQL_PORT" `
--env MSSQL_SCHEME="$MSSQL_SCHEME" `
--env MSSQL_USERNAME="$MSSQL_USERNAME" `
--env REDIS_ADDRESS="${RedisHost}:$RedisPort" `
--name ${GATEWAY_CONTAINER_NAME} `
--network $PerceptiaDockerNet `
--publish "${GatewayPort}:443" `
--restart on-failure `
--mount type=bind,source=$GATEWAY_TLSMOUNTSOURCE,target=/encrypt/,readonly `
${GATEWAY_IMAGE_AND_TAG}

Write-Host "Gateway is listening at https://localhost:${GatewayPort}"