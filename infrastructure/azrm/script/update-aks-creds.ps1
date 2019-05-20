param(
        [string]
        $servicePrincipalClientSecretFile = "$env:SECRET_PERCEPTIA_SERVERS\aks\sp\Secret.txt"
)

# Get service principal ID
$SP_ID=$(az aks show --resource-group perceptiaAks --name perceptiaCluster --query "servicePrincipalProfile.clientId" --output tsv)

# Reset sp password
$SP_SECRET= (Get-Content -Path $servicePrincipalClientSecretFile)

az aks update-credentials `
--resource-group perceptiaAks `
--name perceptiaCluster `
--reset-service-principal `
--service-principal $SP_ID `
--client-secret $SP_SECRET