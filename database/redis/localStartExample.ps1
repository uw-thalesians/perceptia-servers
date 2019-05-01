Param (
    [String]$RedisPortPublish = "6379",
    [String]$PerceptiaDockerNet = "perceptia-net"
)

Set-Variable -Name REDIS_SERVICE_NAME -Value "redis"
Set-Variable -Name REDIS_VOLUME_NAME -Value redis_vol

docker rm --force $REDIS_SERVICE_NAME

docker network create -d bridge $PerceptiaDockerNet


docker run `
--detach `
--name $REDIS_SERVICE_NAME `
--network $PerceptiaDockerNet `
--publish "${RedisPortPublish}:6379" `
--mount type=volume,source=${REDIS_VOLUME_NAME},destination=/data `
redis:5.0.4-alpine

Write-Host "Redis Server is listening inside docker network: ${PerceptiaDockerNet} at: ${REDIS_SERVICE_NAME}:1433"
Write-Host "MsSql Server is listening on the host at: localhost:${MsSqlPortPublish}"
