Param (
    [string]$Build = "232",
    [String]$Branch = "develop",
    [switch]$CurrentBranch,
    [switch]$Latest,
    [String]$GatewayVersion = "0.3.0",
    [String]$MsSqlVersion = "0.7.1",
    [String]$MsSqlDatabase = "Perceptia",
    [String]$MsSqlHost = "mssql",
    [String]$MsSqlSaPassword = "SecureNow!",
    [String]$MsSqlPort = "1433",
    [String]$MsSqlPortPublish = "1401",
    [String]$MsSqlScheme = "sqlserver",
    [String]$MsSqlUsername = "sa",
    [String]$MsSqlGatewaySpUsername = "gateway_sp",
    [String]$MsSqlGatewaySpPassword = "ThisIsReal!",
    [String]$PerceptiaDockerNet = "perceptia-net",
    [String]$GatewayPortPublish = "4443",
    [String]$RedisPort = "6379",
    [String]$RedisPortPublish = "6379",
    [String]$RedisHost = "redis",
    [String]$AqRestPort = "80",
    [String]$AqRestHost = "aqrest",
    [switch]$SkipRedis = $false,
    [switch]$RedisRemoveDbVolume = $false,
    [switch]$MsSqlRemoveDbVolume = $false,
    [switch]$SkipMsSql = $false,
    [switch]$BuildGateway = $false,
    [switch]$CleanUp = $false,
    [switch]$RemoveAllDbVolumes = $false
)

Set-Variable -Name DOCKERHUB_ORG -Value "uwthalesians"

Set-Variable -Name LATEST_COMMIT -Value "$(git rev-parse --short HEAD)"
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"

Set-Variable -Name GATEWAY_CONTAINER_NAME -Value "gateway"
Set-Variable -Name GATEWAY_SERVICE_NAME -Value "gateway-lc-svc"

Set-Variable -Name GATEWAY_TLSCERTPATH -Value "/encrypt/fullchain.pem"
Set-Variable -Name GATEWAY_TLSKEYPATH -Value "/encrypt/privkey.pem"
Set-Variable -Name GATEWAY_TLSMOUNTSOURCE -Value "$(Get-Location)/gateway/encrypt"

Set-Variable -Name GATEWAY_SESSION_KEY -Value "fjsfndreifnfsnm5kngfnklef23kdnfskng"

Set-Variable -Name REDIS_CONTAINER_NAME -Value $RedisHost
Set-Variable -Name REDIS_VOLUME_NAME -Value "${GATEWAY_SERVICE_NAME}_redis_vol"

Set-Variable -Name MSSQL_DATABASE -Value $MsSqlDatabase
Set-Variable -Name MSSQL_HOST -Value $MsSqlHost
Set-Variable -Name MSSQL_PASSWORD -Value $MsSqlGatewaySpPassword
Set-Variable -Name MSSQL_PORT -Value $MsSqlPort
Set-Variable -Name MSSQL_SCHEME -Value $MsSqlScheme
Set-Variable -Name MSSQL_USERNAME -Value $MsSqlGatewaySpUsername

