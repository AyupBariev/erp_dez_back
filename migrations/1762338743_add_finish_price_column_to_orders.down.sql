-- +migrate Down

ALTER TABLE orders
    DROP COLUMN finish_price;

ALTER TABLE engineers
    DROP COLUMN is_working;
