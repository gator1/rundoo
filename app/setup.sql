CREATE DATABASE rundoo;

CREATE ROLE rundoo WITH LOGIN PASSWORD 'uber';

CREATE TABLE IF NOT EXISTS products (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    category text NOT NULL,
    sku text NOT NULL,
    version integer NOT NULL DEFAULT 1
);