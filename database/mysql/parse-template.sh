#! /usr/bin/env bash
echo "generating script with user pass from template"
sed s/\{pass\}/$user_pass/ /docker-entrypoint-initdb.d/create_user.sql.template > /docker-entrypoint-initdb.d/create_user.sql
echo "removing this script and template"
rm /docker-entrypoint-initdb.d/create_user.sql.template /usr/local/bin/parse-template.sh

#echo '/usr/local/bin/docker-entrypoint.sh "$@"' > /usr/local/bin/parse-template.sh
echo "returning entrypoint to previous startup"
ln -s /usr/local/bin/docker-entrypoint.sh /usr/local/bin/parse-template.sh

/usr/local/bin/docker-entrypoint.sh "$@"