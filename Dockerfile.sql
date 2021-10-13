FROM mysql:8.0.26
COPY ./sql/*.sql /docker-entrypoint-initdb.d