Param (
        [Switch]$CleanUp,
        [switch]$Latest,
        [string]$Build = "227",
        [String]$Branch = "develop",
        [string]$GatewayPortPublish = "4443",
        [String]$MsSqlPortPublish = "47011",
        [String]$AqRestPortPublish = "47020",
        [String]$AqMySqlPortPublish = "47021",
        [String]$AqSolrPortPublish = "47022"
)

Set-Variable -Name PERCEPTIA_STACK_NAME -Value perceptia-api

Write-Host "Note, this command requires docker swarm to be initialized"
Write-Host "docker swarm init"

if (!$CleanUp) {
        Write-Host "Remember, you must create the Tls cert and key files in the ./encrypt/ directory"

        Write-Host "Starting up the docker stack: $PERCEPTIA_STACK_NAME"

        # Define Image Tags to use
        Set-Item -Path env:TAG_BUILD -Value $Build # Build number 
        if ($Latest) {
                Write-Host "Latest switch provided, building latest build from"
                Set-Item -Path env:TAG_BUILD -Value "latest"              
        } 

        Set-Item -Path env:TAG_BRANCH -Value $Branch
        Set-Item -Path env:GATEWAY_PORT_PUBLISH -Value $GatewayPortPublish
        Set-Item -Path env:MSSQL_PORT_PUBLISH -Value $MsSqlPortPublish
        Set-Item -Path env:AQREST_PORT_PUBLISH -Value $AqRestPortPublish
        Set-Item -Path env:AQMYSQL_PORT_PUBLISH -Value $AqMySqlPortPublish
        Set-Item -Path env:AQSOLR_PORT_PUBLISH -Value $AqSolrPortPublish
        
        Write-Host "Starting stack using tag: uwthalesians/{imageName}:{version}-$env:TAG_BUILD-$env:TAG_BRANCH"
        docker stack deploy -c perceptia-stack.yml $PERCEPTIA_STACK_NAME


        if (!$?) {
                Write-Host "Docker stack: $PERCEPTIA_STACK_NAME failed to start, see error above."
                exit(1)
        }
        Write-Host "Perceptia backend is listening for requests at: https://localhost:$GatewayPortPublish"
        Write-Host "To test if the gateway is able to process requests, make a GET request to:"
        Write-Host "/api/v1/gateway/health"
        Write-Host "Example request using curl"
        Write-Host 'curl --insecure -X GET "' + "https://localhost:${GatewayPortPublish}/api/v1/gateway/health" + '"'
        Write-Host "To test the status of the services in the stack run: docker stack ps $PERCEPTIA_STACK_NAME"
        Write-Host "To see the logs for a particular service run: docker service logs ${PERCEPTIA_STACK_NAME}_nameOfService"
        Write-Host "For example: docker service logs ${PERCEPTIA_STACK_NAME}_gateway"

} else {
        Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
        docker stack rm $PERCEPTIA_STACK_NAME
}

