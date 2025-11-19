-- +migrate Down

ALTER TABLE engineer_monthly_motivations
    DROP COLUMN base_motivation_percent,
    DROP COLUMN bonus_percent,
    DROP COLUMN aggregator_payout

ALTER TABLE motivation_steps
    CHANGE COLUMN type order_type ENUM('primary', 'repeat') NOT NULL DEFAULT 'primary';

-- 3. orders: убираем из статуса новые значения, оставляем старый набор
ALTER TABLE orders
    MODIFY COLUMN status ENUM('new', 'thinking', 'in_proccess', 'working',
        'closed_without_repeat', 'sent_to_cash', 'canceled')
        DEFAULT 'new';

-- 4. reports: удаляем добавленные колонки
ALTER TABLE reports
    DROP COLUMN gave_cash,
    DROP COLUMN issued_money;

-- 5. удаляем таблицу выплат
DROP TABLE IF EXISTS engineer_payouts;

-- 6. orders: удаляем колонку aggregator_payout
ALTER TABLE orders
    DROP COLUMN aggregator_payout;

-- ====== reports ======
ALTER TABLE reports
    DROP FOREIGN KEY fk_reports_orders;

ALTER TABLE reports
    ADD COLUMN old_erp_number int;

UPDATE reports r
    JOIN orders o ON o.id = r.order_id
SET r.old_erp_number = o.erp_number;

ALTER TABLE reports
    DROP COLUMN order_id,
    CHANGE old_erp_number order_id int NOT NULL;

ALTER TABLE reports
    ADD CONSTRAINT reports_ibfk_1
        FOREIGN KEY (order_id) REFERENCES orders(erp_number)
            ON DELETE CASCADE ON UPDATE CASCADE;


-- ====== report_links ======
ALTER TABLE report_links
    DROP FOREIGN KEY fk_report_links_orders;

ALTER TABLE report_links
    ADD COLUMN old_erp_number int;

UPDATE report_links rl
    JOIN orders o ON o.id = rl.order_id
SET rl.old_erp_number = o.erp_number;

ALTER TABLE report_links
    DROP COLUMN order_id,
    CHANGE old_erp_number order_id int NOT NULL;

ALTER TABLE report_links
    ADD CONSTRAINT report_links_ibfk_1
        FOREIGN KEY (order_id) REFERENCES orders(erp_number)
            ON DELETE CASCADE ON UPDATE CASCADE;