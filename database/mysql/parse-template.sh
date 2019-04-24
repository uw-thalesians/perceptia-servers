#! /usr/bin/env bash
echo "generating script with user pass from template"
sed s/\{pass\}/$user_pass/ /docker-entrypoint-initdb.d/create_user.sql.template > /docker-entrypoint-initdb.d/create_user.sql
echo "removing this script and template"
rm /docker-entrypoint-initdb.d/create_user.sql.template /docker-entrypoint-initdb.d/parse-template.sh
echo "returning to previous startup"

/usr/local/bin/docker-entrypoint.sh mysqld