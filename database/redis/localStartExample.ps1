
docker rm --force redis

Write-Output "running redis image, container name: redis, publish: 6379:6379"
docker run `
--detach `
--name redis `
--publish "6379:6379" `
redis
