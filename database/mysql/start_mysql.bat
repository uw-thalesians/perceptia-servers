docker run -d --name aq-mysql --volume mysql:/var/lib/mysql --cap-add syslog -e MYSQL_ROOT_PASSWORD=mrpw -e user_pass=8aWZjNadxspXQEHu  -m="1g" --memory-swap="1g" --network aq-net uw-thalesians/aq-mysql