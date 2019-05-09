# Redis

This directory contains the code used to manage the redis instances used by the application backend.

## [Contents](#contents)

* [Getting Started](#getting-started)

* [Setup](#setup)

* [Start Server Locally](#start-server-locally)

  * [Start with Script](#start-with-script)

    * [PowerShell](#powershell)

## [Getting Started](#getting-started)

## [Setup](#setup)

## [Start Server Locally](#start-server-locally)

### [Start with Script](#start-with-script)

#### [PowerShell](#powershell)

The PowerShell script, [localStartExample.ps1](./localStartExample.ps1) will run a default redis server container. This script has several command line options which allow you to customize the instance.

##### Comand Line Options

The script accepts several comand line options which can be set when running the script in a PowerShell terminal. No positional options, all options require the explicit flag

Run: `.\locaStartExample.ps1`

`-RedisPortPublish` which is the port the docker container should listen for requests on and send to the redis server, default value is: "6379",

`-PerceptiaDockerNet` which specifies the name to use for the docker network to connect the container to, this should be set to the same docker network that is used by the other backend containers, default value is: "perceptia-net"
