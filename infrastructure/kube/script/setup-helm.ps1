#This script setups the helm rbac settings and inits helm
.\helm-rbac.create.ps1

helm init --service-account tiller