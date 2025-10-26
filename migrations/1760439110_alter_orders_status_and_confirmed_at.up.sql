-- +migrate Up

ALTER TABLE orders
    MODIFY COLUMN status ENUM('new', 'in_progress', 'confirmed', 'closed') DEFAULT 'new',
    ADD COLUMN confirmed_at TIMESTAMP NULL DEFAULT NULL AFTER status;
