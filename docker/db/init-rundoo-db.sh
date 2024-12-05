#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL


    CREATE TABLE public.products (
        id bigint NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        name character varying(255) NULL,
        category character varying(255) NULL,
        sku character varying(255) NULL,
        version bigint NOT NULL
    );

    ALTER TABLE public.products ADD CONSTRAINT products_pkey PRIMARY KEY (id);

    INSERT INTO  "public"."products" ("id", "name", "category", "sku", "version")
        VALUES (1, 'Espresso', 'beverage', 'EV21001', 1);

    INSERT INTO  "public"."products" ("id", "name", "category", "sku", "version")
        VALUES (2, 'Americano', 'beverage', 'EV21002', 2);

    INSERT INTO  "public"."products" ("id", "name", "category", "sku", "version")
        VALUES (3, 'Flat White', 'beverage', 'PP1837638845', 1);

EOSQL
