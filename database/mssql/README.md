# Microsoft SQL Server

This directory contains the code used to manage the mssql databse(s) used by the application backend.

## Getting Started

In order to use a database locally, you will need to run a docker container for the database server and attach a storage medium to that container for the database files. This document provides an overview of the structure of this direcotry, the specific database(s) provided, and how to [setup](#setup-server) the database server.

## Structure

This directory is organized around specific databases, each with their own subdirectory.

### Files

**Dockerfile:** docker file to build custom mssql server image

**entrypoint&#x2e;sh:** bash script run inside custom image on start with no arguments

**setup-db&#x2e;sh:** bash script run inside custom image on start to bootstrap Percepia db

**localStartExample&#x2e;ps1:** PowerShell script providing an example of running the ms sql server container. See [manual-setup](#manual-setup) below for more information about this script

### Perceptia ./Perceptia

Perceptia contains the files necessary to build the Perceptia database. This database is used by the application to manage users.

## [Setup](#setup-server)

We will be using the [Microsoft SQL Server](https://hub.docker.com/_/microsoft-mssql-server) docker image for our local MS SQL Server. For informatin on configuring this container see [this microsoft doc](https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker?view=sql-server-2017). This setup will first give an overview of the configuration items (tools, variables, etc.), then provide informaiton about the custom container that will be used in development, and end with a description of an example [manual setup](#manual-setup).

Note, idealy, running this container will be part of a Kubernetes configuration file, so you should not have to run these commands manually. This section is meant to document what the configuration file would otherwise automate. Additionally, in production our applicatio will use an [Azure SQL Server](https://azure.microsoft.com/en-us/services/sql-database/) to host the application database.

### Image Name with Tag

`mcr.microsoft.com/mssql/server:2017-CU12-ubuntu` is the fully qualified image name and container registry address for the MS SQL server we will be running

#### Image Specific Options

`ACCEPT_EULA=Y` automatically accepts the user agreement

`'SA_PASSWORD=******'` is the password for the system administrator, where the `******` would be supplied using an environment variable. This is the password that the application will also use to access the database. Note, the SA password is being used by the application to simplify the development environment. In a real system the application should have its own service credentials to access the database

### Custom Image

During development, a custom docker container will be run which contains both the ms sql server and the scripts necessary to bootstrap the Perceptia database. Information about this custom image can be found in the Thalesians container registry on DockerHub [uwthalesians/mssql](https://hub.docker.com/r/uwthalesians/mssql).

Please refer to the description on the container registry for specifics on how to configure it. The information below only provides an exmaple setup.

#### Example Setup using uwthalesians/mssql image

First: pull the image from docker

   `docker pull uwthalesians/mssql:0.3.2`

Next: run the container image

   `docker run --env 'ACCEPT_EULA=Y' --env "SA_PASSWORD=$Env:SA_PASSWORD" --publish 1401:1433 --mount type=volume,source=mssql_vol,destination=/var/opt/mssql --detach --name=mssql uwthalesians/mssql:0.3.2`

If you run the above command, the image will be run with the container name `mssql`, listening for requests on the local loopback addresses such as `localhost` at the host port `1401`. This image, by default, drops any existing database found in the docker volume `mssql_vol` and creates a new Perceptia database. Note, you will need to already have set the environment variable SA_PASSWORD, or replaced that variable with an explicit password.

### [Manual Setup](#manual-setup)

This setup only starts a docker container running the mssql server. You would need to manually apply the Perceptia database scripts.

#### Example Docker Run Command

##### PowerShell

See `localStartExample.ps1` for example:

`docker run --env 'ACCEPT_EULA=Y' --env "SA_PASSWORD=$Env:SA_PASSWORD" --publish 1401:1433 --mount type=volume,source=mssql_vol,destination=/var/opt/mssql --detach --name=mssql mcr.microsoft.com/mssql/server:2017-CU12-ubuntu`

###### Docker Options

`--env` indicates the following string is an environment variable that should be made available to the main process that starts in the container

`--publish 1401:1433` tells the docker daemon to bind the port 1401 on the host to the port 1433 inside the container. Port 1433 is the default port MS SQL listens for requests on. The MS SQL Server can be reached at localhost:1404, or in SSMS at the Server name: `localhost, 1401`

`--mount type=volume,source=mssql_vol,destination=/var/opt/mssql` tells the daemon to mount the source volume inside the container at the destination location in the container's file system. The destination location is where the MS SQL server will look for and create database files

`--detach` tells the docker deamon to run the container in the background of the process that ran the command

`--name=mssql` specifies the name that will be given to the container when it runs. If the container were attached to a docker network then it could be reached by other containers on that network by its name