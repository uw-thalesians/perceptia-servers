# Database

The database directory contains the files necessary to build, deploy, and manage the databases used by the application backend.

## [Contents](#contents)

* [Overview](#overview)

* [Servers](#servers)

## [Overview](#overview)

Currently, each database used by a service has its setup and configuration files stored in the subdirectories of this directory. Each subdirectory is for a specific database server, which then includes any necessary files to bootstrap the required databases for that server. Each server directory should have a README to explain how the image is built and used.

## [Servers](#servers)

[Microsoft SQL Server](./mssql/) - Used by the gateway service

[REDIS Server](./redis/) - Used by the gateway service