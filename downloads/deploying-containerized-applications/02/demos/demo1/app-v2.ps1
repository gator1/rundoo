
# create app network:
docker network create wired-brain

# database:
docker run -d --network wired-brain `
 --name products-db-test `
 -e POSTGRES_PASSWORD=wiredtest `
 psdockerrun/products-db

# products API:
docker run -d --network wired-brain `
 --name products-api-test `
 -p 8081:80 `
 -v "$(pwd)/config/products-api:/app/config" `
 psdockerrun/products-api

# stock API:
docker run -d --network wired-brain `
 --name stock-api-test `
 -p 8082:8080 `
 -e POSTGRES_CONNECTION_STRING="host=products-db-test port=5432 user=postgres password=wiredtest dbname=postgres sslmode=disable" `
 psdockerrun/stock-api

# website
docker run -d --network wired-brain `
 --name web-test `
 -p 8080:80 `
 -e Environment=TEST `
 -v "$(pwd)/config/web/config:/app/config" `
 -v "$(pwd)/config/web/secrets:/app/secrets" `
 psdockerrun/web