-- +migrate Up

ALTER TABLE engineer_monthly_motivations
    ADD COLUMN base_motivation_percent DECIMAL(5,2) DEFAULT 0 AFTER gross_profit,
    ADD COLUMN bonus_percent DECIMAL(5,2) DEFAULT 0 AFTER base_motivation_percent,
    ADD COLUMN aggregator_payout DECIMAL(10,2) DEFAULT 0 AFTER bonus_percent;
ALTER TABLE motivation_steps
    CHANGE COLUMN order_type type ENUM('primary', 'repeat', 'bonus') NOT NULL DEFAULT 'primary';

ALTER TABLE orders
    MODIFY COLUMN status ENUM('new', 'thinking', 'in_proccess', 'working',
        'closed_without_repeat', 'sent_to_cash', 'closed_finally', 'canceled') DEFAULT 'new';

ALTER TABLE reports
    ADD COLUMN gave_cash   DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT 'сдано в кассу',
    ADD COLUMN issued_money DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT 'выдано инженеру';

CREATE TABLE engineer_payouts (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  engineer_id BIGINT NOT NULL,
                                  month DATE NOT NULL,          -- первый день месяца
                                  prepayment DECIMAL(10,2) NOT NULL,          -- 50 % от salary
                                  paid_prepayment DECIMAL(10,2) DEFAULT 0,    -- уже выдано авансом
                                  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                  UNIQUE KEY uniq_eng_month (engineer_id, month)
);

ALTER TABLE orders
    ADD COLUMN aggregator_payout DECIMAL(10,2) NOT NULL DEFAULT 0
        COMMENT 'сколько мы платим агрегатору за этот заказ';

-- ====== reports ======

ALTER TABLE reports
    DROP FOREIGN KEY reports_ibfk_1;
ALTER TABLE reports
    ADD COLUMN new_order_id INT;

UPDATE reports r
    JOIN orders o ON o.erp_number = r.order_id
SET r.new_order_id = o.id;

ALTER TABLE reports
    DROP COLUMN order_id,
    CHANGE new_order_id order_id INT NOT NULL;

ALTER TABLE reports
    ADD CONSTRAINT fk_reports_orders
        FOREIGN KEY (order_id) REFERENCES orders(id)
            ON DELETE CASCADE ON UPDATE CASCADE;


-- ====== report_links ======
-- приводим тип к тому же, что и orders.id
ALTER TABLE report_links
    DROP FOREIGN KEY report_links_ibfk_1;

ALTER TABLE report_links
    ADD COLUMN new_order_id INT;

UPDATE report_links rl
    JOIN orders o ON o.erp_number = rl.order_id
SET rl.new_order_id = o.id;

ALTER TABLE report_links
    DROP COLUMN order_id,
    CHANGE new_order_id order_id INT NOT NULL;

ALTER TABLE report_links
    ADD CONSTRAINT fk_report_links_orders
        FOREIGN KEY (order_id) REFERENCES orders(id)
            ON DELETE CASCADE ON UPDATE CASCADE;