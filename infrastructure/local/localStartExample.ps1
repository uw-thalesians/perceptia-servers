Param (
        [Switch]$CleanUp,
        [switch]$Latest,
        [string]$Build = "232",
        [String]$Branch = "develop",
        [switch]$CurrentBranch,
        [String]$GatewayVersion = "0.3.0",
        [string]$GatewayPortPublish = "4443",
        [String]$MsSqlVersion = "0.7.1",
        [String]$MsSqlPortPublish = "47011",
        [switch]$MsSqlResetDb,
        [String]$RedisPortPublish = "47012",
        [String]$AqRestVersion = "1.1.0",
        [String]$AqRestPortPublish = "47020",
        [String]$AqMySqlVersion = "1.0.0",
        [String]$AqMySqlPortPublish = "47021",
        [String]$AqSolrVersion = "1.0.0",
        [String]$AqSolrPortPublish = "47022"
)

Set-Variable -Name DOCKERHUB_ORG -Value "uwthalesians"
Set-Variable -Name PERCEPTIA_STACK_NAME -Value perceptia-api
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name MSSQL_IMAGE_NAME -Value "mssql"
Set-Variable -Name AQREST_IMAGE_NAME -Value "aqrest"
Set-Variable -Name AQMYSQL_IMAGE_NAME -Value "aqmysql"
Set-Variable -Name AQSOLR_IMAGE_NAME -Value "aqsolr"



