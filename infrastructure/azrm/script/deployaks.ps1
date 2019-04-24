Write-Host "Creating service principal for aks cluster"

$sp = New-AzADServicePrincipal -SkipAssignment
Write-Host $sp.ApplicationId

New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaAksCluster `
-Name PerceptiaClusterDeployment `
-TemplateFile ..\template\deployaks.json `
-TemplateParameterFile ..\template\deployaks.parameters.json `
-servicePrincipalClientId $(ConvertTo-SecureString $sp.ApplicationId.ToString() -AsPlainText -Force) `
-servicePrincipalClientSecret $(ConvertTo-SecureString $sp.Secret.ToString() -AsPlainText -Force)
