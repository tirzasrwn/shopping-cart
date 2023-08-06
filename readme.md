# shopping-cart

## About

An shopping-cart API for study case using go, go-gin, swaggo, and postgresql. Main features for this study case are:

- User can view product list by product category
- User can add product to shopping cart
- User can see list of products that have been added to the shopping cart
- User can delete product list in shopping cart
- User can checkout and make payment transactions
- Login and register user

## Stack

- Go
- Gin
- Swagger API Documentation
- JWT
- Postgres
- Database Migration

## Requirement

- Unix based OS (for make command)
- Docker
- Docker Compose
- Make
- Go (for development)

## Running

- Running using docker compose

```sh
git clone https://github.com/tirzasrwn/shopping-cart.git
cd shopping-cart
# start docker compose
make up
# stop docker compose
make down
```

- Running using go command, typically for development

```sh
make start
```

- Running spesific docker image
  There are two docker images, backend go and database postgresql.
  You can run them separately.

```sh
make docker_<service>_build
make docker_<service>_start
make docker_<service>_stop
# where <service> is db for database and be for backend
```
