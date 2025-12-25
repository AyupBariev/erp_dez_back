ALTER TABLE reports
    ADD COLUMN motivation_percent DECIMAL(5,2) DEFAULT 0 AFTER description,
    ADD COLUMN motivation_step_id BIGINT NULL AFTER motivation_percent,
    ADD CONSTRAINT fk_motivation_step FOREIGN KEY (motivation_step_id) REFERENCES motivation_steps(id);

ALTER TABLE motivation_steps ADD COLUMN percent_increment DECIMAL(5,2) DEFAULT 0 AFTER percent;
ALTER TABLE repeat_requests
    ADD COLUMN created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE reports
    ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
