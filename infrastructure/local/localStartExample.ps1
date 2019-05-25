Param (
        [switch]$Latest,
        [string]$Build = "326",
        [String]$Branch = "develop",
        [switch]$CurrentBranch,

        [String]$GatewayVersion = "1",
        [string]$GatewayPortPublish = "4443",

        [String]$MsSqlVersion = "1",
        [String]$MsSqlSaPassword = "SecureNow!",
        [String]$MsSqlPortPublish = "47011",
        [String]$MsSqlGatewaySpUsername = "gateway_sp",
        [String]$MsSqlGatewaySpPassword = "ThisIsReal!",
        [switch]$MsSqlRemoveDbVolume,

        [String]$RedisPortPublish = "47012",
        [switch]$RedisRemoveDbVolume,

        [String]$AqRestVersion = "1.1.0",
        [String]$AqMySqlUserPassword = "8aWZjNadxspXQEHu",
        [String]$AqRestPortPublish = "47020",

        [String]$AqMySqlVersion = "1.0.0",
        [String]$AqMySqlPortPublish = "47021",
        [switch]$AqMySqlRemoveDbVolume,

        [String]$AqSolrVersion = "1.0.0",
        [String]$AqSolrPortPublish = "47022",

        [switch]$RemoveAllDbVolumes,
        [switch]$SkipImageCheck,
        [switch]$RemoveAllContainers,
        [Switch]$CleanUp
)

# Setup Base Veriables
Set-Variable -Name DOCKERHUB_ORG -Value "uwthalesians"

Set-Variable -Name PERCEPTIA_STACK_NAME -Value "perceptia-api"

## Gateway Veriables
Set-Variable -Name GATEWAY_IMAGE_NAME -Value "gateway"
Set-Variable -Name GATEWAY_API_PORT -Value "$GatewayPortPublish"

## Redis Variables
Set-Variable -Name REDIS_VOLUME_NAME -Value "redis_pc_vol"

## MsSql Variables
Set-Variable -Name MSSQL_IMAGE_NAME -Value "mssql"
Set-Variable -Name MSSQL_VOLUME_NAME -Value "mssql_pc_vol"

## AqRest Variables
Set-Variable -Name AQREST_IMAGE_NAME -Value "aqrest"

## AqMySql Variables
Set-Variable -Name AQMYSQL_IMAGE_NAME -Value "aqmysql"
Set-Variable -Name AQMYSQL_VOLUME_NAME -Value "aqmysql_pc_vol"

## AqSolr Variables
Set-Variable -Name AQSOLR_IMAGE_NAME -Value "aqsolr"



