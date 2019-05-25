Param (
    [switch]$Latest,
    [string]$Build = "329",
    [String]$Branch = "develop",
    [switch]$CurrentBranch,

    [String]$PerceptiaDockerNet = "perceptia-net",

    [switch]$BuildGateway = $false,
    [String]$GatewayVersion = "1",
    [String]$GatewayPortPublish = "4443",

    [switch]$SkipMsSql = $false,
    [String]$MsSqlVersion = "1.0.0",
    [String]$MsSqlDatabase = "Perceptia",
    [String]$MsSqlHost = "mssql",
    [String]$MsSqlPort = "1433",
    [String]$MsSqlSaPassword = "SecureNow!",
    [String]$MsSqlPortPublish = "1401",
    [String]$MsSqlScheme = "sqlserver",
    [String]$MsSqlUsername = "sa",
    [String]$MsSqlGatewaySpUsername = "gateway_sp",
    [String]$MsSqlGatewaySpPassword = "ThisIsReal!",
    [switch]$MsSqlRemoveDbVolume = $false,

    [switch]$SkipRedis = $false,
    [String]$RedisPort = "6379",
    [String]$RedisPortPublish = "6379",
    [String]$RedisHost = "redis",
    [switch]$RedisRemoveDbVolume = $false,

    [String]$AqRestPort = "80",
    [String]$AqRestHost = "aqrest",

    [switch]$CleanUp = $false,
    [switch]$RemoveAllDbVolumes = $false
)

# Setup Base Veriables
Set-Variable -Name DOCKERHUB_ORG -Value "uwthalesians"

Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"

Set-Variable -Name DOCKER_NETWORK -Value $PerceptiaDockerNet

## Gateway Veriables
Set-Variable -Name GATEWAY_SERVICE_NAME -Value "gateway-lc-svc"

Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name GATEWAY_CONTAINER_NAME -Value "gateway"
Set-Variable -Name GATEWAY_PORT_PUBLISH -Value $GatewayPortPublish


Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/fullchain.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/privkey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)/gateway/encrypt"

Set-Variable -Name GATEWAY_SESSION_KEY -Value "fjsfndreifnfsnm5kngfnklef23kdnfskng"

Set-Variable -Name GATEWAY_API_PORT -Value "$GatewayPortPublish"

## Redis Variables
Set-Variable -Name REDIS_CONTAINER_NAME -Value $RedisHost
Set-Variable -Name REDIS_VOLUME_NAME -Value "${GATEWAY_SERVICE_NAME}_redis_vol"
Set-Variable -Name REDIS_PORT_PUBLISH -Value $RedisPortPublish
Set-Variable -Name REDIS_PORT -Value $RedisPort


## MsSql Variables
Set-Variable -Name MSSQL_IMAGE_NAME -Value "mssql"
Set-Variable -Name MSSQL_CONTAINER_NAME -Value "mssql"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "${GATEWAY_SERVICE_NAME}_mssql_vol"
Set-Variable -Name MSSQL_PORT_PUBLISH -Value $MsSqlPortPublish

Set-Variable -Name MSSQL_SA_PASSWORD -Value $MsSqlSaPassword
Set-Variable -Name MSSQL_GATEWAY_SP_USERNAME -Value $MsSqlGatewaySpUsername
Set-Variable -Name MSSQL_GATEWAY_SP_PASSWORD -Value $MsSqlGatewaySpPassword

Set-Variable -Name MSSQL_DATABASE -Value $MsSqlDatabase
Set-Variable -Name MSSQL_HOST -Value $MsSqlHost
Set-Variable -Name MSSQL_PORT -Value $MsSqlPort
Set-Variable -Name MSSQL_SCHEME -Value $MsSqlScheme




