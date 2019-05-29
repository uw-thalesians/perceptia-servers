Write-Host "Creating service principal for aks cluster"

$sp = New-AzADServicePrincipal -SkipAssignment -DisplayName PerceptiaCluster -Verbose

Set-Variable -Name SP_VALUES_DIR -Value "$env:SECRET_PERCEPTIA_SERVERS\aks\sp"
Write-Host "Saving SP values to: $SP_VALUES_DIR"
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "Secret.txt" -Value $sp.Secret.ToString()
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "ServicePrincipalNames.txt" -Value $sp.ServicePrincipalNames
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "ApplicationId.txt" -Value $sp.ApplicationId
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "DisplayName.txt" -Value $sp.DisplayName
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "Id.txt" -Value $sp.Id
New-Item -ItemType File -Path $SP_VALUES_DIR -Name "Type.txt" -Value $sp.Type

$sp.ToString()