FROM postgres:15.3-alpine3.18 AS production
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB shopping-cart
COPY ./data/shopping-cart.sql docker-entrypoint-initdb.d/shopping-cart.sql
