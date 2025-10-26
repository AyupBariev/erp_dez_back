-- +migrate Up

ALTER TABLE orders
    ADD COLUMN title VARCHAR(255) NOT NULL AFTER address