if (!$CleanUp) {
        Write-Host "Note, this command requires docker swarm to be initialized"
        Write-Host "To initialize, run: docker swarm init"
        Write-Host "`n"
        Write-Host "Remember, you must create the Tls cert and key files in the ./encrypt/ directory"

        Write-Host "`n"
        
        if ((docker stack ls --format "{{.Name}}") -contains $PERCEPTIA_STACK_NAME) {
                Write-Host "Note, due to issue with bind points (see https://github.com/docker/for-win/issues/1521), must clean up stack before redeployment"
                Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
                docker stack rm $PERCEPTIA_STACK_NAME
                Write-Host "Waiting 10 seconds to allow docker to clean up"
                Start-Sleep -Seconds "10"
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
        } 
        if (($TAG_BUILD).Length -eq 0) {
                Write-Host "Build must be set, but no build set, exiting"
                exit(1)
        }

        Set-Variable -Name BUILD_AND_BRANCH -Value "build-${TAG_BUILD}-branch-${TAG_BRANCH}"

        # Gateway perceptia-stack.yml substitution variables
        Set-Item -Path env:GATEWAY_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${GATEWAY_IMAGE_NAME}:${GatewayVersion}-${BUILD_AND_BRANCH}"
        if (($GatewayVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for gateway, exiting..."
                exit(1)
        }
        Set-Item -Path env:GATEWAY_PORT_PUBLISH -Value $GatewayPortPublish
        # Mssql perceptia-stack.yml substitution variables
        Set-Item -Path env:MSSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${MSSQL_IMAGE_NAME}:${MsSqlVersion}-${BUILD_AND_BRANCH}"
        if (($MsSqlVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for mssql, exiting..."
                exit(1)
        }
        Set-Item -Path env:MSSQL_SKIP_SETUP_IF_EXIST -Value "Y"
        if ($MsSqlResetDb) {
                Set-Item -Path env:MSSQL_SKIP_SETUP_IF_EXIST -Value "N"
        }
        # Redis perceptia-stack.yml substituion variables
        Set-Item -Path env:REDIS_IMAGE_AND_TAG -Value "redis:5.0.4-alpine"
        Set-Item -Path env:REDIS_PORT_PUBLISH -Value $RedisPortPublish
        # Aqrest perceptia-stack.yml substitution variables
        Set-Item -Path env:AQREST_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQREST_IMAGE_NAME}:${AqRestVersion}-${BUILD_AND_BRANCH}"
        if (($AqRestVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for aqrest, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQREST_PORT_PUBLISH -Value $AqRestPortPublish
        # Aqmysql perceptia-stack.yml substitution variables
        Set-Item -Path env:AQMYSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQMYSQL_IMAGE_NAME}:${AqMySqlVersion}-${BUILD_AND_BRANCH}"
        if (($AqMySqlVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided foraqmysql, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQMYSQL_PORT_PUBLISH -Value $AqMySqlPortPublish
        # Aqsolr perceptia-stack.yml substitution variables
        Set-Item -Path env:AQSOLR_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQSOLR_IMAGE_NAME}:${AqSolrVersion}-${BUILD_AND_BRANCH}"
        if (($AqSolrVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for aqsolr, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQSOLR_PORT_PUBLISH -Value $AqSolrPortPublish

        Set-Item -Path env:MSSQL_SA_PASSWORD -Value "SoSecure!"
        Set-Item -Path env:AQMYSQL_USER_PASS -Value "8aWZjNadxspXQEHu"

        # Ensure images exist
        $ALL_IMAGES = @($env:GATEWAY_IMAGE_AND_TAG, $env:MSSQL_IMAGE_AND_TAG, $env:REDIS_IMAGE_AND_TAG, $env:AQREST_IMAGE_AND_TAG, $env:AQMYSQL_IMAGE_AND_TAG, $env:AQSOLR_IMAGE_AND_TAG)
        Write-Host "Checking dockerhub for images..."
        foreach ($IMAGE in $ALL_IMAGES) {
                (docker pull $IMAGE) | Out-Null 
                if (!$?) {
                        Write-Host "Image: $IMAGE not found on dockerhub, exiting"
                        exit(1)
                }
        }
        Write-Host "All images found!"
        Write-Host "`n"
        Write-Host "Starting up the docker stack: $PERCEPTIA_STACK_NAME"
        Write-Host "Starting stack using tag: $DOCKERHUB_ORG/{imageName}:{version-}$BUILD_AND_BRANCH"
        docker stack deploy -c perceptia-stack.yml $PERCEPTIA_STACK_NAME
        if (!$?) {
                Write-Host "Docker stack: $PERCEPTIA_STACK_NAME failed to start, see error above."
                exit(1)
        }
        Write-Host "`n"
        foreach ($IMAGE in $ALL_IMAGES) {
                Write-Host "Image used: $IMAGE"
        }
        Write-Host "`n"


        Write-Host "Perceptia backend is listening for requests at: https://localhost:$GatewayPortPublish"
        Write-Host "To test if the gateway is able to process requests, make a GET request to:"
        Write-Host "/api/v1/gateway/health"
        Write-Host "Example request using curl"
        Write-Host "curl --insecure -X GET `"https://localhost:${GatewayPortPublish}/api/v1/gateway/health`""
        Write-Host "`n"
        Write-Host "To test the status of the services in the stack run: docker stack ps $PERCEPTIA_STACK_NAME"
        Write-Host "To see the logs for a particular service run: docker service logs ${PERCEPTIA_STACK_NAME}_nameOfService"
        Write-Host "For example: docker service logs ${PERCEPTIA_STACK_NAME}_gateway"
        Write-Host "`n"
        Write-Host "Mssql server for the Perceptia databased used by the gateway can be reached at port: $Env:MSSQL_PORT_PUBLISH"
        Write-Host "Perceptia databased used by the gateway has an 'sa' user account with password: $Env:MSSQL_SA_PASSWORD"
        Write-Host "Aqmysql server for the any_quiz_db databased used by the aqrest service can be reached at port: $Env:AQMYSQL_PORT_PUBLISH"
        Write-Host "any_quiz_db databased used by theaqrest service has an 'any_quiz' user account with password: $Env:AQMYSQL_USER_PASS"
        Write-Host "`n"
        docker stack ps $PERCEPTIA_STACK_NAME

} else {
        Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
        docker stack rm $PERCEPTIA_STACK_NAME
}

