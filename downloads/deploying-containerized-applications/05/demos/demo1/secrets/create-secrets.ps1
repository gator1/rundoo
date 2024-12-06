
# set your kubectl context *BEFORE* you run this, e.g.
# az aks get-credentials -g ps-docker-run -n ps-docker-run-aks

# prep the files in ./sources with PG connection and password - e.g.
# az postgres server show-connection-string -s ps-docker-run-postgres -u pguser -p *PASSWORD* -d wb-prod

kubectl create secret generic `
  -n wb-prod products-api-db `
  --from-file=sources/application.properties

kubectl create secret generic `
  -n wb-prod web-api `
  --from-file=sources/api.json

kubectl create secret generic `
  -n wb-prod stock-api-connection `
  --from-env-file=sources/stock-api.env