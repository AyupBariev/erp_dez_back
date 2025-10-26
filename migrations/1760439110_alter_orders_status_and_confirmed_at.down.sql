-- +migrate Down

ALTER TABLE orders
    DROP COLUMN confirmed_at,
    MODIFY COLUMN status ENUM('new', 'in_progress', 'closed') DEFAULT 'new';
