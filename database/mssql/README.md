# Microsoft SQL Server

This directory contains the code used to manage the mssql databse(s) used by the application backend.

## [Contents](#contents)

* [Getting Started](#getting-started)

* [Structure](#structure)

  * [Config and Setup Files](#config-and-setup-files)

  * [Databases](#databases)

* [Setup Server](#setup-server)

  * [Base Image](#base-image)

  * [Custom Image](#custom-image)

* [Start Server Locally](#start-server-locally)

  * [Start with Script](#start-with-script)

  * [Start with Docker Commands](#start-with-docker-commands)

## [Getting Started](#getting-started)

In order to use a database locally, you will need to run a docker container for the database server and attach a storage medium to that container for the database files (if you want the data to persist). This document provides an overview of the structure of this direcotry, the specific database(s) provided, and how to [setup](#setup-server) the database server.

## [Structure](#structure)

This directory is organized around specific database(s), each with their own subdirectory. The files in this directory support building those database(s).

### [Config and Setup Files](#config-and-setup-files)

[Dockerfile:](./Dockerfile) docker file to build custom mssql server image

[entrypoint.sh:](./entrypoint.sh) bash script run inside custom image on start with no arguments

[setup-db.sh:](./setup-db.sh) bash script run inside custom image on start to bootstrap Percepia db

[localStartExample.ps1:](./localStartExample.ps1) PowerShell script providing an example of running the ms sql server container. See [Start with Script](#start-with-script) below for more information about this script

[.dockerignore:](./.dockerignore)

### [Databases](#databases)

Each database is contained in a subdirectory of this directory. Each directory contains the necessary files, such as sql files, to bootstrap the given database. This includes the database schema and any stored procedures.

#### [Perceptia](#perceptia)

Perceptia contains the files necessary to build the Perceptia database. This database is used by the gateway service to manage users.

## [Setup Server](#setup-server)

We will be using the [Microsoft SQL Server](https://hub.docker.com/_/microsoft-mssql-server) docker image for our local MS SQL Server. For informatin on configuring this container see [this microsoft doc](https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker?view=sql-server-2017). This setup will first give an overview of the configuration items (tools, variables, etc.), then provide informaiton about the custom container that will be used in development, and end with a description of an example [Start with Script](#start-with-script).

Note, idealy, running this container will be part of a Kubernetes configuration file, so you should not have to run these commands manually. This section is meant to document what the configuration file would otherwise automate. Additionally, in production our applicatio will use an [Azure SQL Server](https://azure.microsoft.com/en-us/services/sql-database/) to host the application database.

### [Base Image](#base-image)

`mcr.microsoft.com/mssql/server:2017-CU12-ubuntu` is the fully qualified image name and container registry address for the MS SQL server we will be running

#### [Base Image Specific Options](#base-image-specific-options)

This section lists specific configuration options for the base Ms Sql image.

##### [Base Image Container Environment Variables](#base-image-container-environment-variables)

`ACCEPT_EULA=Y` automatically accepts the user agreement

`'SA_PASSWORD=******'` is the password for the system administrator, where the `******` would be supplied using an environment variable. This is the password that the application will also use to access the database. Note, the SA password is being used by the application to simplify the development environment. In a real system the application should have its own service credentials to access the database

##### [Custom Mount Points](#custom-mount-points)

To save and/or load a database for use beyond the lifecycle of one container, you should mount a volume at the location where the server saves the database files. The Ms Sql server maintains the database files at this location: `/var/opt/mssql`

Example mount option to docker: `--mount type=volume,source={MSSQL_VOLUME_NAME},destination=/var/opt/mssql`

If a volume is mounted at this location, as long as the volume is not deleted, any databases created can be loaded by a fresh container using this same mount point.

### [Custom Image](#custom-image)

During development, a custom docker container will be run which contains both the ms sql server and the scripts necessary to bootstrap the Perceptia database. Information about this custom image can be found in the Thalesians container registry on DockerHub [uwthalesians/mssql](https://hub.docker.com/r/uwthalesians/mssql).

Please refer to the description on the [container registry](https://hub.docker.com/r/uwthalesians/mssql) for specifics on how to configure it. The information below only provides an exmaple setup.

#### [Custom Image Specific Options](#custom-image-specific-options)

This section list any configuration options for the custom image in addition to any [options from the base image](#base-image-specific-options).

##### [Custom Container Environment Variables](#custom-container-environment-variables)

All [environment variables for the base image](#base-image-container-environment-variables) apply in addition to the container environment variables listed below:

 `SKIP_SETUP=Y` (optional) if value is "Y" will skip running setup-db.sh script which bootstraps the database schema, any other value besides "Y" will be ignored, as if "SKIP_SETUP" was not set

`SKIP_SETUP_IF_EXISTS=Y` (optional) if value is "Y" will skip running setup-db.sh script which bootstraps the database schema if the Perceptia database already exists. Any other value besides `Y` will be ignored, as if `SKIP_SETUP_IF_EXISTS` was not set

`GATEWAY_SP_USERNAME={username}` (required) sets the name to use to create the service principal user that the gateway service will use to connect to the gateway

`GATEWAY_SP_PASSWORD={password}` (required) sets the password to use to authenticate using the service principal user for the gateway service

`MSSQL_ENVIRONMENT={production|development}` (optional) sets the environment the server is running in, production: limited logging, development: verbose logging, default development

#### [Example Setup using uwthalesians mssql image](#example-setup-using-uwthalesians-mssql-image)

1. pull the image from docker (check [registry](https://hub.docker.com/r/uwthalesians/mssql) for latest images)

   `docker pull uwthalesians/mssql:0.7.1-build-latest-branch-develop`

2. run the container image

   `docker run --env 'ACCEPT_EULA=Y' --env "SA_PASSWORD=$SA_PASSWORD" --publish 1401:1433 --mount type=volume,source=mssql_vol,destination=/var/opt/mssql --detach --name=mssql uwthalesians/mssql:0.7.1-build-latest-branch-develop`

If you run the above command, the image will be run with the container name `mssql`, listening for requests on the local loopback addresses such as `localhost` at the host port `1401`. This image, by default, drops any existing database found in the docker volume `mssql_vol` and creates a new Perceptia database. Note, you will need to already have set the environment variable SA_PASSWORD, or replaced that variable with an explicit password.

A note about versions: currently, for versions starting with 0.7.1 of the custom image, the version refers the the specific version of the [stored procedures](./Perceptia/procedure.sql) the database supports. Versions below 1 may make breaking changes. For versions below 0.7.1, the version refered the the specific Perceptia schema the database was built for.

## [Start Server Locally](#start-server-locally)

This setup explains how to build and start the server locally.

### [Start with Script](#start-with-script)

#### PowerShell

The PowerShell script, [localStartExample.ps1](./localStartExample.ps1) will build the custom mssql image, loading in the Perceptia database files and bootstrap scripts. This script has several command line options which allow you to customize the instance.

##### Comand Line Options

The script accepts several comand line options which can be set when running the script in a PowerShell terminal. No positional options, all options require the explicit flag

Run: `.\locaStartExample.ps1 -MsSqlSkipSetupIfExist Y`

`-MsSqlSaPassword` (string) which is used to pass in either the password to use to secure the mssql server, the default value is: "SecureNow!"

`-MsSqlPortPublish` (string) which is the port the docker container should listen for requests on and send to the mssql server, default value is: "1401"

`-MsSqlGatewaySpUsername` (string) which is the username the gateway will use to connect to the Perceptia database as

`-MsSqlGatewaySpPassword` (string) which is the password for the user the gateway will connect to the Perceptia database as

`-MsSqlSkipSetupIfExist` (string) which allows setting what value is passed for the custom image environment variable SKIP_SETUP_IF_EXISTS (see [custom image env vars](#custom-image-env-vars)), default value is: "Y", meaning setup will run (unless another option over rules this)

`-MsSqlSkipSetup` (string) which allows setting what value is passed for the custom image environment variable SKIP_SETUP (see [custom image env vars](#custom-image-env-vars)), default value is: "N", meaning setup will run (unless another option over rules this)

`-PerceptiaDockerNet` which specifies the name to use for the docker network to connect the container to, this should be set to the same docker network that is used by the other backend containers, default value is: "perceptia-net"

`-MsSqlRemoveDbVolume` (switch) when set, removes any existing volume created for the mssql container

`-RemoveAllDbVolumes` (switch) when set, removes any existing volume created by this script

`-CleanUp` (switch) when set, removes any container(s) created by this script

##### Docker Options Explained

This subsection explains the meaning of the various docker options supplied to the docker run command in the local start script.

`--env` indicates the following string is an environment variable that should be made available to the main process that starts in the container

`--publish 1401:1433` tells the docker daemon to bind the port 1401 on the host to the port 1433 inside the container. Port 1433 is the default port MS SQL listens for requests on. The MS SQL Server can be reached at localhost:1404, or in SSMS at the Server name: `localhost, 1401`

`--mount type=volume,source=mssql_vol,destination=/var/opt/mssql` tells the daemon to mount the source volume inside the container at the destination location in the container's file system. The destination location is where the MS SQL server will look for and create database files

`--detach` tells the docker deamon to run the container in the background of the process that ran the command

`--name=mssql` specifies the name that will be given to the container when it runs. If the container were attached to a docker network then it could be reached by other containers on that network by its name

`--network "perceptia-net"` specifies the docker network the container should be attached to. By default, all containers in a docker network can communicate with all other containers

### [Start with Docker Commands](#start-with-docker-commands)

TODO