<#
 .SYNOPSIS
    Deploys a template to Azure

 .DESCRIPTION
    Deploys an Azure Resource Manager template

 .PARAMETER subscriptionId
    The subscription id where the template will be deployed.

 .PARAMETER resourceGroupName
    The resource group where the template will be deployed. Can be the name of an existing or a new resource group.

 .PARAMETER resourceGroupLocation
    A resource group location. If specified, will try to create a new resource group in this location. If not specified, assumes resource group is existing.

 .PARAMETER deploymentName
    The deployment name.

 .PARAMETER templateFilePath
    Path to the template file. Defaults to template.json.

 .PARAMETER parametersFilePath
    Path to the parameters file. Defaults to parameters.json. If file is not found, will prompt for parameter values based on template.

 .PARAMETER servicePrincipalClientSecretFile
    Path to the file containing the sp secret file. Defaults to parameters.json. If file is not found, will prompt for parameter values based on template.
#>

param(
 [string]
 $resourceGroupName = "perceptiaAks",

 [string]
 $deploymentName = "PerceptiaClusterDeployment",

 [string]
 $templateFilePath = ".\..\template\deployaks.json",

 [string]
 $parametersFilePath = ".\..\parameter\deployaks.e2sv3.json",

 [string]
 $servicePrincipalClientSecretFile = "$env:SECRET_PERCEPTIA_SERVERS\aks\sp\Secret.txt"
)


#******************************************************************************
# Script body
# Execution begins here
#******************************************************************************
$ErrorActionPreference = "Stop"

$SECRET = (Get-Content -Path $servicePrincipalClientSecretFile)

# Start the deployment
Write-Host "Starting deployment...";
New-AzResourceGroupDeployment -ResourceGroupName $resourceGroupName -Name $deploymentName -TemplateFile $templateFilePath -TemplateParameterFile $parametersFilePath -servicePrincipalClientSecret $SECRET;
