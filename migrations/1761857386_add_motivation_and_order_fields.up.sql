-- +migrate Up

-- 1. Добавляем колонки в orders
ALTER TABLE orders
    ADD COLUMN repeat_id INT NULL,
    ADD COLUMN repeat_description VARCHAR(255) NULL,
    ADD COLUMN repeated_by VARCHAR(50) NULL;


-- 2. Добавляем внешний ключ
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_repeat
        FOREIGN KEY (repeat_id) REFERENCES orders(id)
            ON DELETE SET NULL;

-- 2. Создаем таблицу motivation_steps
CREATE TABLE motivation_steps (
                                  id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                  name VARCHAR(50) NOT NULL,       -- первичка/повторка/бонус
                                  min_amount DECIMAL(15,2) NOT NULL,
                                  percent DECIMAL(5,2) NOT NULL,
                                  sort INT NOT NULL,
                                  order_type ENUM('primary', 'repeat', 'bonus') NOT NULL DEFAULT 'primary',
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO motivation_steps (name, min_amount, percent, sort, order_type)
VALUES
    ('Первичка', 20000, 10, 1, 'primary'),
    ('Повторка', 10000, 15, 2, 'repeat'),
    ('Первичка', 50000, 20, 3, 'primary'),
    ('Повторка', 20000, 25, 4, 'repeat'),
    ('Бонус за 100к+', 100000, 5, 5, 'bonus');

-- 3. Создаем таблицу engineer_motivation_targets
CREATE TABLE engineer_motivation_targets (
                                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                             engineer_id BIGINT NOT NULL,
                                             motivation_percent DECIMAL(5,2) NOT NULL,  -- 10,15,20,25,30
                                             target_month DATE NOT NULL,
                                             confirmed_by_admin BOOLEAN DEFAULT FALSE,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 4. Создаем таблицу engineer_monthly_motivations
CREATE TABLE engineer_monthly_motivations (
                                              id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                              engineer_id BIGINT NOT NULL,
                                              month DATE NOT NULL, -- первый день месяца

                                              reports_count INT NOT NULL DEFAULT 0,           -- всего отчетов
                                              primary_orders_count INT NOT NULL DEFAULT 0,     -- количество первичных заказов
                                              repeat_orders_count INT NOT NULL DEFAULT 0,      -- количество повторных заказов
                                              orders_total_amount DECIMAL(15,2) NOT NULL DEFAULT 0, -- сумма первичных заказов
                                              repeat_orders_amount DECIMAL(15,2) NOT NULL DEFAULT 0, -- сумма повторов
                                              gross_profit DECIMAL(15,2) NOT NULL DEFAULT 0,
                                              average_check DECIMAL(15,2) NOT NULL DEFAULT 0, -- по первичным заказам

                                              motivation_percent DECIMAL(5,2) NOT NULL DEFAULT 0, -- итоговый %
                                              total_motivation DECIMAL(15,2) NOT NULL DEFAULT 0,  -- сумма мотивации
                                              confirmed_by_admin BOOLEAN DEFAULT FALSE,

                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                              UNIQUE KEY engineer_month_unique (engineer_id, month)
);