# Microservice
This service is responsible for visitor counts related operation.

# RESTful API
- `POST /visitor/counts`, Update visitor counts based on all data in the database

  Example request:
  ```bash
  curl --location --request POST 'http://localhost:4005/visitor/counts' \
  --header 'Content-Type: application/json' \
  --data-raw '{
	
  }'
  ```

- `GET /visitor/counts/{fname}`, Reade visitor counts

  Example request:
  ```bash
  curl --location --request GET 'http://localhost:4005/visitor/counts'
  ```
# Sample Visistors Table

| Date | counts |
|---|---|---|
| "03-02-2017" | 19|
