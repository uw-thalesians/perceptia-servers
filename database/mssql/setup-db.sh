#!/usr/bin/env bash

# Skip setup if specified
# Defailt: Don't skip

SKIP_SETUP=$SKIP_SETUP

if [ SKIP_SETUP -eq "Y"]
then
        echo "SKIP_SETUP=Y, skipping setup-db.sh script"
        exit 0
fi

# wait for the SQL Server to come up
sleep 30s

# Skip setup if database exists
# Default: don't skip

SKIP_SETUP_IF_EXISTS=$SKIP_SETUP_IF_EXISTS

if [ SKIP_SETUP_IF_EXISTS -eq "Y"]
then
        /opt/mssql-tools/bin/sqlcmd \
        -S localhost -U sa -P ${SA_PASSWORD} -d 'master' -b \
        -Q "IF NOT EXISTS(SELECT [name] FROM master.dbo.sysdatabases WHERE [name] = 'Perceptia') THROW 50100, 'Database does not exist', 1"
        if [ $? -eq 0 ]
        then
                echo "SKIP_SETUP_IF_EXISTS=Y, skipping setup-db.sh script"
                exit 0
        fi
fi

# run the setup script to create the DB and the schema in the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d master -i /script/schema.sql

# run the setup script to create the stored procedures for the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d master -i /script/procedure.sql
