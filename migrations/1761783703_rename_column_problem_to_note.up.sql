-- +migrate Up
#
ALTER TABLE orders RENAME COLUMN problem TO note;
ALTER TABLE orders RENAME COLUMN title TO work_volume;
ALTER TABLE orders RENAME COLUMN source_id TO aggregator_id;

CREATE TABLE report_links (
      id BIGINT AUTO_INCREMENT PRIMARY KEY,
      order_id INT NOT NULL,
      engineer_id INT NOT NULL,
      token VARCHAR(255) NOT NULL UNIQUE,
      expires_at DATETIME NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (order_id) REFERENCES orders(erp_number) ON DELETE CASCADE,
      FOREIGN KEY (engineer_id) REFERENCES engineers(id) ON DELETE CASCADE
);

CREATE TABLE reports (
                         id BIGINT AUTO_INCREMENT PRIMARY KEY,
                         order_id INT NOT NULL,
                         engineer_id INT NOT NULL,
                         has_repeat BOOLEAN NOT NULL DEFAULT FALSE,
                         repeat_date DATETIME NULL,
                         repeat_note TEXT NULL,
                         description TEXT NULL,
                         created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (order_id) REFERENCES orders(erp_number) ON DELETE CASCADE,
                         FOREIGN KEY (engineer_id) REFERENCES engineers(id) ON DELETE CASCADE
);
