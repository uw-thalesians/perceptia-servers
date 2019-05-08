#! /usr/bin/env bash

echo "generating scripts with user pass from template"
sed s/\{pass\}/$user_pass/ /var/www/html/py/config.py.template > /var/www/html/py/config.py
sed s/\{pass\}/$user_pass/ /var/www/html/sql/config.php.template > /var/www/html/sql/config.php

if [ -z ${dev+x} ]; then
    echo "removing templates"
    rm /var/www/html/py/config.py.template
    rm /var/www/html/sql/config.php.template
fi

echo "removing this script"
rm /usr/local/bin/parse-template.sh

echo "returning entrypoint to previous startup"
ln -s /usr/local/bin/docker-php-entrypoint /usr/local/bin/parse-template.sh

/usr/local/bin/docker-php-entrypoint "$@"