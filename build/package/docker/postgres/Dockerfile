FROM postgres:14

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=12345
ENV POSTGRES_DB=campaing-consumer-api-db

EXPOSE 5432:5432

ARG LOCAL_SQL=postgres/scripts/local
ARG MIGRATION_SQL=postgres/scripts/migrations

COPY $LOCAL_SQL /docker-entrypoint-initdb.d/
COPY $MIGRATION_SQL /docker-entrypoint-initdb.d/