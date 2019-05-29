FROM mysql:5.5
SHELL ["/bin/bash", "-c"]

# If there are any .sh or .sql scripts you want to run on the database immediately after it has been created, you can put them into a host directory and then mount the directory at /docker-entrypoint-initdb.d/ inside the container.
# https://dev.mysql.com/doc/refman/5.6/en/docker-mysql-more-topics.html#docker-additional-init

COPY create_tables.sql create_user.sql.template /docker-entrypoint-initdb.d/
COPY parse-template.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/parse-template.sh
RUN printf "[mysqld]\nlog_error=/var/lib/mysql/mysql_error.log\ngeneral_log_file=/var/lib/mysql/mysql_general.log\ngeneral_log=1\nmax_allowed_packet=100M\n" >> /etc/mysql/my.cnf

ENTRYPOINT [ "/usr/local/bin/parse-template.sh" ]
CMD ["mysqld"]