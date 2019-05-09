# AnyQuiz REST API Container

This image provides a customized version of the official `php:7.2-apache` image which installs python2.7, pip and enables some apache functionality like rewrite, and the cgi module, as well as the PHP PDO extension.

## Building the Docker Image

Building the docker image may require more memory and/or swap space be allowed to the Docker Service to allow containers to use more.

![Docker setting screenshot](/AdvDockerSettings.PNG)

## Example to run locally

`start_net.bat` - example to create a bridge network that all the docker containers connect to


`make_aq.bat` - example to create docker image


`start_aq.bat` - example to start up the AnyQuiz container with an example local bind to any_quiz repo code on windows, drive share needs to be enabled in docker

## Minimum requirements
While it may run with somewhat less, building the container appears to require around **2.5G RAM and Swap each** allocated to Docker Daemon. When running it is presumed that this amount will still hold true for loading NLP packages into memory in python et c.

These flags can be provided when first running the container if Docker Daemon is allocated more than 2.5G itself to tell Docker this is an explicit amount to allocate for this container.

`--memory="2.5G" --memory-swap="2.5G"`

## Parameters to pass to container on `docker run`

The environment variable `user_pass` must be passed to the container when it is first run. This **must** match the user_pass passed when first running the MySQL container. It is the application user password for MySQL.

`-e user_pass=8aWZjNadxspXQEHu`


## Putting it all together

`docker run -d --name aqrest -p 8082:80 -e user_pass=8aWZjNadxspXQEHu --network perceptia-net --memory="2.5G" --memory-swap="2.5G" uw-thalesians/aqrest`
