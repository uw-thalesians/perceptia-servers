Param (
        [Switch]$DeleteDeployment,
        [switch]$DeletePVC,
        [switch]$Prod
)

Set-Variable -Name NAMESPACE -Value "development"
Set-Variable -Name DEPLOY_DIR -Value "dev"
if ($Prod) {
        Set-Variable -Name NAMESPACE -Value "production" 
        Set-Variable -Name DEPLOY_DIR -Value "prod"
}



if ($DeletePVC) {
        kubectl delete -f "./../setup/pvc.yaml" --namespace $NAMESPACE
}

if ($DeleteDeployment) {
        Write-Host "Deleting Deployment for namespace: $NAMESPACE"
        kubectl delete -f "./../deploy/common" -f "./../deploy/$DEPLOY_DIR" --namespace $NAMESPACE
} else {
        Write-Host "Deploying for namespace: $NAMESPACE"
        kubectl apply -f "./../setup/pvc.yaml" -f "./../deploy/common" -f "./../deploy/$DEPLOY_DIR" --namespace $NAMESPACE
}