
-- +migrate Down

ALTER TABLE orders DROP FOREIGN KEY fk_orders_repeat;
-- 1. Удаляем колонки из orders
ALTER TABLE orders
    DROP COLUMN repeat_id,
    DROP COLUMN repeat_description,
    DROP COLUMN repeated_by;


-- 2. Удаляем таблицы
DROP TABLE IF EXISTS motivation_steps;
DROP TABLE IF EXISTS engineer_motivation_targets;
DROP TABLE IF EXISTS engineer_monthly_motivations;
