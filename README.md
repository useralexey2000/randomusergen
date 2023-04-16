Client for https://randomuser.me/api

accepts http POST request with optional query
and makes http  GET request to the resource
and saves the result to PostgreSQL

query keys can be found https://randomuser.me/documentation

run:
    docker-compose up -d

request example:
curl -i \
    -H "Accept: application/json" \
    -X POST -d '{"query": {"results": "3"}}'
    http://localhost:8080/api/v1/create