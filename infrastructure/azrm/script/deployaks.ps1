Write-Host "This script is not working. Creates SP secret wrong"
exit(1)
Write-Host "Creating service principal for aks cluster"

$sp = New-AzADServicePrincipal -SkipAssignment -DisplayName perceptiaCluster
Write-Host $sp.ApplicationId

Write-Host "Deploying AKS Cluster: perceptiaCluster"
New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaAksCluster `
-Name PerceptiaClusterDeployment `
-TemplateFile ..\template\deployaks.json `
-TemplateParameterFile ..\template\deployaks.parameters.json `
-servicePrincipalClientId $(ConvertTo-SecureString $sp.ApplicationId.ToString() -AsPlainText -Force) `
-servicePrincipalClientSecret $(ConvertTo-SecureString $sp.Secret.ToString() -AsPlainText -Force)
