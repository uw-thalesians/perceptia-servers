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

The local start script [./localStartExample.ps1](./localStartExample.ps1) wraps the docker stack commands to start and stop the Perceptia backend. See the [docker stack config section below](#docker-stack-config) for how to connect to the backend once its running.

Note, you will need to create the Tls certs using the script in the encrypt directory before running the local start script for the first time.

1. Read the [README in ./encrypt/](./encrypt/README.md) and follow the instructions there for running the createTlsCert.sh script to generate the Tls certificate and private key. These files will be used by the backend to accept requests using a Tls (secure) connection.

2. Run `docker swarm init` if you haven't already started or attached a swarm master.

        `docker swarm init`

3. Ensure Google Api Key is provided: either set the PERCEPTIA_AQREST_GOOGLE_API_KEY environment variable, or provide the key using the -AqRestGoogleApiKey option when running the script (see step 4 below)

4. Run the localStatExample.ps1 script, see below:

        To start the backend run the script with no options:

        `./localStartExample.ps1`

        If you want to run the latest tagged image for each service built from the develop branch, then add the -Latest switch

        `./localStartExample.ps1 -Latest`

        Finally, if you want to run the latest build from your current branch, use the -Branch option to specify the name of your branch after the feature part. For example, if your branch is "feature/peacock-local-start". Note, if you have incremented the version of your service being built in that branch you can use the version options to specify the specific version for that image to pull in conjunction with the -Branch or -CurrentBranch options

        `./localStartExample.ps1 -Latest -Branch peacock-local-start`

        `./localStartExample.ps1 -Latest -CurrentBranch`

5. To remove the deployment (stop the containers) use the `-CleanUp` flag

        To stop the backend, run the script with the option `-CleanUp`

        `./localStartExample.ps1 -CleanUp`

        To stop the backend and remove all database volumes run the script with the options `-CleanUp` and `-ResetAllDbVolumes`

        `./localStartExample.ps1 -CleanUp -RemoveAllDbVolumes`

### Comand Line Options

`-CleanUp` (Switch) when used, stops the services started by docker stack deploy, default false

`-Latest` (Switch) when used, starts the stack using the latest images for the given version of each image built from the develop branch, default false

`-Build` (String) specify the build number to use image builds from, default is a known working build for all images used, will be ignored if -Latest is also set

`-Branch` (String) specify the branch to use image builds from, default is "develop"

`-CurrentBranch` (Switch) when set, uses the name of the current branch to specify the images to use, if on branch "feature/peacock-local-start" would use images with the branch tag "peacock-local-start", default false

`-GatewayVersion` (String) sets the version of the gateway image to use, default is a known stable version of the image

`-GatewayPortPublish` (String) sets the port that requests can be made to the gateway on, default is "4443"

`-RedisPortPublish` (String) sets the port that the redis service can recieve requests on, default is "47012" which maps to "6379" inside the container (Note, the redis service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-RedisRemoveDbVolume` (switch) when set, will remove the volume the database was created in, forcing it to be rebuilt, default false, meaning the database will not reset if it already exists, but if it doesn't it will be created

`-MsSqlVersion` (String) sets the version of the mssql image to use, default is a known stable version of the image

`-MsSqlSaPassword` (string) which is used to pass in either the password to use to secure the mssql server, the default value is: "SecureNow!"

`-MsSqlGatewaySpUsername` (string) which is the username the gateway will use to connect to the Perceptia database as

`-MsSqlGatewaySpPassword` (string) which is the password for the user the gateway will connect to the Perceptia database as

`-MsSqlPortPublish` (String) sets the port that the mssql service can recieve requests on, default is "47011" which maps to "1433" inside the container (Note, the mssql service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-MsSqlRemoveDbVolume` (switch) when set, will remove the volume the database was created in, forcing it to be rebuilt, default false, meaning the database will not reset if it already exists, but if it doesn't it will be created

`-AqRestVersion` (String) sets the version of the aqrest image to use, default is a known stable version of the image

`-AqRestGoogleApiKey` (String) sets the api key to be used by aqrest to find images for quizzes, using the Google Custom Search API, if this value is not set, will check for the environment variable PERCEPTIA_AQREST_GOOGLE_API_KEY, if neither are found the script will exit. Setting the AqRestGoogleApiKey will ignore whatever is set for the environment variable

`-AqRestPortPublish` (String) sets the port that the aqrest service can recieve requests on, default is "47020" which maps to "80" inside the contianer (Note, the aqrest service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-AqMySqlVersion` (String) sets the version of the aqmysql image to use, default is a known stable version of the image

`-AqMySqlUserPassword` (String) sets the password used to connect to the any_quiz_db database.

`-AqMySqlPortPublish` (String) sets the port that the aqmysql service can recieve requests on, default is "47021" which maps to "3306" inside the contianer (Note, the aqmysql service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-AqMySqlRemoveDbVolume` (switch) when set, will remove the volume the database was created in, forcing it to be rebuilt, default false, meaning the database will not reset if it already exists, but if it doesn't it will be created

`-AqSolrVersion` (String) sets the version of the aqsolr image to use, default is a known stable version of the image

`-AqSolrPortPublish` (String) sets the port that the aqsolr service can recieve requests on, default is "47022" which maps to "8983" inside the contianer (Note, the aqsolr service is exposed to make it easier to make a direct connection to the container, and is not necessary for the service to function properly)

`-RemoveAllDbVolumes` (switch) will cause all volumes used by services labeled as "app.perceptia.info/component=database" to be deleted, will remove the containers that mounted the volumes before removing the volumes, default is false, meaning if volumes exist they will not be removed

`-RemoveAllContainers` (switch) will cause all containers used by the stack labeled as "app.perceptia.info/part-of=perceptia-api" to be deleted

`-SkipImageCheck` (switch) will skip checking and pulling the images before starting the stack (note, if stack can't find images the stack will succeed but the containers won't be started because the image could not be found, so use this only if you know the branches exist)

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