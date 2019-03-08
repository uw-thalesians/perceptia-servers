# API Gateway Service

Last updated: 2019-02-21

The API Gateway service serves as the primary entry point for the Perceptia application. The Gateway provides several key services in addition to routing requests to the responsible microservice. These services include: CORS middleware, Session Authentication, Sign on, Account creation and management. 

## Getting Started

The Gateway service is designed to run within a linux container. This README will describe the key files used to build and run this service in a container. Additionally, there are certain environment variables that the application expects to be present in order to run. These environment variables will also be described in this document. 

## Setup
### Directory Structure
#### Root ./

The root of the gateway directory contains the supporting files for building the application.

**Dockerfile:** multi-stage docker file to build the gateway executable and the gateway image

**.dockerignore:** identifies which files should be accessible to commands in the Dockerfile

**.gitignore:** identifies which files in the gateway directory should be tracked by git

**gateway-service-api.yaml:** documents the public REST based APIs provided by the gateway service directly

**localStartExample.ps1:** is meant for local testing of the gateway in a docker container

#### Gateway ./gateway/

 Directory containing the source code for the gateway service.
 
 **main.go:** the source code containing the main function for the gateway service
 
 **go.mod:** file containing the modules used by the application and its dependencies. Used by the go command line tools to identify which packages and their specific versions to retrieve when building the gateway service into an executable 
 
 **go.sum:** used to track and ensure validity of retrieved package files listed in go.mod
 
 ### Building the container image
 
Builds of this container image are automatically triggered by pushes to the GitHub repository. Currently, the latest commit-hash is appended to the image tag and then the image is pushed to the container registry for uwthalesians on DockerHub. The version part of the tag is based on the [semver](https://semver.org/) format.
 
 ### Running the Gateway Locally
 
 For testing the gateway locally, the localStartExample.ps1 script can be used. This script assumes that docker is already installed and running on the system and that the TLS cert and key have been generated in the ./gateway/encrypt/ subdirectory. Note, the script is a PowerShell script and thus requires a PowerShell shell. Additionally, PowerShell will not run unsigned scripts by default, therefore you may need to enable running unsigned scripts to use it. 
 
 The example script also builds the docker container on each run. In the future it will instead pull from our container registry the latest image. This has not been done yet as our specific tagging and use of the container registry has not been defined yet. 
 
 ### Configuration
 
 #### Requirements
 
 TODO
 
 #### Environment Variables
 
 Use the following variables to configure the gateway for the given environment.
 
 `GATEWAY_LISTEN_ADDR` (OPTIONAL) identifies what [[host]:[port]] the gateway should listen for requests on. If this variable is not set the gateway will default to ":443".
 
 `GATEWAY_TLSCERTPATH` (REQUIRED) identifies the absolute path to the certificate file to be used by the gateway to make TLS connections. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container.
 
 `GATEWAY_TLSKEYPATH` (REQUIRED) identifies the absolute path to the key file for the certificate identified by the "GATEWAY_TLSCERTPATH" variable. This path is based on where the gateway executable is being run, so if it is being run in a container, the path referenced must be accessible within the container.
 
 #### Integration
 
 TODO
 
 #### Testing
 
 Go build tags are used to identify test types and allow selectively running tests, such as unit tests and integration tests. Testing files have the go build directive options to identify which build tags to run the test for:
 
 `// +build tag_example all unit etc`
 
 With this build directive at the top of the `some_test.go` file, followed by a blank line, when go test is run only the explicit tags provided to go test that match a build directive tag will be run. Example: 
 
 `go test -tags=unit ./...` 
 
 This command (assuming it is run from the same directory as the root go.mod file) will run all test files that contain the `unit` build tag.