-- +migrate Up

CREATE TABLE problems (
                          id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(255) NOT NULL UNIQUE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE orders
    ADD COLUMN problem_id BIGINT UNSIGNED NULL,
    ADD CONSTRAINT fk_orders_problem
        FOREIGN KEY (problem_id) REFERENCES problems(id)
            ON DELETE SET NULL;

-- если source_id не связан с aggregators
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_aggregator
        FOREIGN KEY (source_id) REFERENCES aggregators(id)
            ON DELETE RESTRICT;

ALTER TABLE orders
    ADD COLUMN price VARCHAR(256) DEFAULT 0 AFTER our_percent;


INSERT INTO problems (name) VALUES
    ('Кроты'), ('Землеройки'), ('Крысы'), ('Мыши'), ('Осы'),
    ('Шершни'), ('Запах'), ('Холодильник'), ('Труп'), ('Пожар'),
    ('Комары'), ('Аккарицидная обработка'), ('Клещи'), ('Муравьи'),
    ('Моль'), ('Прочие насекомые'), ('Кожеед'), ('Мукоед'), ('Слизни'),
    ('Улитки'), ('Мухи'), ('Клопы'), ('Тараканы'), ('Блохи'),
    ('Плесень'), ('Змеи'), ('Демеркуризация'), ('Обработка деревьев'),
    ('Обработка кустов'), ('Борщевик'), ('Одуванчики'), ('Другие сорняки');

INSERT INTO aggregators (name) VALUES
    ('Бугатти'),('Maxon'), ('Володя'), ('Заявочная'), ('Ооодез'),('Витамин'),
    ('Serezha'), ('Arseny'),('АлександрК'),('ИльяС'),('Балаш'), ('Denisгранит'),
    ('Сергей Александрович');

ALTER TABLE orders
    MODIFY COLUMN status ENUM('new', 'thinking', 'in_proccess', 'working',
        'closed_without_repeat', 'closed_finally', 'canceled') DEFAULT 'new';
