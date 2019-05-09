REM using path to this batch file as root for aq-rest
@echo off
REM 
REM add using path for output of dev config.php (generated from script), instead of current working directory script
REM is called from...
docker run -d --cap-add syslog --name aqrest -p 8082:80 -e dev=1 -e user_pass=8aWZjNadxspXQEHu --network aqnet --memory="2.5G" --memory-swap="2.5G" --mount type=bind,source=%~dp0aq,target=/var/www/html uw-thalesians/aqrest