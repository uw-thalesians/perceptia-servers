#!/usr/bin/env bash

# Skip setup if specified
# Defailt: Don't skip


if [ "$MSSQL_ENVIRONMENT" == "" ]
then
        echo "MSSQL_ENVIRONMENT not set, defaulting to development"
        export MSSQL_ENVIRONMENT="development"
fi
echo "MSSQL_ENVIRONMENT=$MSSQL_ENVIRONMENT"

# wait for the SQL Server to come up
if [ "$MSSQL_ENVIRONMENT" == "production" ]
then
        echo "Sleeping for 30s to allow server to start"
        sleep 30s
elif [ "$MSSQL_ENVIRONMENT" == "development" ]
then
        echo "Sleeping for 20s to allow server to start"
        sleep 20s
fi

echo "SKIP_SETUP=$SKIP_SETUP"
echo "SKIP_SETUP_IF_EXISTS=$SKIP_SETUP_IF_EXISTS"

# Setup Database to use containered databases
echo "Setting server to allow contained databases"
/opt/mssql-tools/bin/sqlcmd \
-S localhost -U sa -P ${SA_PASSWORD} -b \
-Q "EXECUTE sp_configure 'contained database authentication', 1; RECONFIGURE;"
if [ $? -eq 1 ]
then
        echo "Unable to set Database Engine to use contained database authentication."
        exit 1
fi

echo "Checking if SKIP_SETUP variable was set to Y"
if [ "$SKIP_SETUP" == "Y" ]
then
        echo "SKIP_SETUP=Y, skipping setup-db.sh script"
        exit 0
fi



# Skip setup if database exists
# Default: don't skip

SKIP_SETUP_IF_EXISTS=$SKIP_SETUP_IF_EXISTS

echo "Checking if SKIP_SETUP_IF_EXISTS was set to Y"
if [ "$SKIP_SETUP_IF_EXISTS" == "Y" ]
then
        echo "Checking to see if Perceptia database exists"
        /opt/mssql-tools/bin/sqlcmd \
        -S localhost -U sa -P ${SA_PASSWORD} -d 'master' -b \
        -Q "IF NOT EXISTS(SELECT [name] FROM master.dbo.sysdatabases WHERE [name] = 'Perceptia') THROW 50100, 'Database does not exist', 1"
        if [ $? -eq 0 ]
        then
                echo "SKIP_SETUP_IF_EXISTS=Y, skipping setup-db.sh script"
                exit 0
        fi
fi

# create Perceptia database
echo "Creating Perceptia database on server"
/opt/mssql-tools/bin/sqlcmd \
-S localhost -U sa -P ${SA_PASSWORD} -d 'master' -b \
-Q "If Exists(SELECT [name] FROM master.dbo.sysdatabases WHERE [name] = 'Perceptia')
Begin
	USE [master]
	ALTER DATABASE [Perceptia] SET SINGLE_USER WITH ROLLBACK IMMEDIATE
	DROP DATABASE [Perceptia]
End
;
CREATE DATABASE [Perceptia]
	CONTAINMENT = PARTIAL
	COLLATE Latin1_General_100_CI_AS_SC
;"
if [ $? -eq 1 ]
then
        echo "Unable to create $GATEWAY_SP_USERNAME user."
        exit 1
fi


# run the setup script to create the DB and the schema in the DB
echo "Applying schema.sql to Perceptia database"
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d Perceptia -i /script/Perceptia/schema.sql

# run the setup script to create the stored procedures for the DB
echo "Applying procedure.sql to Perceptia database"
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d Perceptia -i /script/Perceptia/procedure.sql

# run the setup script to populate the DB
echo "Applying procedure.sql to Perceptia database"
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P ${SA_PASSWORD} -d Perceptia -i /script/Perceptia/populate.sql

# add SQL database user
echo "Adding $GATEWAY_SP_USERNAME user to Perceptia database"
/opt/mssql-tools/bin/sqlcmd \
-S localhost -U sa -P ${SA_PASSWORD} -d 'Perceptia' -b \
-Q "CREATE USER $GATEWAY_SP_USERNAME WITH PASSWORD = '$GATEWAY_SP_PASSWORD';"
if [ $? -eq 1 ]
then
        echo "Unable to create $GATEWAY_SP_USERNAME user."
        exit 1
fi

echo "Adding $GATEWAY_SP_USERNAME to the role: RL_ExecuteAllProcedures"
/opt/mssql-tools/bin/sqlcmd \
-S localhost -U sa -P ${SA_PASSWORD} -d 'Perceptia' -b \
-Q "ALTER ROLE RL_ExecuteAllProcedures ADD MEMBER $GATEWAY_SP_USERNAME;"
if [ $? -eq 1 ]
then
        echo "Unable to add $GATEWAY_SP_USERNAME user to Role RL_ExecuteAllProcedures."
        exit 1
fi

# Wait for database to finish loading
echo "Sleeping for 20 seconds to allow the database to finish loading..."
sleep 20s

# Remove Database scripts
if [ "$MSSQL_ENVIRONMENT" == "production" ]
then
        rm /script/Perceptia/schema.sql /script/Perceptia/procedure.sql /script/Perceptia/populate.sql
fi

# Ensure Database was created successfully
echo "Ensuring database was created successfully"
/opt/mssql-tools/bin/sqlcmd \
-S localhost -U sa -P ${SA_PASSWORD} -d 'master' -b \
-Q "IF NOT EXISTS(SELECT [name] FROM master.dbo.sysdatabases WHERE [name] = 'Perceptia') THROW 50100, 'Database does not exist', 1"
if [ $? -eq 1 ]
then
        echo "Database was not created successfully!"
        exit 1
else 
        echo "Database found, Perceptia database created successfully!"
fi