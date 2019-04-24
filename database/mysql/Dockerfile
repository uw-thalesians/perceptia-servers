FROM mysql:5.5
SHELL ["/bin/bash", "-c"]

# If there are any .sh or .sql scripts you want to run on the database immediately after it has been created, you can put them into a host directory and then mount the directory at /docker-entrypoint-initdb.d/ inside the container.
# https://dev.mysql.com/doc/refman/5.6/en/docker-mysql-more-topics.html#docker-additional-init

COPY create_tables.sql create_user.sql.template parse-template.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/parse-template.sh
ENTRYPOINT [ "/docker-entrypoint-initdb.d/parse-template.sh" ]
