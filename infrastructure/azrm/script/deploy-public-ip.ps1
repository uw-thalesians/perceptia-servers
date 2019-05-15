New-AzResourceGroupDeployment `
-ResourceGroupName MC_perceptiaAksCluster_perceptiaCluster_westus2 `
-Name PerceptiaPublicIpDeployment `
-TemplateFile ..\template\deploypublicip.json `
-TemplateParameterFile ..\template\deploypublicip.parameters.json