Param (
        [Switch]$CleanUp
)

Set-Variable -Name PERCEPTIA_STACK_NAME -Value perceptia-api

if (!$CleanUp) {
        Write-Host "Remember, you must create the Tls cert and key files in the ./encrypt/ directory"

        Write-Host "Starting up the docker stack: $PERCEPTIA_STACK_NAME"

        docker stack deploy -c perceptia-stack.yml $PERCEPTIA_STACK_NAME
} else {
        Write-Host "Cleaning up the docker stack: $PERCEPTIA_STACK_NAME"
        docker stack rm $PERCEPTIA_STACK_NAME
}

