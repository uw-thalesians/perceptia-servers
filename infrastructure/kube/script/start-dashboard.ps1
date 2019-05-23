# param(
#         [switch]$Stop
# )

az aks browse --resource-group perceptiaAks --name perceptiaCluster

# if ($Stop) {
#         Stop-AzAksDashboard
# } else {
#         Start-AzAksDashboard -ResourceGroupName perceptiaAks -Name perceptiaCluster
# }