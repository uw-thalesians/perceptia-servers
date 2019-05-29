#! /usr/bin/env bash

/var/www/html/py/persistent_server.py &
echo "py server start status $?"

/usr/local/bin/docker-php-entrypoint "$@"