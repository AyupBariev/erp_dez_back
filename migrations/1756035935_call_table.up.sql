-- +migrate Up

CREATE TABLE IF NOT EXISTS calls (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     order_id INT NULL,
                                     caller_role ENUM('si','logist','admin') NOT NULL,
                                     target_role ENUM('client','logist','admin') NOT NULL,
                                     caller_number VARCHAR(20) NOT NULL,
                                     target_number VARCHAR(20) NOT NULL,
                                     call_sid VARCHAR(255) NULL, -- идентификатор звонка в Telphin
                                     status ENUM('queued','ringing','completed','failed') DEFAULT 'queued',
                                     started_at TIMESTAMP NULL,
                                     ended_at TIMESTAMP NULL,
                                     recording_url VARCHAR(500) NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     FOREIGN KEY (order_id) REFERENCES orders(id)
);