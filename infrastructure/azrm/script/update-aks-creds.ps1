# Get service principal ID
$SP_ID=$(az aks show --resource-group perceptiaAks --name perceptiaCluster --query "servicePrincipalProfile.clientId" --output tsv)

# Reset sp password
$SP_SECRET=$(az ad sp credential reset --name $SP_ID --query password -o tsv)

az aks update-credentials `
--resource-group perceptiaAks `
--name perceptiaCluster `
--reset-service-principal `
--service-principal $SP_ID `
--client-secret $SP_SECRET