Param (
        [Switch]$DeleteDeployment
)

if ($DeleteDeployment) {
        kubectl delete -f .\..\deploy\common -f .\..\deploy\dev --namespace development
} else {
        kubectl apply -f .\..\deploy\common -f .\..\deploy\dev --namespace development
}