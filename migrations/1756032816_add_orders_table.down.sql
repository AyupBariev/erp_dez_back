-- +migrate Down

DROP TABLE IF EXISTS payouts;
DROP TABLE IF EXISTS repeat_requests;
DROP TABLE IF EXISTS aggregators;

ALTER TABLE orders
    DROP COLUMN erp_number,
    DROP COLUMN source_id,
    DROP COLUMN our_percent,
    DROP COLUMN client_name,
    DROP COLUMN phones,
    DROP COLUMN address,
    DROP COLUMN problem,
    DROP COLUMN scheduled_at,
    DROP COLUMN status;

ALTER TABLE orders
    ADD COLUMN title VARCHAR(255) NOT NULL AFTER id,
    ADD COLUMN description TEXT AFTER title,
    ADD COLUMN status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' AFTER description,
    ADD COLUMN order_date DATE NOT NULL AFTER status,
    ADD COLUMN order_time TIME NOT NULL AFTER order_date,
    ADD COLUMN is_repeat BOOLEAN DEFAULT FALSE AFTER order_time,
    ADD COLUMN payment_type ENUM('cash', 'card') DEFAULT 'cash' AFTER is_repeat;
