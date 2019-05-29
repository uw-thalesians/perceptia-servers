Param (
        [Switch]$DeleteDeployment
)

if ($DeleteDeployment) {
        kubectl delete -f .\..\deploy\common -f .\..\deploy\prod --namespace production
} else {
        kubectl apply -f .\..\deploy\common -f .\..\deploy\prod --namespace production
}