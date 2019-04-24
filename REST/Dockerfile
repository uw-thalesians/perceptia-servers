FROM php:7.2-apache
SHELL ["/bin/bash", "-c"]

# enable some used apache modules
RUN ln -s ../mods-available/{expires,headers,rewrite}.load /etc/apache2/mods-enabled/

# allow directory override/rewrite
RUN sed -e '/<Directory \/var\/www\/>/,/<\/Directory>/s/AllowOverride None/AllowOverride All/' -i /etc/apache2/apache2.conf

#copy default production php.ini as the php.ini to use
RUN cp /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini

# install python 2.7 and pip, includes a few utilities for debugging and reading configs
RUN apt update && apt install -y default-mysql-client less python2.7 python-pip

# install mysql php database object extension to php, handles enabling the extension in php.ini
RUN docker-php-ext-install pdo_mysql

# enable cgi to allow python scripts to be run as cgi scripts
RUN a2enmod cgi

EXPOSE 80/tcp
EXPOSE 3306/tcp