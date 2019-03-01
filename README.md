# perceptia-servers

[![Build Status](https://dev.azure.com/uw-thalesians/Capstone%202019/_apis/build/status/uw-thalesians.perceptia-servers?branchName=master)](https://dev.azure.com/uw-thalesians/Capstone%202019/_build/latest?definitionId=1&branchName=master)

This repository contains the source files for the services that make up the Perceptia application.

## [Getting Started](#getting-started)

Each service that makes up the Perceptia application is contained in its own subdirectory from the root of the repository (see [Structure](#structure) below).

## [Setup](#setup)

## [Structure](#structure)

The Perceptia application's backend is developed using a microservices architecture. This is reflected in the organization of this repository, with each subdirectory roughly cooresponding to one service of the application. There are additional subdirectories to maintain supporting code. Each subdirectory should have a README.md file which provides additional information about the files in that directory and how to use them.

**./infrastructure/** which contains the supporting code for building and deploying the application

**azure-pipelines.yml** which defines the continuous integration pipeline for the application, including automated testing and artifact building

## [Azure Boards Integration](#azure-boards-integration)

To have commits and PRs for this repository appear as a link in an ADO work-item you have to use a specific syntax in your commit and PR messages. Read more about this proccess [here.](https://docs.microsoft.com/en-us/azure/devops/boards/github/link-to-from-github?view=vsts)

Note, in order for ADO to know to watch for the Azure Board tag, the repository must already be selected as a connection in the [ADO project settings.](https://dev.azure.com/uw-thalesians/Capstone%202019/_settings/boards-external-integration) Instructions for how to set this up can be found [here.](https://docs.microsoft.com/en-us/azure/devops/boards/github/index?view=vsts)

### Commit Format

AB#{ID}

If you include the above, where {ID} is replaced with the work-item id from ADO, in your commit or PR, the coresponding ADO work-item will attach a link to the commit or PR to the work-item. Note, there are additional keywords that ADO will watch for in a commit message with the AB#{ID} format, and take specific actions. See [here](https://docs.microsoft.com/en-us/azure/devops/boards/github/link-to-from-github?view=vsts) for more information.  

## [Public Repository Security Considerations](#security-considerations)

This is a public repository. Do not store any sensitive information in this repository, such as secure API access tokens, certificates, private keys, etc. If your build process depends on this content, be sure to add the file to the .gitignore before saving it to the local clone of the repository, or load the information by an envirnment variable. Sensitive informaiton should be stored in the Azure Pipelines (AZP) library, or in cloud provider vaults.
