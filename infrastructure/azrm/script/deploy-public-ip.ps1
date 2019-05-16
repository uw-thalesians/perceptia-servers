# deploy-public-ip.ps1
# Deploy public ip to be used by api services in cluster

Write-Host "Deploying public ip named: api"
New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaApi `
-Name PerceptiaPublicIpApiDeployment `
-TemplateFile ..\template\deploypublicip.json `
-TemplateParameterFile ..\parameter\deploypublicip.api.json

Write-Host "Deploying public ip named: api-dev"
New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaApi `
-Name PerceptiaPublicIpApiDevDeployment `
-TemplateFile ..\template\deploypublicip.json `
-TemplateParameterFile ..\parameter\deploypublicip.api.dev.json