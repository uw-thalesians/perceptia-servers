# add-roll-sp.ps1
# Adds role to sp used by Aks cluster
Write-Host "Adding Network Contributor role to perceptiaCluster sp for perceptiaApi rg"

New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaApi `
-Name PerceptiaAksSpRoleApiDeployment `
-TemplateFile ..\template\addrole.json `
-TemplateParameterFile "..\parameter\addrole.perceptiaApi.json" `
-roleNameGuid  $(New-Guid)