if (!$CleanUp) {

    if (($GatewayVersion).Length -eq 0) {
        Write-Host "Version must be provided, but no version provided for gateway, exiting..."
        exit(1)
}
    # Define Image Tags to use
    Set-Variable -Name TAG_BRANCH -Value $Branch
    if ($CurrentBranch) {
            Set-Variable -Name TAG_BRANCH -Value ((git rev-parse --abbrev-ref HEAD) -replace "^(?:(?:[^//]{0,})[/]{1,1}){1,}")
            Write-Host "CurrentBranch switch provided, using build from branch: $TAG_BRANCH"
    }

    if (($TAG_BRANCH).Length -eq 0) {
            Write-Host "Branch must be set, but no branch set, exiting"
            exit(1)
    }

    
    Set-Variable -Name TAG_BUILD -Value $Build # Build number 
    
    if ($Latest) {
            Write-Host "Latest switch provided, using latest build from branch: $TAG_BRANCH"
            Set-Variable -Name TAG_BUILD -Value "latest"              
    } else {
            Write-Host "Using build $TAG_BUILD from branch: $TAG_BRANCH"
            Set-Variable -Name TAG_BUILD -Value "$Build"              
    }
    if (($TAG_BUILD).Length -eq 0) {
            Write-Host "Build must be provided, but no build provided, exiting..."
            exit(1)
    } 
    Write-Host "`n"
    
    Set-Variable -Name BUILD_AND_BRANCH -Value "build-${TAG_BUILD}-branch-${TAG_BRANCH}"

    if ($BuildGateway) {
        Write-Host "BuildGateway switch provided, building gateway locally, branch and build flags will be ignored for gateway"
        Set-Variable -Name GATEWAY_IMAGE_TAG -Value "${LATEST_COMMIT}"
        Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${GATEWAY_IMAGE_NAME}:${GATEWAY_IMAGE_TAG}"
        Write-Host "Building gateway image: $GATEWAY_IMAGE_AND_TAG"
        docker build --tag "${GATEWAY_IMAGE_AND_TAG}" --no-cache .
    } else {
        
        Set-Variable -Name GATEWAY_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${GATEWAY_IMAGE_NAME}:${GatewayVersion}-${BUILD_AND_BRANCH}"
        Write-Host "Using gateway image: $GATEWAY_IMAGE_AND_TAG"
        Write-Host "`n"
    } 
    
    # Remove Existing Containers
    if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
        Write-Host "Removing all containers started by this script..."
        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
    }

    #Remove existing volumes
    if ($RemoveAllDbVolumes -and (((docker volume ls --format "{{.Name}}") -Match "${GATEWAY_SERVICE_NAME}"))) {
        Write-Host "Database volume reset requested"
        if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
                Write-Host "Removing all containers before removing volumes"
                docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        }  
        docker volume rm (docker volume ls --format "{{.Name}}" --filter "name=${GATEWAY_SERVICE_NAME}")
    }

    if (($MsSqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "$MSSQL_VOLUME_NAME" ))){
        Write-Host "-MsSqlRemoveDbVolume option set, removing volume: $MSSQL_VOLUME_NAME"
        if ((docker ps -aq --filter "label=label.perceptia.info/name=mssql" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
                docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=mssql"  --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        } 
        Start-Sleep -Seconds 2
        docker volume rm $MSSQL_VOLUME_NAME
    }
    if (($RedisRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "$REDIS_VOLUME_NAME" ))) {
        Write-Host "-RedisRemoveDbVolume option set, removing volume: $REDIS_VOLUME_NAME"
        if ((docker ps -aq --filter "label=label.perceptia.info/name=redis" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
                docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=redis"  --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAMEE}")
        } 
        
        Start-Sleep -Seconds 2
        docker volume rm $REDIS_VOLUME_NAME
    }
    
    if (!((docker network ls) -Match $DOCKER_NETWORK)) {
        Write-Host "Docker network $DOCKER_NETWORK does not exist, creating now..."
        docker network create -d bridge $DOCKER_NETWORK 
        Write-Host "`n"
    }
    


    
    if (!$SkipRedis) {
        Write-Host "SkipRedis option false, starting redis dependency"

        Set-Variable -Name REDIS_IMAGE_AND_TAG -Value "redis:5.0.4-alpine"
        # Check if image exists
        (docker pull $REDIS_IMAGE_AND_TAG) | Out-Null 
        if (!$?) {
                Write-Host "Image: $REDIS_IMAGE_AND_TAG not found on dockerhub, exiting"
                exit(1)
        }
    
        docker run `
        --detach `
        --label "label.perceptia.info/name=${REDIS_CONTAINER_NAME}" `
        --label "label.perceptia.info/instance=${REDIS_CONTAINER_NAME}-1" `
        --label "label.perceptia.info/managed-by=docker" `
        --label "label.perceptia.info/component=database" `
        --label "label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}" `
        --label "label.perceptia.info/sub-of=${GATEWAY_SERVICE_NAME}" `
        --name $REDIS_CONTAINER_NAME `
        --network $DOCKER_NETWORK `
        --publish "${REDIS_PORT_PUBLISH}:6379" `
        --mount type=volume,source=${REDIS_VOLUME_NAME},destination=/data `
        redis:5.0.4-alpine
        Write-Host "`n"
        Write-Host "Redis Server is listening inside docker network: ${DOCKER_NETWORK} at: ${REDIS_CONTAINER_NAME}:6379"
        Write-Host "Redis Server is listening on the host at: localhost:${REDIS_PORT_PUBLISH}"
        Write-Host "`n"
    } else {
        Write-Host "Be sure to start mssql dependency, see README"
    }
    Write-Host "`n"
    if (!$SkipMsSql) {
        Write-Host "SkipMsSql option false, starting mssql dependency"

        if (($MsSqlVersion).Length -eq 0) {
            Write-Host "Version must be provided, but no version provided for mssql, exiting..."
            exit(1)
        }

        Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${MSSQL_IMAGE_NAME}:${MsSqlVersion}-${BUILD_AND_BRANCH}"
        Write-Host "Using mssql image: $MSSQL_IMAGE_AND_TAG"

        # Check if image exists
        Write-Host "Checking Dockerhub for image..."
        (docker pull $MSSQL_IMAGE_AND_TAG) | Out-Null 
        if (!$?) {
                Write-Host "Image: $MSSQL_IMAGE_AND_TAG not found on dockerhub, exiting"
                exit(1)
        }
        Write-Host "Mssql image found!"
    
        docker run `
        --detach `
        --env 'ACCEPT_EULA=Y' `
        --env "GATEWAY_SP_USERNAME=${MSSQL_GATEWAY_SP_USERNAME}" `
        --env "GATEWAY_SP_PASSWORD=${MSSQL_GATEWAY_SP_PASSWORD}" `
        --env "MSSQL_ENVIRONMENT=development" `
        --env "SA_PASSWORD=$MSSQL_SA_PASSWORD" `
        --env "SKIP_SETUP_IF_EXISTS=Y" `
        --label "label.perceptia.info/name=${MSSQL_CONTAINER_NAME}" `
        --label "label.perceptia.info/instance=${MSSQL_CONTAINER_NAME}-1" `
        --label "label.perceptia.info/managed-by=docker" `
        --label "label.perceptia.info/component=database" `
        --label "label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}" `
        --label "label.perceptia.info/sub-of=${GATEWAY_SERVICE_NAME}" `
        --mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
        --name=${MSSQL_CONTAINER_NAME} `
        --network $DOCKER_NETWORK `
        --publish "${MSSQL_PORT_PUBLISH}:1433" `
        ${MSSQL_IMAGE_AND_TAG}
        Write-Host "`n"
        Write-Host "MsSql Server is listening inside docker network: ${DOCKER_NETWORK} at: ${MSSQL_CONTAINER_NAME}:1433"
        Write-Host "MsSql Server is listening on the host at: localhost:${MSSQL_PORT_PUBLISH}"
        Write-Host "`n"
    } else {
        Write-Host "Be sure to start mssql dependency, see README"
    }
    Write-Host "Using gateway image: $GATEWAY_IMAGE_AND_TAG"
    if (!$BuildGateway) {
        Write-Host "Checking Dockerhub for image..."
        # Check if image exists
        (docker pull $GATEWAY_IMAGE_AND_TAG) | Out-Null 
        if (!$?) {
                Write-Host "Image: $GATEWAY_IMAGE_AND_TAG not found on dockerhub, exiting"
                exit(1)
        }
        Write-Host "Gateway image found!"
    }
    
    docker run `
    --detach `
    --env AQREST_HOSTNAME="$AqRestHost" `
    --env AQREST_PORT="$AqRestPort" `
    --env GATEWAY_API_PORT=$GATEWAY_API_PORT `
    --env GATEWAY_ENVIRONMENT=development `
    --env GATEWAY_SESSION_KEY="$GATEWAY_SESSION_KEY" `
    --env GATEWAY_TLSCERTPATH="$GATEWAY_TLSCERTPATH" `
    --env GATEWAY_TLSKEYPATH="$GATEWAY_TLSKEYPATH" `
    --env MSSQL_DATABASE="$MSSQL_DATABASE" `
    --env MSSQL_HOST="$MSSQL_HOST" `
    --env MSSQL_GATEWAY_SP_PASSWORD="$MSSQL_GATEWAY_SP_PASSWORD" `
    --env MSSQL_PORT="$MSSQL_PORT" `
    --env MSSQL_SCHEME="$MSSQL_SCHEME" `
    --env MSSQL_GATEWAY_SP_USERNAME="$MSSQL_GATEWAY_SP_USERNAME" `
    --env REDIS_ADDRESS="${REDIS_CONTAINER_NAME}:$REDIS_PORT" `
    --label "label.perceptia.info/name=${GATEWAY_CONTAINER_NAME}" `
    --label "label.perceptia.info/instance=gateway-1" `
    --label "label.perceptia.info/managed-by=docker" `
    --label "label.perceptia.info/component=server" `
    --label "label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}" `
    --name ${GATEWAY_CONTAINER_NAME} `
    --network $DOCKER_NETWORK `
    --publish "${GATEWAY_PORT_PUBLISH}:443" `
    --restart on-failure `
    --mount type=bind,source=$GATEWAY_TLSMOUNTSOURCE,target=/encrypt/,readonly `
    ${GATEWAY_IMAGE_AND_TAG}
    
    Write-Host "`n"
    Write-Host "Gateway is listening at https://localhost:${GATEWAY_PORT_PUBLISH}"
    Write-Host "`n"
    Write-Host "To test if the gateway is able to process requests, make a GET request to:"
    Write-Host "/api/v1/gateway/health"
    Write-Host "Example request using curl"
    Write-Host "curl --insecure -X GET `"https://localhost:${GATEWAY_PORT_PUBLISH}/api/v1/gateway/health`""
    Write-Host "`n"
    Write-Host "Mssql server for the Perceptia databased used by the gateway can be reached at port: $MSSQL_PORT_PUBLISH"
    Write-Host "Perceptia databased used by the gateway has an 'sa' user account with password: $MSSQL_SA_PASSWORD"
    Write-Host "`n"
    docker ps --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}"
} else {

    if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
        Write-Host "Removing all contaienrs started by this script..."
        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
    }
    if ($RemoveAllDbVolumes -and (((docker volume ls --format "{{.Name}}") -Match "${GATEWAY_SERVICE_NAME}"))) {
        Write-Host "Database volume reset requested"
        if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
                Write-Host "Removing all containers before removing volumes"
                docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        }  
        Start-Sleep -Seconds 2
        docker volume rm (docker volume ls --format "{{.Name}}" --filter "name=${GATEWAY_SERVICE_NAME}")
    }
    if (($MsSqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "$MSSQL_VOLUME_NAME" ))){
        Write-Host "-MsSqlRemoveDbVolume option set, removing volume: $MSSQL_VOLUME_NAME"
        if ((docker ps -aq --filter "label=label.perceptia.info/name=mssql" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
            Write-Host "Removing any mssql containers started by this script to allow removing the database volume: $MSSQL_VOLUME_NAME"
            docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=mssql"  --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        } 
        Start-Sleep -Seconds 2
        Write-Host "Removing mssql volume..."
        docker volume rm $MSSQL_VOLUME_NAME
    }
    if (($RedisRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "$REDIS_VOLUME_NAME" ))) {
        Write-Host "-RedisRemoveDbVolume option set, removing volume: $REDIS_VOLUME_NAME"
        if ((docker ps -aq --filter "label=label.perceptia.info/name=redis" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")) {
            Write-Host "Removing any redis containers started by this script to allow removing the database volume: $REDIS_VOLUME_NAME"
            docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=redis"  --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAMEE}")
        } 
        
        Start-Sleep -Seconds 2
        Write-Host "Removing redis volume..."
        docker volume rm $REDIS_VOLUME_NAME
    }
    Write-Host "CleanUp Complete!"
}

