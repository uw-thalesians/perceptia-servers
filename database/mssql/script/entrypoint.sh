#!/usr/bin/env bash

# Important: '/opt/mssql/bin/sqlservr' must be the last command to run.
# The last command run controls the container lifecycle, therefore
# sqlservr must be run last.
/script/setup-db.sh & /opt/mssql/bin/sqlservr