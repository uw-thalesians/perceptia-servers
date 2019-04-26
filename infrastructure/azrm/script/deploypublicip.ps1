New-AzResourceGroupDeployment `
-ResourceGroupName perceptiaAksCluster `
-Name PerceptiaPublicIpDeployment `
-TemplateFile ..\template\deploypublicip.json `
-TemplateParameterFile ..\template\deploypublicip.parameters.json