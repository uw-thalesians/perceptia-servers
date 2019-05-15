# deployrg.ps1
# Deploys resource groups used by perceptia backend
Write-Host "Running script: $(Split-Path $MyInvocation.PSCommandPath -Leaf)"
Write-Host $MyInvocation.ScriptName
Write-Host "Deploying resource group: perceptiaAks"
New-AzDeployment `
-Name PerceptiaAksRgDeployment `
-Location "westus2" `
-TemplateFile ..\template\deployrg.json `
-TemplateParameterFile ..\parameter\deployrg.perceptiaAks.json

Write-Host "Deploying resource group: perceptiaApi"
New-AzDeployment `
-Name PerceptiaApiRgDeployment `
-Location "westus2" `
-TemplateFile ..\template\deployrg.json `
-TemplateParameterFile ..\parameter\deployrg.perceptiaApi.json