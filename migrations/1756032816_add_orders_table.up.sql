-- +migrate Up

-- === Обновляем таблицу orders ===
ALTER TABLE orders
    DROP COLUMN title,
    DROP COLUMN payment_type,
    DROP COLUMN order_date,
    DROP COLUMN order_time,
    DROP COLUMN description,
    DROP COLUMN is_repeat,
    DROP COLUMN status;

ALTER TABLE orders
    ADD COLUMN erp_number INT NOT NULL UNIQUE AFTER id,  -- Заполняется в Go отсчет начинается с 100 000 +
    ADD COLUMN source_id INT NOT NULL AFTER erp_number,
    ADD COLUMN our_percent DECIMAL(5,2) NOT NULL AFTER source_id,
    ADD COLUMN client_name VARCHAR(255) NOT NULL AFTER our_percent,
    ADD COLUMN phones JSON AFTER client_name,
    ADD COLUMN address VARCHAR(500) AFTER phones,
    ADD COLUMN problem TEXT AFTER address,
    ADD COLUMN scheduled_at TIMESTAMP AFTER problem,
    ADD COLUMN status ENUM('new', 'in_progress', 'closed') DEFAULT 'new' AFTER scheduled_at;

-- === Таблица агрегаторов ===
CREATE TABLE IF NOT EXISTS aggregators (
                                           id INT AUTO_INCREMENT PRIMARY KEY,
                                           name VARCHAR(255) NOT NULL,
                                           priority INT DEFAULT 0,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- === Таблица запросов повтора ===
CREATE TABLE IF NOT EXISTS repeat_requests (
                                               id INT AUTO_INCREMENT PRIMARY KEY,
                                               order_id INT NOT NULL,
                                               engineer_id INT NOT NULL,
                                               description TEXT,
                                               requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               scheduled_at TIMESTAMP,
                                               confirmed BOOLEAN DEFAULT FALSE,
                                               confirmed_at TIMESTAMP NULL,
                                               repeat_order_id INT NULL,
                                               FOREIGN KEY (order_id) REFERENCES orders(id),
                                               FOREIGN KEY (engineer_id) REFERENCES engineers(id),
                                               FOREIGN KEY (repeat_order_id) REFERENCES orders(id)
);

-- === Таблица выплат агрегаторам ===
CREATE TABLE IF NOT EXISTS payouts (
                                       id INT AUTO_INCREMENT PRIMARY KEY,
                                       aggregator_id INT NOT NULL,
                                       amount DECIMAL(10,2) NOT NULL,
                                       payout_date DATE NOT NULL,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       FOREIGN KEY (aggregator_id) REFERENCES aggregators(id)
);
