# API Gateway Service

Last updated: 2019-04-30

The API Gateway service serves as the primary entry point for the Perceptia application. The Gateway provides several key services in addition to routing requests to the responsible microservice. These services include: CORS middleware, Session Authentication, Sign on, Account creation and management.

## [Contents](#contents)

* [Getting Started](#getting-started)

* [Structure](#structure)

  * [Config and Setup Files](#config-and-setup-files)

  * [Gateway Source](#gateway-source)

* [Setup Server](#setup-server)

  * [Building the Image](#building-the-image)

  * [Custom Image](#custom-image)

    * [Image Specific Options](#image-specific-options)

* [Start Server Locally](#start-server-locally)

  * [Start with Script](#start-with-script)

  * [Start with Docker Commands](#start-with-docker-commands)

* [Testing](#testing)
  
## [Getting Started](#getting-started)

The Gateway service is designed to run within a linux container. This README will describe the key files used to build and run this service in a container. Additionally, there are certain environment variables that the application expects to be present in order to run. These environment variables will also be described in this document.

## [Structure](#structure)

The root of the gateway directory contains the supporting files for building the application.

### [Config and Setup Files](#config-and-setup-files)

[Dockerfile:](./Dockerfile) multi-stage docker file to build the gateway executable and the gateway image

[.dockerignore:](./.dockerignore) identifies which files should be accessible to commands in the Dockerfile

[.gitignore:](./.gitignore) identifies which files in the gateway directory should not be tracked by git

[gateway-service-api.yaml:](./gateway-service-api.yam) documents the public REST based APIs provided by the gateway service directly. For specific versions, see the [api directory](./../api/).

[localStartExample.ps1:](./localStartExample.ps1) is meant for local testing of the gateway in a docker container

[testGatewayUnit.ps1:](./testGatewayUnit.ps1) is meant for runing the gateway unit tests locally with coverage

### [Gateway Source](#gateway-source)

The [gateway](./gateway/) directory contains the source files for the gateway executable.

[main.go:](./gateway/main.go) the source code containing the main function (entrypoint) for the gateway service

[go.mod:](./gateway/go.mod) file containing the modules used by the application and its dependencies. Used by the go command line tools to identify which packages and their specific versions to retrieve when building the gateway service into an executable

[go.sum:](./gateway/go.sum) used to track and ensure validity of retrieved package files listed in go.mod

## [Setup Server](#setup-server)

The gateway executable is designed to be deployed using a linux container. The following subsections explain how this container is built and how to use it. The gateway executable is currently built on the [GoLang](https://hub.docker.com/_/golang) image `golang:1.12`, and the gateway image is then based off the [Alpine Linux](https://hub.docker.com/_/alpine) image `alpine:3.9`.

### [Building the Image](#building-the-image)

Builds of this container image are automatically triggered by pushes to the GitHub repository.

Builds are tagged based on the version of the API the gateway implements (as defined in an variable in the azure-pipelines.yml file in the root of this repository which should reflect the API version listed in the gateway-service-api.yaml file in this directory). For a complete description of the possible tags see the [gateway container repository](https://hub.docker.com/r/uwthalesians/gateway) on the container registry DockerHub.

#### [Build](#build)

The image can be built locally using the docker build command. This command should be run from this directory (where the Dockerfile is located). See the local start script for an automated build and run.

Example docker build command: `docker build --tag "${GATEWAY_IMAGE_AND_TAG}" --no-cache .`

Commands:

  `--tag` is the name the image will be given (the name used to then run the built image as a container)

  `--no-cache` ensures that source code is always updated in container

  `.` the final period in the command indicates the root directory to send to the docker deamon for the build process, this should be the directory where the Dockerfile is located

### [Custom Image](#custom-image)

The Gateway image will be used during development and production. Information about this custom image can be found in the Thalesians container registry on DockerHub [uwthalesians/gateway](https://hub.docker.com/r/uwthalesians/gateway).

Please refer to the description on the [container registry](https://hub.docker.com/r/uwthalesians/gateway) for specifics on how to configure it. The information below only provides an exmaple setup.

#### [Image Specific Options](#image-specific-options)

This section list any configuration options for the custom image.

##### [Container Environment Variables](#custom-image-env-vars)

Environment variables must be passed to be accessible to the gateway executable (such as using the --env flag with the docker run commnad)

Use the following variables to configure the gateway for the given environment.

`GATEWAY_LISTEN_ADDR=[[<host>]:[<port>]]` (OPTIONAL) identifies what host and port the gateway should listen for requests on. If this variable is not set the gateway will default to ":443".

`GATEWAY_TLSCERTPATH=<pathToCert>` (REQUIRED) identifies the absolute path to the certificate file to be used by the gateway to make TLS connections. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container

`GATEWAY_TLSKEYPATH=<pathToCertKey>` (REQUIRED) identifies the absolute path to the key file for the certificate identified by the "GATEWAY_TLSCERTPATH" variable. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container

`GATEWAY_SESSION_KEY=<sessionkey>` (REQUIRED) the session key used to sign login sessions

`MSSQL_SCHEME=<scheme>` (REQUIRED) identifies the scheme to use to connect to the mssql database

`MSSQL_USERNAME=<username>` (REQUIRED) identifies the username to login to the mssql database with

`MSSQL_PASSWORD=<password>` (REQUIRED) the password used to login to the mssql database

`MSSQL_HOST=<host>` (REQUIRED) the host where the mssql database can be reached

`MSSQL_PORT=<port>` (REQUIRED) the port where the mssql database can be reached

`MSSQL_DATABASE=<database>` (REQUIRED) the database to use for the connection

`AQREST_HOSTNAME=<hostname>` (REQUIRED) the hostname of the aqrest service

`AQREST_PORT=<port>` (REQUIRED) the port that the aqrest service is listening on

`REDIS_ADDRESS=<hostname:port>` (REQUIRED) the hostname and port the redis server is listening on

`GATEWAY_API_PORT={port}` (optional) identifies the external port that clients reach the gateway from, default 443

`GATEWAY_API_HOST={hostname}` (optional) identifies the external hostname that clients reach the gateway from, default localhost

`GATEWAY_API_SCHEME={scheme}` (optional) identifies the external scheme that clients reach the gateway from, default https

## [Start Server Locally](#start-server-locally)

This setup explains how to build and start the server locally.

### [Start with Script](#start-with-script)

Building and starting the gateway container locally can be more involed than running one script. The following must be setup:

1. Populate the [./gateway/encrypt](./gateway/encrypt) directory with the apprioriate TLS certificates for localhost. See [./gateway/encrypt/README.md](./gateway/encrypt/README.md) for instructions.

2. Ensure service dependencies are set up to run. If using the PowerShell script, the dependencies should be started by default. (Meaing, if you are using the PowerShell script with the default options you should not need to start the dependencies)

   * [mssql](https://hub.docker.com/r/uwthalesians/mssql): image with version 0.7.1 bootstraped (see [local start script](./../database/mssql/localStartExample.ps1) for example)

   * [redis](https://hub.docker.com/_/redis): redis:5.0.4-alpine (see [local start script](./../database/mssql/localStartExample.ps1) for example)

#### PowerShell

For testing the gateway locally, the [localStartExample.ps1](./localStartExample.ps1) script can be used. This script assumes that docker is already installed and running on the system and that the TLS cert and key have been generated (see note above). Note, the script is a PowerShell script and thus requires a PowerShell shell. Additionally, PowerShell will not run unsigned scripts by default, therefore you [may need to enable running unsigned scripts](https://superuser.com/questions/106360/how-to-enable-execution-of-powershell-scripts) to use it.

The PowerShell script, [localStartExample.ps1](./localStartExample.ps1) will run the gateway image as a container inside a docker network and expose it to localhost. This script has several command line options which allow you to customize the instance. By default this script will also start the gateway dependencies (redis and mssql servers).

##### Command Line Options

The script accepts several comand line options which can be set when running the script in a PowerShell terminal. No positional options.

Unless you need to run your own mssql container or build the gateway container locally, you should not have to provide any options to the local start script.

    Run: `.\locaStartExample.ps1`

However, **if you want to remove redis and mssql databases between runs of the containers**, you need to include the -RedisRemoveDbVolume and -MsSqlRemoveDbVolume switch parameters or the -RemoveAllDbVolumes parameter.

    Run: `.\locaStartExample.ps1 -RedisRemoveDbVolume -MsSqlRemoveDbVolume`

To clean up (take down the containers started by running this script):

    Run: `.\locaStartExample.ps1 -CleanUp -RemoveAllDbVolumes`

`-Latest` (Switch) when used, starts the stack using the latest images for the given version of each image built from the develop branch, default false

`-Build` (String) specify the build number to use image builds from, default is a known working build for all images used, will be ignored if -Latest is also set

`-Branch` (String) specify the branch to use image builds from, default is "develop"

`-CurrentBranch` (Switch) when set, uses the name of the current branch to specify the images to use, if on branch "feature/peacock-local-start" would use images with the branch tag "peacock-local-start", default false

`-GatewayVersion` (String) sets the version of the gateway image to use, default is a known stable version of the image

`-GatewayPortPublish` which is the port the gateway service should be exposed on the host machine, default value is "4443"

`-BuildGateway` (switch) will build the gateway using the local source, default is: false. To set true, include the switch

`-MsSqlVersion` (String) sets the version of the mssql image to use, default is a known stable version of the image

`-MsSqlDatabase` (string) which is used to pass in the name of the database to use for connections to the server, default value is: "Perceptia"

`-MsSqlHost` (string) which is used to pass in the hostname of the mssql database server, default value is: "mssql"

`-MsSqlSaPassword` (string) which is used to pass in the password to use to secure the mssql server, the default value is: "SecureNow!"

`-MsSqlPort` (string) which is the port the docker container should listen for requests on and send to the mssql server, default value is: "1433",

`-MsSqlPortPublish` (string) which is the port to publish the mssql container if run by this script, default is: "1401"

`-MsSqlScheme` (string) which is the URL type scheme used to connect to the mssql database server, default value is: "sqlserver"

`-MsSqlUsername` (string) which is the username used to connect to the mssql database server, default value is: "sa"

`-MsSqlGatewaySpUsername` (string) which is the username the gateway will use to connect to the Perceptia database as

`-MsSqlGatewaySpPassword` (string) which is the password for the user the gateway will connect to the Perceptia database as

`-PerceptiaDockerNet` (string) which is the name of the docker network the container should be attached to when run, default value is "perceptia-net"

`-RedisPort` which is the port the gateway should reach the redis container at, default is: "6379

`-RedisPortPublish` which is the port to publish the redis container if run by this script, default is: "6379"

`-SkipRedis` (switch) will skip starting the redis dependency, default is: false. If you are starting your own redis instance include the option by including the switch

`-SkipMsSql` (switch) will skip starting the mssql dependency, default is: false. If you are starting your own mssql instance, set option by including the switch

`-MsSqlRemoveDbVolume` (switch) will remove any existing named volume used by the mssql service, default is false

`-RedisRemoveDbVolume` (switch) will remove any existing named volume used by the redis service, default is false

`-RemoveAllDbVolumes` (switch) will remove any existing named volume used by the redis or mssql service, default is false

`-CleanUp` (switch) will remove any existing container started by this script

### [Start with Docker Commands](#start-with-docker-commands)

For directions to start the container locally using a script, see [Start Server Locally](#start-server-locally).

1. pull the image from docker (check [registry](https://hub.docker.com/r/uwthalesians/gateway) for latest images)

   `docker pull uwthalesians/0.2.0-build-latest-branch-develop`

2. run the container image (replace variables with the correct values)

   See end of [localStartExample.ps1](./localStartExample.ps1) for example docker run

## [Testing](#testing)

### [Go Test Command](#testing-go-test)

Go build tags are used to identify test types and allow selectively running tests, such as unit tests and integration tests. Testing files have the go build directive options to identify which build tags to run the test for:

`// +build tag_example all unit etc`

With this build directive at the top of the `some_test.go` file, followed by a blank line, when go test is run only the explicit tags provided to go test that match a build directive tag will be run. Example:

`go test -tags=unit ./...`

This command (assuming it is run from the same directory as the root go.mod file) will run all test files that contain the `unit` build tag.

### [Go Unit Test Scripts](#testing-go-unit-test-scripts)

The [./testGatewayUnit.ps1](./testGatewayUnit.ps1) script will run all unit tests for the gateway source code, with coverage, and load a webpage with the results of the coverage test.