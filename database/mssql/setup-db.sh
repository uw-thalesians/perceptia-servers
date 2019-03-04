#!/usr/bin/env bash

# wait for the SQL Server to come up
sleep 30s

# run the setup script to create the DB and the schema in the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d master -i /script/schema.sql

# run the setup script to create the stored procedures for the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d master -i /script/procedure.sql
