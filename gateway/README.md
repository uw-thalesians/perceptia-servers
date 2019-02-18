#API Gateway Service

The API Gateway service serves as the primary entry point for the Perceptia application. The Gateway provides several key services in addition to routing requests to the responsible microservice. These services include: CORS middleware, Session Authentication, Sign on, Account creation and management. 

##Getting Started

The Gateway service is designed to run within a linux container. This README will describe the key files used to build and run this service in a container. Additionally, there are certain environment variables that the application expects to be present in order to run. These environment variables will also be described in this document. 

##Setup
###Directory Structure
####Root ./
The root of the gateway directory contains the supporting files for building the application.

**Dockerfile:** multi-stage docker file to build the gateway executable and the gateway image

**.dockerignore:** identifies which files should be accessible to commands in the Dockerfile

**.gitignore:** identifies which files in the gateway directory should be tracked by git

**gateway-service-api.yaml** documents the public REST based APIs provided by the gateway service directly

####Gateway ./gateway/
 Directory containing the source code for the gateway service.
 
 **main.go:** the source code containing the main function for the gateway service
 
 **go.mod:** file containing the modules used by the application and its dependencies. Used by the go command line tools to identify which packages and their specific versions to retrieve when building the gateway service into an executable 
 
 **go.sum:** use to track and ensure validity of retrieved package files listed in go.mod
 
 ###Building the container image
 
 
 ###Configuration
 ####Requirements
 
 
 ####Environment Variables
 Use the following variables to configure the gateway for the given environment.
 
 `GATEWAY_LISTEN_ADDR` (OPTIONAL) identifies what [[host]:[port]] the gateway should listen for requests on. If this variable is not set the gateway will default to ":443".
 
 `GATEWAY_TLSCERTPATH` (REQUIRED) identifies the absolute path to the certificate file to be used by the gateway to make TLS connections. 
 
 `GATEWAY_TLSKEYPATH` (REQUIRED) identifies the absolute path to the key file for the certificate identified by the "GATEWAY_TLSCERTPATH" variable.