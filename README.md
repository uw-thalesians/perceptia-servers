# perceptia-servers
This repository contains the source files for the services that make up the Perceptia application.

## [Getting Started](#getting-started)
Each service that makes up the Perceptia application is contained in its own subdirectory from the root of the repository (see [Structure](#structure) below).

## [Setup](#setup)

## [Structure](#structure)


## [Azure Boards Integration](#azure-boards-integration)
To have commits and PRs for this repository appear as a link in an ADO work-item you have to use a specific syntax in your commit and PR messages. Read more about this proccess [here.](https://docs.microsoft.com/en-us/azure/devops/boards/github/link-to-from-github?view=vsts)

Note, in order for ADO to know to watch for the Azure Board tag, the repository must already be selected as a connection in the [ADO prject settings.](https://dev.azure.com/uw-thalesians/Capstone%202019/_settings/boards-external-integration) Instructions for how to set this up can be found [here.](https://docs.microsoft.com/en-us/azure/devops/boards/github/index?view=vsts) 

### Commit Format
AB#{ID}

If you include the above, where {ID} is replaced with the work-item id from ADO, in your commit or PR, the coresponding ADO work-item will attach a link to the commit or PR to the work-item. Note, there are additional keywords that ADO will watch for in a commit message with the AB#{ID} format, and take specific actions. See [here](https://docs.microsoft.com/en-us/azure/devops/boards/github/link-to-from-github?view=vsts) for more information.  


## [Public Repository Security Considerations](#security-considerations)
This is a public repository. Do no store any sensitive information in this repository, such as secure API access tokens, certificates, private keys, etc. If your build process depends on this content, be sure to add the file to the .gitignore before saving it to the local clone of the repository, or load the information by an envirnment variable. At this time there is not an agreed upon location to store this information. (Suggestion: Team Drive Folder).
