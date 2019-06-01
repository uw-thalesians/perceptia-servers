Param (
        [Switch]$DeleteDeployment,
        [switch]$DeletePVC,
        [switch]$Prod
)


if ($Prod) {
        Set-Variable -Name NAMESPACE -Value "production" 
        Set-Variable -Name DEPLOY_DIR -Value "prod"
        Write-Host "Prod switch set, using $NAMESPACE namespace for all actions..."
} else {
        Set-Variable -Name NAMESPACE -Value "development"
        Set-Variable -Name DEPLOY_DIR -Value "dev"
        Write-Host "Prod switch not set, defaulting to $NAMESPACE namespace for all actions..."
}





if ($DeleteDeployment) {
        Write-Host "`n"
        Write-Host "Deleting Deployment for namespace: $NAMESPACE"
        kubectl delete -f "./../deploy/common" -f "./../deploy/$DEPLOY_DIR" --namespace $NAMESPACE
        
        if ($DeletePVC) {
                Write-Host "`n"
                Write-Host "Sleeping 15 seconds to allow deployment to finish deleting..."
                Start-Sleep -Seconds 15
                Write-Host "DeletePVC switch set, deleting all PVCs created with the ./../setup/pvc.yaml config..."
                Write-Host "`n"
                kubectl delete -f "./../setup/pvc.yaml" --namespace $NAMESPACE
        } else {
                Write-Host "`n"
                Write-Host "DeletePVC switch not set, existing PVCs created with the ./../setup/pvc.yaml config will not be deleted..."
                Write-Host "`n"
        }
} else {
        Write-Host "`n"
        Write-Host "Deploying for namespace: $NAMESPACE"
        kubectl apply -f "./../setup/pvc.yaml" -f "./../deploy/common" -f "./../deploy/$DEPLOY_DIR" --namespace $NAMESPACE
}