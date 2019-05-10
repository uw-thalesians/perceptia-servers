Param (
        [Switch]$CleanUp
)

Set-Variable -Name PERCEPTIA_STACK_NAME -Value perceptia-api

Write-Host "Note, this command requires docker swarm to be initialized"

if (!$CleanUp) {
        Write-Host "Remember, you must create the Tls cert and key files in the ./encrypt/ directory"

        Write-Host "Starting up the docker stack: $PERCEPTIA_STACK_NAME"

        docker stack deploy -c perceptia-stack.yml $PERCEPTIA_STACK_NAME
        Write-Host "Perceptia backend is listening for requests at: https://localhost:4443"
        Write-Host "To test if the gateway is able to process requests, make a GET request to:"
        Write-Host "/api/v1/gateway/health"
        Write-Host "Example request using curl"
        Write-Host 'curl -X GET "https://localhost:4443/api/v1/gateway/health"'

} else {
        Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
        docker stack rm $PERCEPTIA_STACK_NAME
}

