-- +migrate Up

ALTER TABLE orders
    ADD COLUMN finish_price VARCHAR(256) DEFAULT 0 AFTER price;

ALTER TABLE engineers
    ADD COLUMN is_working BOOLEAN DEFAULT 0;