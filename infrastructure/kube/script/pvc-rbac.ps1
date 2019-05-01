# The purpose of applying this config is to enable the cluster to create secrets to store
# creds needed to access the volume created dynamically by the pvc.
kubectl apply -f .\..\setup\pvc-rbac.yaml