# Backend Golang for Developtment
start: swagger
	go run ./cmd/api/
swagger:
	swag fmt
	swag init -g ./cmd/api/main.go  --output ./docs/ --parseDependency

# Docker Backend Golang
docker_be_build:
	docker build . -t shopping-cart-backend
docker_be_start:
	docker run --name shopping-cart-backend -p 4000:4000 -d shopping-cart-backend
docker_be_stop:
	docker stop shopping-cart-backend
	docker rm shopping-cart-backend

# Docker Database Postgresql
docker_db_build:
	docker build -f ./db/Dockerfile ./db/ -t shopping-cart-postgres
docker_db_start:
	docker run --restart unless-stopped --name shopping-cart-postgres -p 5432:5432 -d shopping-cart-postgres
	# Use this command if you want to save the db data.
	# docker run --restart unless-stopped --name shopping-cart-postgres -p 5432:5432 -v ./db/data/postgres/:/var/lib/postgresql/data -d shopping-cart-postgres
docker_db_stop:
	docker stop shopping-cart-postgres
	docker rm shopping-cart-postgres

# Docker Compose Backend and Database Postgresql
down:
	docker compose down
up:
	docker compose up --build -d
