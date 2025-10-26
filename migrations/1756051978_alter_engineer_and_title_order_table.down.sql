-- +migrate Down
ALTER TABLE orders
    DROP COLUMN title;