CREATE TABLE IF NOT EXISTS "warehouses" (
    "id" uuid PRIMARY KEY,
    "name" varchar NOT NULL,
    "street" varchar NOT NULL,
    "city" varchar NOT NULL,
    "state" varchar NOT NULL,
    "zip_code" varchar NOT NULL,
    "is_main_warehouse" boolean NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "warehouse_products" (
    "id" uuid PRIMARY KEY,
    "warehouse_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "product_name" varchar NOT NULL,
    "product_quantity" integer NOT NULL,
    "description" text NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "stock_movement" (
    "id" uuid PRIMARY KEY,
    "product_id" uuid NOT NULL,
    "product_name" varchar NOT NULL,
    "quantity" integer NOT NULL,
    "from_warehouse_id" uuid NOT NULL,
    "to_warehouse_id" uuid,
    "to_user_id" uuid,
    "created_at" timestamp NOT NULL
);

CREATE INDEX `warehouse_is_main_warehouse_idx` ON `warehouses` (`is_main_warehouse`);

CREATE INDEX `warehouse_products_product_id_idx` ON `warehouse_products` (`product_id`);
CREATE INDEX `warehouse_products_warehouse_id_idx` ON `warehouse_products` (`warehouse_id`);

ALTER TABLE `warehouse_products` ADD FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE `stock_movement` ADD FOREIGN KEY (`from_warehouse_id`) REFERENCES `warehouses` (`id`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `stock_movement` ADD FOREIGN KEY (`to_warehouse_id`) REFERENCES `warehouses` (`id`) ON UPDATE CASCADE ON DELETE CASCADE;