Set-Variable -Name MSSQL_IMAGE_NAME -Value "mssql"
Set-Variable -Name MSSQL_CONTAINER_NAME -Value "mssql"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "${GATEWAY_SERVICE_NAME}_mssql_vol"

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
    } 
    
    if ((docker ps -aq --filter "label=label.perceptia.info/name=$GATEWAY_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=$GATEWAY_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
    }
    
    docker network create -d bridge $PerceptiaDockerNet


    
    if (!$SkipRedis) {
        Write-Host "SkipRedis option false, starting redis dependency"
    
        if ((docker ps -aq --filter "label=label.perceptia.info/name=$REDIS_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
            docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=$REDIS_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        }
        
    
        if ($RedisRemoveDbVolume -or $RemoveAllDbVolumes) {
            Write-Host "RedisRemoveDbVolume option true, clearing old sessions"
            docker volume rm $REDIS_VOLUME_NAME
        }

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
        --network $PerceptiaDockerNet `
        --publish "${RedisPortPublish}:6379" `
        --mount type=volume,source=${REDIS_VOLUME_NAME},destination=/data `
        redis:5.0.4-alpine
        Write-Host "Redis Server is listening inside docker network: ${PerceptiaDockerNet} at: ${REDIS_CONTAINER_NAME}:1433"
        Write-Host "Redis Server is listening on the host at: localhost:${MsSqlPortPublish}"
    } else {
        Write-Host "Be sure to start mssql dependency, see README"
    }
    
    if (!$SkipMsSql) {
        Write-Host "SkipMsSql option false, starting mssql dependency"
    
        if ((docker ps -aq --filter "label=label.perceptia.info/name=$MSSQL_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
            docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=$MSSQL_CONTAINER_NAME" --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
        }
        
    
        if ($MsSqlRemoveDbVolume -or $RemoveAllDbVolumes) {
            Write-Host "MsSqlRemoveDbVolume option true, removing previous database"
            docker volume rm ${MSSQL_VOLUME_NAME}
        }

        if (($MsSqlVersion).Length -eq 0) {
            Write-Host "Version must be provided, but no version provided for mssql, exiting..."
            exit(1)
        }

        Set-Variable -Name MSSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${MSSQL_IMAGE_NAME}:${MsSqlVersion}-${BUILD_AND_BRANCH}"
        Write-Host "Using mssql image: $MSSQL_IMAGE_AND_TAG"

        # Check if image exists
        (docker pull $MSSQL_IMAGE_AND_TAG) | Out-Null 
        if (!$?) {
                Write-Host "Image: $MSSQL_IMAGE_AND_TAG not found on dockerhub, exiting"
                exit(1)
        }
    
        docker run `
        --detach `
        --env 'ACCEPT_EULA=Y' `
        --env "GATEWAY_SP_USERNAME=${MSSQL_USERNAME}" `
        --env "GATEWAY_SP_PASSWORD=${MSSQL_PASSWORD}" `
        --env "MSSQL_ENVIRONMENT=development" `
        --env "SA_PASSWORD=$MsSqlSaPassword" `
        --env "SKIP_SETUP_IF_EXISTS=Y" `
        --label "label.perceptia.info/name=${MSSQL_CONTAINER_NAME}" `
        --label "label.perceptia.info/instance=${MSSQL_CONTAINER_NAME}-1" `
        --label "label.perceptia.info/managed-by=docker" `
        --label "label.perceptia.info/component=database" `
        --label "label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}" `
        --label "label.perceptia.info/sub-of=${GATEWAY_SERVICE_NAME}" `
        --mount type=volume,source=${MSSQL_VOLUME_NAME},destination=/var/opt/mssql `
        --name=${MSSQL_CONTAINER_NAME} `
        --network $PerceptiaDockerNet `
        --publish "${MsSqlPortPublish}:1433" `
        ${MSSQL_IMAGE_AND_TAG}
    
        Write-Host "MsSql Server is listening inside docker network: ${PerceptiaDockerNet} at: ${MSSQL_CONTAINER_NAME}:1433"
        Write-Host "MsSql Server is listening on the host at: localhost:${MsSqlPortPublish}"
    } else {
        Write-Host "Be sure to start mssql dependency, see README"
    }

    if (!$BuildGateway) {
        # Check if image exists
        (docker pull $GATEWAY_IMAGE_AND_TAG) | Out-Null 
        if (!$?) {
                Write-Host "Image: $GATEWAY_IMAGE_AND_TAG not found on dockerhub, exiting"
                exit(1)
        }
    }
    
    docker run `
    --detach `
    --env AQREST_HOSTNAME="$AqRestHost" `
    --env AQREST_PORT="$AqRestPort" `
    --env GATEWAY_ENVIRONMENT=development `
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
    --label "label.perceptia.info/name=gateway" `
    --label "label.perceptia.info/instance=gateway-1" `
    --label "label.perceptia.info/managed-by=docker" `
    --label "label.perceptia.info/component=server" `
    --label "label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}" `
    --name ${GATEWAY_CONTAINER_NAME} `
    --network $PerceptiaDockerNet `
    --publish "${GatewayPortPublish}:443" `
    --restart on-failure `
    --mount type=bind,source=$GATEWAY_TLSMOUNTSOURCE,target=/encrypt/,readonly `
    ${GATEWAY_IMAGE_AND_TAG}
    
    Write-Host "Gateway is listening at https://localhost:${GatewayPortPublish}"
    docker ps --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}"
} else {

    if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}").Length -gt 0) {
        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${GATEWAY_SERVICE_NAME}")
    }
    
    
    if ($RedisRemoveDbVolume -or $RemoveAllDbVolumes) {
        Write-Host "RedisRemoveDbVolume option true, clearing old sessions"
        docker volume rm $REDIS_VOLUME_NAME
    }

    if ($MsSqlRemoveDbVolume -or $RemoveAllDbVolumes) {
        Write-Host "MsSqlRemoveDbVolume option true, removing previous database"
        docker volume rm ${MSSQL_VOLUME_NAME}
    }

}

