docker_build:
	docker build . -t shopping-cart-postgres
docker_start:
	docker run --name shopping-cart-postgres -p 5432:5432 -v ./data/postgres/:/var/lib/postgresql/data -d shopping-cart-postgres
docker_stop:
	docker stop shopping-cart-postgres
	docker rm shopping-cart-postgres

