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

The local start script [./localStartExample.ps1](./localStartExample.ps1) wraps the docker stack commands to start and stop the Perceptia backend. See the [docker stack config section below](#docker-stack-config) for how to connect to the backend.

To start the backend run the script with no options:

`./localStartExample.ps1`

To stop the backend, run the script with the option `-CleanUp`

`./localStartExample.ps1 -CleanUp`

### Comand Line Options

`-CleanUp` (Switch) when used, stops the services started by docker stack deploy

### Docker Stack Config

The docker stack config makes certain assumptions about the images being run and the local ip environemnt on the device where the script is being run.

The application backend is listening for requests on the localhost using https at the port 4443.

Scheme: `https`

Host: `localhost`

Port: `4443`