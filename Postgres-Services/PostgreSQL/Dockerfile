
FROM postgres:latest


ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=example_db


COPY init.sql /docker-entrypoint-initdb.d/


EXPOSE 5432
