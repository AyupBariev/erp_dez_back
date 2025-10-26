-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     first_name VARCHAR(100) NOT NULL,
                                     second_name VARCHAR(100) NOT NULL,
                                     login VARCHAR(100) NOT NULL,
                                     password VARCHAR(255) NOT NULL,  -- Хэшированный пароль
                                     role ENUM('admin', 'logist') NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- таблица инженеров
CREATE TABLE IF NOT EXISTS engineers (
                                         id INT AUTO_INCREMENT PRIMARY KEY,
                                         first_name VARCHAR(100),
                                         second_name VARCHAR(100),
                                         username VARCHAR(100) NOT NULL,
                                         phone VARCHAR(20),
                                         telegram_id BIGINT UNIQUE,
                                         is_approved BOOLEAN DEFAULT FALSE,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Таблица заказов
CREATE TABLE IF NOT EXISTS orders (
                                      id INT AUTO_INCREMENT PRIMARY KEY,
                                      title VARCHAR(255) NOT NULL,
                                      description TEXT,
                                      status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
                                      engineer_id INT NOT NULL,
                                      admin_id INT NOT NULL,
                                      order_date DATE NOT NULL,
                                      order_time TIME NOT NULL,
                                      is_repeat BOOLEAN DEFAULT FALSE,
                                      payment_type ENUM('cash', 'card') DEFAULT 'cash',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      FOREIGN KEY (engineer_id) REFERENCES engineers(id),
                                      FOREIGN KEY (admin_id) REFERENCES users(id)
);

-- Таблица логов (пример третьей таблицы)
CREATE TABLE IF NOT EXISTS action_logs (
                                           id INT AUTO_INCREMENT PRIMARY KEY,
                                           user_id INT NOT NULL,
                                           action VARCHAR(50) NOT NULL,
                                           details TEXT,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           FOREIGN KEY (user_id) REFERENCES users(id)
);