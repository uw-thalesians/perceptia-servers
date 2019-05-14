# Local

The purpose of this direcotry is to organize files that support the local running of the perceptia backend. These files are not a part of the actual application.

## Contents

* [Setup](#setup)

* [Structure](#structure)

* [Run Local Start Script](#run-local-start-script)

## Setup

## Structure

The following list describes the purpose of the directories and files in this subdirectory:

[./encrypt/](./encrypt) contains a script to generate a self signed tls cert for the gateway

[./localStartExample.ps1](./localStartExample.ps1) is a PowerShell script to start the backend using docker stack

[./perceptia-stack.yml](./perceptia-stack.yml) is the docker compose file used to define the docker stack deployment

## Run Local Start Script

The local start script [./localStartExample.ps1](./localStartExample.ps1) wraps the docker stack commands to start and stop the Perceptia backend. See the [docker stack config section below](#docker-stack-config) for how to connect to the backend one its running.

Note, you will need to create the Tls certs using the script in the encrypt directory before running the local start script for the first time.

1. Read the [README in ./encrypt/](./encrypt/README.md) and follow the instructions there for running the createTlsCert.sh script to generate the Tls certificate and private key. These files will be used by the backend to accept requests using Tls (secure) connection.

2. Run `docker swarm init` if you haven't already started or attached a swarm master

3. If the stack is already running (you can run `docker stack ls` to see running deployments), you should run the localStatExample.ps1 script with the -CleanUp flag before running it again. The stack does not handle in place upgrades (rerunning the deployment) well. Once you run the script with the -CleanUp flag, you should give it at least 10 seconds to clean everything up, otherwise the deployment might error out if run too soon. (If it does error out, just wait a few more seconds before running again)

   `localStartExample.ps1 -CleanUp`

4. Run the localStatExample.ps1 script, see below:

To start the backend run the script with no options:

`./localStartExample.ps1`

To stop the backend, run the script with the option `-CleanUp`

`./localStartExample.ps1 -CleanUp`

If you want to run the latest tagged image for each service built from the develop branch, then add the -Latest switch

`./localStartExample.ps1 -Latest`

Finally, if you want to run the latest build from your current branch, use the -Branch option to specify the name of your branch after the feature part. For example, if your branch is "feature/peacock-local-start"

`./localStartExample.ps1 -Latest -Branch peacock-local-start`

### Comand Line Options

`-CleanUp` (Switch) when used, stops the services started by docker stack deploy

`-Latest` (Switch) when used, starts the stack using the latest images for the given version of each image built from the develop branch

`-Build` (String) specify the build number to use image builds from, default is a known working build for all images used, will be ignored if -Latest is also set

`-Branch` (String) specify the branch to use image builds from, default is "develop"

`-GatewayPortPublish` (String) sets the port that requests can be made to the gateway on, default is "4443"

`-MsSqlPortPublish` (String) sets the port that the mssql service can recieve requests on, default is "47011" which maps to "1433" inside the container (Note, the mssql service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-AqRestPortPublish` (String) sets the port that the aqrest service can recieve requests on, default is "47020" which maps to "80" inside the contianer (Note, the aqrest service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-AqMySqlPortPublish` (String) sets the port that the aqmysql service can recieve requests on, default is "47021" which maps to "3306" inside the contianer (Note, the aqmysql service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-AqSolrPortPublish` (String) sets the port that the aqsolr service can recieve requests on, default is "47022" which maps to "8983" inside the contianer (Note, the aqsolr service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

### Docker Stack Config

The docker stack config makes certain assumptions about the images being run and the local ip environemnt on the device where the script is being run.

The application backend is listening for requests on the localhost using https at the port 4443.

Scheme: `https`

Host: `localhost`

Port: `4443`

Additionally, each service is also exposed to make it easier to interact with the service for troublshooting. These ports can be customized using the options listed above.

### Checking the status of a deployment

To test the status of the services in the stack run: `docker stack ps perceptia-api`

To see the logs for a particular service run: `docker service logs perceptia-api_nameOfService`

For example: `docker service logs perceptia-api_gateway`