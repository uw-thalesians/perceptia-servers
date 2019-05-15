Write-Host "Creating service principal for aks cluster"

$sp = New-AzADServicePrincipal -SkipAssignment -DisplayName perceptiaCluster -Verbose
Write-Host 