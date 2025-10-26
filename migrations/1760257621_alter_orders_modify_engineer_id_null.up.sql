-- +migrate Up

ALTER TABLE orders MODIFY COLUMN engineer_id INT NULL;