    CREATE TABLE public.products (
        id bigint NOT NULL,
        name character varying(255) NULL,
        price numeric(19,2) NULL,
        stock bigint NOT NULL
    );

    ALTER TABLE public.products ADD CONSTRAINT products_pkey PRIMARY KEY (id);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (1, 'Espresso', 5.00, 600);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (2, 'Americano', 6.00, 400);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (3, 'Flat White', 7.50, 750);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (4, 'Cappuccino', 5.50, 100);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (5, 'Latte', 6.50, 100);

    INSERT INTO "public"."products" ("id", "name", "price", "stock")
        VALUES (6, 'Cortado', 6.00, 100);