if (!$CleanUp) {
        Write-Host "Note, this script requires docker swarm to be initialized"
        Write-Host "To initialize, run: docker swarm init"
        Write-Host "`n"
        Write-Host "Remember, you must create the Tls cert and key files in the ./encrypt/ directory"

        Write-Host "`n"

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

        # Setup Environment Variables
        Set-Item -Path env:PERCEPTIA_STACK_NAME -Value $PERCEPTIA_STACK_NAME

        ## Gateway perceptia-stack.yml substitution variables
        Set-Item -Path env:GATEWAY_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${GATEWAY_IMAGE_NAME}:${GatewayVersion}-${BUILD_AND_BRANCH}"
        if (($GatewayVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for gateway, exiting..."
                exit(1)
        }
        Set-Item -Path env:GATEWAY_PORT_PUBLISH -Value $GatewayPortPublish
        Set-Item -Path env:GATEWAY_API_PORT -Value $GATEWAY_API_PORT


        ## Mssql perceptia-stack.yml substitution variables
        Set-Item -Path env:MSSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${MSSQL_IMAGE_NAME}:${MsSqlVersion}-${BUILD_AND_BRANCH}"
        if (($MsSqlVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for mssql, exiting..."
                exit(1)
        }
        Set-Item -Path env:MSSQL_SA_PASSWORD -Value $MsSqlSaPassword
        Set-Item -Path env:MSSQL_GATEWAY_SP_PASSWORD -Value $MsSqlGatewaySpPassword
        Set-Item -Path env:MSSQL_GATEWAY_SP_USERNAME -Value $MsSqlGatewaySpUsername
        if ($GatewayVersion -eq "0.3.0") {
                Set-Item -Path env:MSSQL_GATEWAY_SP_PASSWORD -Value $MsSqlSaPassword
                Set-Item -Path env:MSSQL_GATEWAY_SP_USERNAME -Value "sa"
        }
        Set-Item -Path env:MSSQL_PORT_PUBLISH -Value $MsSqlPortPublish


        ## Redis perceptia-stack.yml substituion variables
        Set-Item -Path env:REDIS_IMAGE_AND_TAG -Value "redis:5.0.4-alpine"
        Set-Item -Path env:REDIS_PORT_PUBLISH -Value $RedisPortPublish


        ## Aqrest perceptia-stack.yml substitution variables
        Set-Item -Path env:AQREST_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQREST_IMAGE_NAME}:${AqRestVersion}-${BUILD_AND_BRANCH}"
        if (($AqRestVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for aqrest, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQREST_PORT_PUBLISH -Value $AqRestPortPublish


        ## Aqmysql perceptia-stack.yml substitution variables
        Set-Item -Path env:AQMYSQL_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQMYSQL_IMAGE_NAME}:${AqMySqlVersion}-${BUILD_AND_BRANCH}"
        if (($AqMySqlVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided foraqmysql, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQMYSQL_PORT_PUBLISH -Value $AqMySqlPortPublish
        Set-Item -Path env:AQMYSQL_USER_PASS -Value $AqMySqlUserPassword


        ## Aqsolr perceptia-stack.yml substitution variables
        Set-Item -Path env:AQSOLR_IMAGE_AND_TAG -Value "${DOCKERHUB_ORG}/${AQSOLR_IMAGE_NAME}:${AqSolrVersion}-${BUILD_AND_BRANCH}"
        if (($AqSolrVersion).Length -eq 0) {
                Write-Host "Version must be provided, but no version provided for aqsolr, exiting..."
                exit(1)
        }
        Set-Item -Path env:AQSOLR_PORT_PUBLISH -Value $AqSolrPortPublish







        # Remove stack to redeploy        
        if ((docker stack ls --format "{{.Name}}") -Match $PERCEPTIA_STACK_NAME) {
                Write-Host "`n"
                Write-Host "Note, due to issue with bind points (see https://github.com/docker/for-win/issues/1521), must clean up stack before redeployment"
                Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
                docker stack rm $PERCEPTIA_STACK_NAME
                Write-Host "Waiting 15 seconds to allow docker to clean up"
                Start-Sleep -Seconds "15"
        }

        if ($RemoveAllDbVolumes -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}"))) {
                Write-Host "Database volume reset requested"
                if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        Write-Host "Removing all containers before removing volumes"
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
                docker volume rm (docker volume ls --format "{{.Name}}" --filter "name=${PERCEPTIA_STACK_NAME}")
        }
        if (($MsSqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME" ))){
                Write-Host "-MsSqlRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=mssql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=mssql"  --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                } 
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME
        }
        if (($RedisRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME" ))) {
                Write-Host "-RedisRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=redis" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=redis"  --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                } 
                
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME
        }
        if (($AqMySqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME" ))) {
                Write-Host "-AqMySqlRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=aqmysql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=aqmysql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
                
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME
        }
        if ($RemoveAllContainers) {
                if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
        }

        # Ensure images exist
        if (!$SkipImageCheck) {
                $ALL_IMAGES = @($env:GATEWAY_IMAGE_AND_TAG, $env:MSSQL_IMAGE_AND_TAG, $env:REDIS_IMAGE_AND_TAG, $env:AQREST_IMAGE_AND_TAG, $env:AQMYSQL_IMAGE_AND_TAG, $env:AQSOLR_IMAGE_AND_TAG)
                Write-Host "Checking dockerhub for images..."
                foreach ($IMAGE in $ALL_IMAGES) {
                        (docker pull $IMAGE) | Out-Null 
                        if (!$?) {
                                Write-Host "Image: $IMAGE not found on dockerhub, exiting"
                                exit(1)
                        }
                }
        }
       
        Write-Host "All images found!"
        Write-Host "`n"
        Write-Host "Starting up the docker stack: $PERCEPTIA_STACK_NAME"
        Write-Host "Starting stack using tag: $DOCKERHUB_ORG/{imageName}:{version-}$BUILD_AND_BRANCH"
        docker stack deploy -c perceptia-stack.yml $PERCEPTIA_STACK_NAME
        if (!$?) {
                Write-Host "Docker stack: $PERCEPTIA_STACK_NAME failed to start, see error above."
                Write-Host "`n"
                Write-Host "In most cases rerunning this script a second time will resolve error."
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
        if ((docker stack ls --format "{{.Name}}") -Match $PERCEPTIA_STACK_NAME) {
                Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
                docker stack rm $PERCEPTIA_STACK_NAME
                Write-Host "Waiting 5 seconds to allow docker to clean up"
                Start-Sleep -Seconds "5"
        }
        if ($RemoveAllDbVolumes -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}"))) {
                Write-Host "Database volume reset requested"
                if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        Write-Host "Removing all containers before removing volumes"
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
                docker volume rm (docker volume ls --format "{{.Name}}" --filter "name=${PERCEPTIA_STACK_NAME}")
        }
        if (($MsSqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME" ))){
                Write-Host "-MsSqlRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=mssql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=mssql"  --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                } 
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$MSSQL_VOLUME_NAME
        }
        if (($RedisRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME" ))) {
                Write-Host "-RedisRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=redis" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=redis"  --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                } 
                
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$REDIS_VOLUME_NAME
        }
        if (($AqMySqlRemoveDbVolume) -and (((docker volume ls --format "{{.Name}}") -Match "${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME" ))) {
                Write-Host "-AqMySqlRemoveDbVolume option set, removing volume: ${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME"
                if ((docker ps -aq --filter "label=label.perceptia.info/name=aqmysql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/name=aqmysql" --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
                
                Start-Sleep -Seconds 2
                docker volume rm ${PERCEPTIA_STACK_NAME}_$AQMYSQL_VOLUME_NAME
        }
        if ($RemoveAllContainers) {
                if ((docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")) {
                        docker rm --force (docker ps -aq --filter "label=label.perceptia.info/part-of=${PERCEPTIA_STACK_NAME}")
                }  
        }
        Write-Host "Docker stack should be cleaned up..."
}

