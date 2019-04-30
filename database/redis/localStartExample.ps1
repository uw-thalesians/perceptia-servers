Param (
    [String]$RedisPort = "6379",
    [String]$PerceptiaDockerNet = "perceptia-net"
)

Set-Variable -Name REDIS_VOLUME_NAME -Value redis_vol

docker rm --force redis

docker network create -d bridge $PerceptiaDockerNet


docker run `
--detach `
--name redis `
--network $PerceptiaDockerNet `
--publish "${RedisPort}:6379" `
--mount type=volume,source=${REDIS_VOLUME_NAME},destination=/data `
redis:5.0.4-alpine
