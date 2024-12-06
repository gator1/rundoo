
#ensure rg:
. ./resource-group.ps1

# create cluster:
$cluster = "$rg-aks"
$kubeversion = '1.19.7'
az aks create -g $rg -n $cluster `
 --node-count 2 --kubernetes-version $kubeversion

 # use this to find supported kubernetes versions:
 # az aks get-versions --location $region -o table