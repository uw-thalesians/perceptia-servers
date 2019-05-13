FROM php:7.2-apache
SHELL ["/bin/bash", "-c"]

# enable some used apache modules
RUN ln -s ../mods-available/{expires,headers,rewrite}.load /etc/apache2/mods-enabled/

# allow directory override/rewrite
RUN sed -e '/<Directory \/var\/www\/>/,/<\/Directory>/s/AllowOverride None/AllowOverride All/' -i /etc/apache2/apache2.conf

#copy default production php.ini as the php.ini to use
RUN cp /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini

# install python 2.7 and pip, includes a few utilities for debugging and reading configs
RUN apt update && apt install -y python2.7 python-pip

# install mysql php database object extension to php, handles enabling the extension in php.ini
RUN docker-php-ext-install pdo_mysql

# enable cgi to allow python scripts to be run as cgi scripts
RUN a2enmod cgi

RUN pip install --upgrade pip
RUN apt remove -y python-pip

# install python modules used by AnyQuiz

# install a few packages at a time to lower env mem/swap requirements (avoid error 137)
RUN pip install wikipedia mysql-connector
RUN pip install numpy
RUN pip install wasabi cymem preshed pathlib srsly murmurhash
RUN pip install blis plac tqdm thinc functools32 jsonschema
RUN pip install spacy 
RUN pip install nltk

# download used nltk data/packages
RUN python -m nltk.downloader averaged_perceptron_tagger -d /usr/share/nltk_data
RUN python -m nltk.downloader punkt -d /usr/share/nltk_data
RUN python -m nltk.downloader wordnet -d /usr/share/nltk_data

# RUN python -m spacy download en_vectors_web_lg

RUN python -m spacy download en_core_web_md

# download spacy en model alternative method
# With external URL
#pip install https://github.com/explosion/spacy-models/releases/download/en_core_web_sm-2.1.0/en_core_web_sm-2.1.0.tar.gz

#add first run script to fill out template files with mysql app user password provided from environment variable
COPY parse-template.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/parse-template.sh

COPY aq/ /var/www/html/

RUN chmod +rx /var/www/html/py/*.py

EXPOSE 80/tcp
EXPOSE 3306/tcp

ENTRYPOINT ["parse-template.sh"]
CMD ["apache2-foreground"]