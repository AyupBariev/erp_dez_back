-- +migrate Up

ALTER TABLE repeat_requests
    ADD COLUMN confirmed_by BIGINT NULL,
    ADD COLUMN status VARCHAR(255) DEFAULT 'pending',
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE INDEX idx_repeat_requests_deleted_at ON repeat_requests(deleted_at);
