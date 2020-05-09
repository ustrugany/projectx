# Projectx

### Starting backing services
1. for cassandra run `docker-compose -f docker-compose.backing.yml up --build `

### Starting api 
1. docker-compose
    - run `docker-compose build`
    - run `docker-compose up`
2. docker
    - update cassandra connection details in `.env` 
    - run ` docker run -p 8090:8090 --env-file=.env ${IMAGE}:latest server --port 8090`
        - where `${IMAGE}` e.g. `projectx_api` or `piotras/projectx:latest`


### Postman collection
- import from `./projectx.postman_collection.json`

