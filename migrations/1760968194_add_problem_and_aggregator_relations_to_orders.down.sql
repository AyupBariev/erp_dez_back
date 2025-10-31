-- +migrate Down

-- 1. Удаляем внешние ключи из orders
ALTER TABLE orders DROP FOREIGN KEY fk_orders_problem;
ALTER TABLE orders DROP FOREIGN KEY fk_orders_aggregator;

-- 2. Удаляем добавленные колонки из orders
ALTER TABLE orders DROP COLUMN problem_id;
ALTER TABLE orders DROP COLUMN price;

-- 3. Удаляем таблицу problems
DROP TABLE IF EXISTS problems;
