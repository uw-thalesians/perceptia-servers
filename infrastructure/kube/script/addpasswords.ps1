Set-Variable -Name AQMYSQL_ANYQUIZ_USER_PASSWORD -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keys\AqmysqlAnyQuizUserPassword.txt)
Set-Variable -Name AQMYSQL_ROOT_PASSWORD -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keys\AqmysqlRootPassword.txt)
Set-Variable -Name GATEWAY_SESSION_KEY -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keys\GatewaySessionKey.txt)
Set-Variable -Name MSSQL_SA_PASSWORD -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keys\MssqlSaPassword.txt)

Set-Variable -Name AQMYSQL_ANYQUIZ_USER_PASSWORD_DEV -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keysdev\AqmysqlAnyQuizUserPassword.txt)
Set-Variable -Name AQMYSQL_ROOT_PASSWORD_DEV -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keysdev\AqmysqlRootPassword.txt)
Set-Variable -Name GATEWAY_SESSION_KEY_DEV -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keysdev\GatewaySessionKey.txt)
Set-Variable -Name MSSQL_SA_PASSWORD_DEV -Value (Get-Content -Path $Env:SECRET_PERCEPTIA_SERVERS\keysdev\MssqlSaPassword.txt)

# Add development secrets
kubectl create secret generic aqmysql --type=string --from-literal=user_password=$AQMYSQL_ANYQUIZ_USER_PASSWORD_DEV `
--from-literal=root_password=$AQMYSQL_ROOT_PASSWORD_DEV --namespace development 

kubectl create secret generic gateway --type=string --from-literal=session-key=$GATEWAY_SESSION_KEY_DEV --namespace development 

kubectl create secret generic mssql --type=string --from-literal=sa_password=$MSSQL_SA_PASSWORD_DEV --namespace development 

# Add production secrets

kubectl create secret generic aqmysql --type=string --from-literal=user_password=$AQMYSQL_ANYQUIZ_USER_PASSWORD `
--from-literal=root_password=$AQMYSQL_ROOT_PASSWORD --namespace production

kubectl create secret generic gateway --type=string --from-literal=session-key=$GATEWAY_SESSION_KEY --namespace production 

kubectl create secret generic mssql --type=string --from-literal=sa_password=$MSSQL_SA_PASSWORD --namespace production 
