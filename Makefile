# Backend Go for development.
start: swagger
	go run ./cmd/api/
swagger:
	swag fmt
	swag init -g ./cmd/api/main.go  --output ./docs/ --parseDependency

# Docker network.
docker_net_start:
	- docker network create --driver bridge shopping-cart
docker_net_stop:
	- docker network remove shopping-cart

# Docker backend Go.
docker_be_build:
	docker build . -t shopping-cart-backend:v1.2.0
docker_be_start: docker_net_start
	docker run --network=shopping-cart --name shopping-cart-backend -p 4000:4000 -d shopping-cart-backend:v1.2.0
docker_be_stop:
	- docker stop shopping-cart-backend
	- docker rm shopping-cart-backend

# Docker database PostgreSLQ.
docker_db_build:
	docker build -f ./db/Dockerfile ./db/ -t shopping-cart-postgres:v1.2.0
docker_db_start: docker_net_start
	docker run --hostname shopping-cart-postgres --network=shopping-cart --restart unless-stopped --name shopping-cart-postgres -p 5432:5432 -d shopping-cart-postgres:v1.2.0
	# Use this command if you want to save the db data.
	# docker run --hostname shopping-cart-postgres --network=shopping-cart --restart unless-stopped --name shopping-cart-postgres -p 5432:5432 -v ./db/data/.postgres/:/var/lib/postgresql/data -d shopping-cart-postgres:v1.2.0
docker_db_stop:
	- docker stop shopping-cart-postgres
	- docker rm shopping-cart-postgres

docker_stop: docker_be_stop docker_db_stop docker_net_stop

docker_start: docker_net_start docker_db_build docker_db_start docker_be_build docker_be_start

# Docker compose backend Go and database PostgreSQL.
down:
	docker compose down
up:
	docker compose up --build -d
