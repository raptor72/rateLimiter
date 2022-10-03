#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
     create schema if not exists $SCHEMA;
     GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;
     CREATE TABLE IF NOT EXISTS white_lists (
         id SERIAL PRIMARY KEY,
         address INET
     );
     CREATE TABLE IF NOT EXISTS black_lists (
         id SERIAL PRIMARY KEY,
         address INET
     );
EOSQL