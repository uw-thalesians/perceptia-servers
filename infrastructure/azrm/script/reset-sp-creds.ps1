# Get service principal ID
$SP_APP_ID="c7fca245-ee6b-486c-9437-bb1b60907bcd"

# Reset sp password
$SP_SECRET=$(az ad sp credential reset --name $SP_APP_ID --query password -o tsv)

$env:SP_SECRET=$SP_SECRET

Write-Host $env:SP_SECRET