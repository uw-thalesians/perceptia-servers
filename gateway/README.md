# API Gateway Service

Last updated: 2019-04-30

The API Gateway service serves as the primary entry point for the Perceptia application. The Gateway provides several key services in addition to routing requests to the responsible microservice. These services include: CORS middleware, Session Authentication, Sign on, Account creation and management.

## [Contents](#contents)

* [Getting Started](#getting-started)

* [Structure](#structure)

  * [Config and Setup Files](#structure-files)

  * [Gateway Source](#structure-gateway-source)

* [Setup Server](#setup-server)

  * [Building the Image](#setup-server-build-image)

  * [Custom Image](#setup-server-custom-image)

    * [Image Specific Options](#custom-image-specific-options)

* [Start Server Locally](#start-local)

  * [Start with Script](#start-local-script)

  * [Start with Docker Commands](#start-local-docker-commands)

* [Testing](#testing)
  
## [Getting Started](#getting-started)

The Gateway service is designed to run within a linux container. This README will describe the key files used to build and run this service in a container. Additionally, there are certain environment variables that the application expects to be present in order to run. These environment variables will also be described in this document.

## [Structure](#structure)

The root of the gateway directory contains the supporting files for building the application.

### [Config and Setup Files](#structure-files)

[Dockerfile:](./Dockerfile) multi-stage docker file to build the gateway executable and the gateway image

[.dockerignore:](./.dockerignore) identifies which files should be accessible to commands in the Dockerfile

[.gitignore:](./.gitignore) identifies which files in the gateway directory should not be tracked by git

[gateway-service-api.yaml:](./gateway-service-api.yam) documents the public REST based APIs provided by the gateway service directly. For specific versions, see the [api directory](./../api/).

[localStartExample.ps1:](./localStartExample.ps1) is meant for local testing of the gateway in a docker container

[testGatewayUnit.ps1:](./testGatewayUnit.ps1) is meant for runing the gateway unit tests locally with coverage

### [Gateway Source](#structure-gateway-source)

The [gateway](./gateway/) directory contains the source files for the gateway executable.

[main.go:](./gateway/main.go) the source code containing the main function (entrypoint) for the gateway service

[go.mod:](./gateway/go.mod) file containing the modules used by the application and its dependencies. Used by the go command line tools to identify which packages and their specific versions to retrieve when building the gateway service into an executable

[go.sum:](./gateway/go.sum) used to track and ensure validity of retrieved package files listed in go.mod

## [Setup Server](#setup-server)

The gateway executable is designed to be deployed using a linux container. The following subsections explain how this container is built and how to use it. The gateway executable is currently built on the [GoLang](https://hub.docker.com/_/golang) image `golang:1.12`, and the gateway image is then based off the [Alpine Linux](https://hub.docker.com/_/alpine) image `alpine:3.9`.

### [Building the Image](#setup-server-build-image)

Builds of this container image are automatically triggered by pushes to the GitHub repository.

Builds are tagged based on the version of the API the gateway implements (as defined in an variable in the azure-pipelines.yml file in the root of this repository which should reflect the API version listed in the gateway-service-api.yaml file in this directory). For a complete description of the possible tags see the [gateway container repository](https://hub.docker.com/r/uwthalesians/gateway) on the container registry DockerHub.

#### [Build](#setup-server-build-image)

The image can be built locally using the docker build command. This command should be run from this directory (where the Dockerfile is located). See the local start script for an automated build and run.

Example docker build command: `docker build --tag "${GATEWAY_IMAGE_AND_TAG}" --no-cache .`

Commands:

  `--tag` is the name the image will be given (the name used to then run the built image as a container)

  `--no-cache` ensures that source code is always updated in container

  `.` the final period in the command indicates the root directory to send to the docker deamon for the build process, this should be the directory where the Dockerfile is located

### [Custom Image](#setup-server-custom-image)

The Gateway image will be used during development and production. Information about this custom image can be found in the Thalesians container registry on DockerHub [uwthalesians/mssql](https://hub.docker.com/r/uwthalesians/gateway).

Please refer to the description on the [container registry](https://hub.docker.com/r/uwthalesians/gateway) for specifics on how to configure it. The information below only provides an exmaple setup.

#### [Image Specific Options](#custom-image-specific-options)

This section list any configuration options for the custom image.

##### [Container Environment Variables](#custom-image-env-vars)

Environment variables must be passed to be accessible to the gateway executable (such as using the --env flag with the docker run commnad)

Use the following variables to configure the gateway for the given environment.

`GATEWAY_LISTEN_ADDR=[[<host>]:[<port>]]` (OPTIONAL) identifies what host and port the gateway should listen for requests on. If this variable is not set the gateway will default to ":443".

`GATEWAY_TLSCERTPATH=<pathToCert>` (REQUIRED) identifies the absolute path to the certificate file to be used by the gateway to make TLS connections. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container.

`GATEWAY_TLSKEYPATH=<pathToCertKey>` (REQUIRED) identifies the absolute path to the key file for the certificate identified by the "GATEWAY_TLSCERTPATH" variable. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container.

`GATEWAY_SESSION_KEY=<sessionkey>` (REQUIRED) the session key used to sign login sessions

`MSSQL_SCHEME=<scheme>` (REQUIRED) identifies the scheme to use to connect to the mssql database

`MSSQL_USERNAME=<username>` (REQUIRED) identifies the username to login to the mssql database with

`MSSQL_PASSWORD=<password>` (REQUIRED) the password used to login to the mssql database

`MSSQL_HOST=<host>` (REQUIRED) the host where the mssql database can be reached

`MSSQL_PORT=<port>` (REQUIRED) the port where the mssql database can be reached

`MSSQL_DATABASE=<database>` (REQUIRED) the database to use for the connection

## [Start Server Locally](#start-local)

This setup explains how to build and start the server locally.

### [Start with Script](#start-local-script)

Building and starting the gateway container locally can be more involed than running one script. The following must be setup:

1. Populate the [./gateway/encrypt](./gateway/encrypt) directory with the apprioriate TLS certificates for localhost. See [./gateway/encrypt/README.md](./gateway/encrypt/README.md) for instructions.

2. Ensure service dependencies are set up to run. If using the PowerShell script, the dependencies should be started by default. (Meaing, if you are using the PowerShell script with the default options you should not need to start the dependencies)

   * [mssql](https://hub.docker.com/r/uwthalesians/mssql): image with version 0.7.1 bootstraped (see [local start script](./localStartExample.ps1) for example)

   * [redis](https://hub.docker.com/_/redis): redis:5.0.4-alpine (see [local start script](./localStartExample.ps1) for example)

#### PowerShell

For testing the gateway locally, the [localStartExample.ps1](./localStartExample.ps1) script can be used. This script assumes that docker is already installed and running on the system and that the TLS cert and key have been generated (see note above). Note, the script is a PowerShell script and thus requires a PowerShell shell. Additionally, PowerShell will not run unsigned scripts by default, therefore you may need to enable running unsigned scripts to use it.

The PowerShell script, [localStartExample.ps1](./localStartExample.ps1) will run the gateway image as a container inside a docker network and exposed to localhost. This script has several command line options which allow you to customize the instance.

##### Comand Line Options

The script accepts several comand line options which can be set when running the script in a PowerShell terminal. No positional options.

Unless you need to run your own mssql container or build the gateway container locally, you should not have to provide any options to the local start script.

However, if you want to retian redis and mssql database between runs of the containers, you need to include the 

Run: `.\locaStartExample.ps1 `

`-MsSqlDatabase` (string) which is used to pass in the name of the database to use for connections to the server, default value is: "Perceptia"

`-MsSqlHost` (string) which is used to pass in the hostname of the mssql database server, default value is: "mssql"

`-MsSqlPassword` (string) which is used to pass in either the password to use to secure the mssql server, the default value is: "SecureNow!"

`-MsSqlPort` (string) which is the port the docker container should listen for requests on and send to the mssql server, default value is: "1433",

`-MsSqlScheme` (string) which is the URL type scheme used to connect to the mssql database server, default value is: "sqlserver"

`-MsSqlUsername` (string) which is the username used to connect to the mssql database server, default value is: "sa"

`-PerceptiaDockerNet` (string) which is the name of the docker network the container should be attached to when run, default value is "perceptia-net"

`-GatewayPort` which is the port the gateway service should be exposed on the host machine, default value is "4443"

`-SkipRedis` (switch) will skip starting the redis dependency, default is: false. If you are starting your own redis instance include the option by including the switch

`-SkipMsSql` (switch) will skip starting the mssql dependency, default is: false. If you are starting your own mssql instance, set option by including the switch

`-KeepMsSqlDb` (switch) will start the mssql dependency with an existing database if it already exists, default is: false. If you want to retain a previously created database, set the option by including the switch

`-KeepRedisDb` (switch) will start the redis dependency with an existing database if it already exists, default is: false. If you want to retain a previously created database, set the option by including the switch

`-BuildGateway` (switch) will build the gateway using the local source, default is: false. To set true, include the switch

### [Start with Docker Commands](#start-local-docker-commands)

For directions to start the container locally using a script, see [Start Server Locally](#start-local).

1. pull the image from docker (check [registry](https://hub.docker.com/r/uwthalesians/gateway) for latest images)

   `docker pull uwthalesians/0.1.1-build-latest-branch-develop`

2. run the container image (replace variables with the correct values)

   `docker run --detach --env GATEWAY_SESSION_KEY="$GATEWAY_SESSION_KEY" --env GATEWAY_TLSCERTPATH="$GATEWAY_TLSCERTPATH" --env GATEWAY_TLSKEYPATH="$GATEWAY_TLSKEYPATH" --env MSSQL_DATABASE="$MSSQL_DATABASE" --env MSSQL_HOST="$MSSQL_HOST" --env MSSQL_PASSWORD="$MSSQL_PASSWORD" --env MSSQL_PORT="$MSSQL_PORT" --env MSSQL_SCHEME="$MSSQL_SCHEME" --env MSSQL_USERNAME="$MSSQL_USERNAME" --name ${GATEWAY_CONTAINER_NAME} --network $PerceptiaDockerNet --publish "${GatewayPort}:443" --restart on-failure --mount type=bind,source="$GATEWAY_TLSMOUNTSOURCE",target="/encrypt/",readonly 0.1.1-build-latest-branch-develop`

## [Testing](#testing)

### [Go Test Command](#testing-go-test)

Go build tags are used to identify test types and allow selectively running tests, such as unit tests and integration tests. Testing files have the go build directive options to identify which build tags to run the test for:

`// +build tag_example all unit etc`

With this build directive at the top of the `some_test.go` file, followed by a blank line, when go test is run only the explicit tags provided to go test that match a build directive tag will be run. Example:

`go test -tags=unit ./...`

This command (assuming it is run from the same directory as the root go.mod file) will run all test files that contain the `unit` build tag.

### [Go Unit Test Scripts](#testing-go-unit-test-scripts)

The [./testGatewayUnit.ps1](./testGatewayUnit.ps1) script will run all unit tests for the gateway source code, with coverage, and load a webpage with the results of the coverage test.