-- +migrate Down

-- 1. Переименовываем колонки обратно в исходные
ALTER TABLE orders RENAME COLUMN note TO problem;
ALTER TABLE orders RENAME COLUMN work_volume TO title;
ALTER TABLE orders RENAME COLUMN aggregator_id TO source_id;

-- 2. Удаляем таблицу report_links, если существует
DROP TABLE IF EXISTS report_links;

-- 3. Удаляем таблицу reports, если существует
DROP TABLE IF EXISTS reports;
