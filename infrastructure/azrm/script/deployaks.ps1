New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaAksCluster `
-Name PerceptiaClusterDeployment `
-TemplateFile ..\template\deployaks.json `
-TemplateParameterFile ..\template\deployaks.parameters.json
