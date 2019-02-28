# Microsoft SQL Server

This directory contains the code used to manage the mssql databses used by the application back-end.

## Getting Started

In order to use a database locally, you will need to run a docker container for the database server and attach a storage medium to that container for the database files. This document provides an overview of the structure of this direcotry, the specific database(s) provided, and how to [setup](#setup-server) the database server.

## Structure

This directory is organized around specific databases, each with their own subdirectory.

### Perceptia

Perceptia contains the files necessary to build the Perceptia database.

## [Setup](#setup-server)

We will be using the [Microsoft SQL Server](https://hub.docker.com/_/microsoft-mssql-server) docker image for our local MS SQL Server. For informatin on configuring this container see [this microsoft doc](https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker?view=sql-server-2017). This setup will first given an overview of the configuration items (tools, variables, etc.) and then provide an example docker run command. 

Note, idealy, running this container will be part of a Kubernetes configuration file, so you should not have to run these commands manually. This section is meant to document what the configuration file would otherwise automate. 

### Image Name with Tag

`mcr.microsoft.com/mssql/server:2017-CU12-ubuntu` is the fully qualified image name and container registry address for the MS SQL server we will be running

#### Image Specific Options

`ACCEPT_EULA=Y` automatically accepts the user agreement

`'SA_PASSWORD=******'` is the password for the system administrator, where the `******` would be supplied using an environment variable. This is the password that the application will also use to access the database. Note, the SA password is being used by the application to simplify the development environment. In a real system the application should have its own service credentials to access the database

### Example Docker Run Command

**Bash**

`docker run --env 'ACCEPT_EULA=Y' --env 'SA_PASSWORD=$(SA_PASSWORD)' --publish 1401:1433 --mount type=volume,source=mssql_vol,destination=/var/opt/mssql --detach --name=mssql mcr.microsoft.com/mssql/server:2017-CU12-ubuntu`

#### Docker Options

`--env` indicates the following string is an environment variable that should be made available to the main process of the container

`--publish 1401:1433` tells the docker daemon to bind the port 1401 on the host to the port 1433 inside the container. Port 1433 is the default port MS SQL listens for requests on

`--mount type=volume,source=mssql_vol,destination=/var/opt/mssql` tells the daemon to mount the source volume inside the container at the destination location in the container's file system. The destination location is where the MS SQL server will look for and create database files

`--detach` tells the docker deamon to run the container in the background of the process that ran the command

`--name=mssql` specifies the name that will be given to the container when it runs. If the container were attached to a docker network then it could be reached by other containers on that network by its name