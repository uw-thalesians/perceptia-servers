REM example user pass is 8aWZjNadxspXQEHu
docker run -d --name aqrest -p 8082:80 -e user_pass=8aWZjNadxspXQEHu --network perceptia-net --memory="2.5G" --memory-swap="2.5G" uw-thalesians/aqrest