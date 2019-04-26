# This scripts create the service account for tiller on the k8 cluster
kubectl apply -f .\..\setup\helm-rbac.yaml
